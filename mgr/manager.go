/*
*	Copyright (C) 2025 Kendall Tauser
*
*	This program is free software; you can redistribute it and/or modify
*	it under the terms of the GNU General Public License as published by
*	the Free Software Foundation; either version 2 of the License, or
*	(at your option) any later version.
*
*	This program is distributed in the hope that it will be useful,
*	but WITHOUT ANY WARRANTY; without even the implied warranty of
*	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*	GNU General Public License for more details.
*
*	You should have received a copy of the GNU General Public License along
*	with this program; if not, write to the Free Software Foundation, Inc.,
*	51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.
 */

package mgr

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/vault/api"
	k8sauth "github.com/hashicorp/vault/api/auth/kubernetes"
	"k8s.io/klog/v2"
)

// Global method used for initializing an APIManager instance. This includes registering
// signal handlers, creating SysAPI, and starting up all the required subsystems as according to the
// provided registrar.
func (m *APIManager) Initialize(registrar *SystemRegistrar) {
	if registrar == nil {
		klog.Error("nil registrar pointer provided to the process")
		return
	}

	// Tell the runtime to forward signals from the OS to this channel for downstream processing.
	signal.Notify(m.sigHandle)

	m.registrar = registrar

	for _, sys := range registrar.Systems {
		m.systems[sys.Name()] = sys
	}

	for name, sys := range m.systems {
		configs := *sys.Configs()
		klog.V(5).Infof("registering %d config keys for subsystem %s", len(configs), name)

		m.ckeys = append(m.ckeys, configs...)

		secrets := *sys.Secrets()
		klog.V(5).Infof("registering %d secret keys for subsystem %s", len(secrets), name)

		m.skeys = append(m.skeys, secrets...)
	}

	//
	if m.opts.EnableSysAPI {
		m.ckeys = append(m.ckeys, sysAPIListenAddress)
		m.ckeys = append(m.ckeys, sysAPIListenPort)
		m.ckeys = append(m.ckeys, sysAPIConcurrency)
		m.ckeys = append(m.ckeys, sysAPIIdleTimeout)
		m.ckeys = append(m.ckeys, sysAPIReadTimeout)
		m.ckeys = append(m.ckeys, sysAPIWriteTimeout)
		m.ckeys = append(m.ckeys, sysAPIWriteBufferSize)
		m.ckeys = append(m.ckeys, sysAPIReadBufferSize)
	}

	// read in configuration and secrets before booting further, or at least attempt to.
	m.initConfigs()

	if m.opts.EnableVault {
		m.initVault()
	}

	// Set up the sysAPI and all its handlers.
	if m.opts.EnableSysAPI {
		m.initSysAPI()
	}

	m.preInit()
	m.initializeSubsystems(registrar)
	m.postInit()

	// register collectors with the registry
	if registrar.Registration != nil && m.registry != nil {
		klog.V(5).Infof("registering collectors with manager registry")
		registrar.Registration.RegisterCollectors(m.registry)
	}
}

// Sync start process should never return until process shutdown has been confirmed and all
// subsystems have exited as gracefully as possible.
func (m *APIManager) SyncStartProcess() {
	go m.handleSignals() // start signal handler

	if m.opts.EnableSysAPI {
		go m.startSysAPI() // start sysAPI.
	}

	for name, sys := range m.systems {
		klog.V(4).Infof("synchronously starting subsystem %s", name)
		go sys.SyncStart()
	}

	// Async config watcher
	go m.watchConfig()

	// Start renewer for vault creds if we were able to make a renewer
	if m.secretRenewer != nil {
		go m.secretRenewer.Start()
		defer m.secretRenewer.Stop()

		for {
			select {
			case <-m.shutdown:
				return
			case done := <-m.secretRenewer.DoneCh():
				klog.Errorf("received error for vault credential renewals: %s", done)
			case renew := <-m.secretRenewer.RenewCh():
				klog.Infof("successfully renewed vault credentials at %s for %d seconds", renew.RenewedAt, renew.Secret.LeaseDuration)
			}
		}
	}

	<-m.shutdown
}

func (m *APIManager) GetVaultDbCreds() (*api.Secret, *api.LifetimeWatcher, error) {
	if m.vault != nil {
		secret, e := m.vault.Logical().Read(fmt.Sprintf("%s/creds/%s", m.config.GetString("vaultDbMountPath"), m.config.GetString("vaultDbRole")))
		if e != nil {
			return nil, nil, e
		}

		watcher, e := m.vault.NewLifetimeWatcher(&api.LifetimeWatcherInput{
			Secret:        secret,
			RenewBehavior: api.RenewBehaviorIgnoreErrors,
			Increment:     3600, // have them last an hour so they disappear quickly once the process dies.
		})
		if e != nil {
			return nil, nil, e
		}

		return secret, watcher, nil
	} else {
		return nil, nil, errors.New("vault not initialized, cannot retrieve credentials")
	}
}

// loads in configuration/secrets to override default values with the given application.
func (m *APIManager) initConfigs() {
	// Configure config file initialization first.
	m.config.AddConfigPath("/etc/" + m.registrar.AppName + "/config")
	m.config.AddConfigPath("test")
	m.config.SetConfigName("config")

	// Register default values into config map.
	for _, key := range m.ckeys {
		m.config.SetDefault(key.key, key.defaultVal)
	}

	if e := m.config.ReadInConfig(); e != nil {
		klog.Errorf("ALERT: unable to read in configuration file! Relying on system defaults. Error: %v", e)
	}

	// Register defualt values into secrets map.
	for _, key := range m.skeys {
		m.secrets.SetDefault(key.key, key.defaultVal)
	}

	m.secrets.AddConfigPath("/etc/" + m.registrar.AppName + "/secrets")
	m.secrets.AddConfigPath("test")
	m.secrets.SetConfigName("secrets")

	if e := m.secrets.ReadInConfig(); e != nil {
		klog.Errorf("ALERT: unable to read in secrets file! Relying on system defaults. Error: %v", e)
	}
}

func (m *APIManager) initVault() {
	var client *api.Client

	if m.vault != nil {
		client = m.vault
	} else {
		conf := api.DefaultConfig()
		conf.Address = m.config.GetString("vaultAddress")
		insecure := m.config.GetBool("vaultSslInsecure")

		if insecure {
			conf.ConfigureTLS(&api.TLSConfig{
				Insecure: true,
			})
		} else {
			conf.ConfigureTLS(&api.TLSConfig{
				CAPath:        m.config.GetString("vaultCAPath"),
				TLSServerName: m.config.GetString("vaultSNIName"),
				Insecure:      false,
			})
		}

		client, _ = api.NewClient(conf)
	}

	// Default to using vault token, otherwise try and login with k8s.
	vaultToken := os.Getenv("VAULT_TOKEN")

	if vaultToken == "" {
		klog.Info("no VAULT_TOKEN provided, attempting to login with k8s")
		opts := []k8sauth.LoginOption{k8sauth.WithMountPath(m.config.GetString("vaultK8sAuthMountPath"))}

		// Check if we have the serviceaccount token stored somewhere else, if so, add it to the options.
		saToken := os.Getenv("VAULT_SA_TOKEN")
		if saToken != "" {
			klog.Info("found VAULT_SA_TOKEN, using for k8s auth process")
			opts = append(opts, k8sauth.WithServiceAccountToken(saToken))
		}

		k8s, e := k8sauth.NewKubernetesAuth(m.config.GetString("vaultK8sRole"), opts...)
		if e != nil {
			klog.Errorf("unable to retrieve kubernetes auth credentials: %v", e)
		}
		secret, e := client.Auth().Login(context.Background(), k8s)
		if e != nil {
			klog.Errorf("unable to login with kubernetes auth: %v", e)
		}

		renewer, e := client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
			Secret:        secret,
			RenewBehavior: api.RenewBehaviorIgnoreErrors,
			Increment:     3600, // have them last an hour so they disappear quickly once the process dies.
		})

		if e == nil {
			klog.V(3).Info("created renewer for automatically renewing client credentials")
			m.secretRenewer = renewer
		} else {
			klog.Errorf("unable to create renewer for k8s auth: %s", e)
		}
	} else {
		klog.Info("found VAULT_TOKEN, using for vault auth")
		client.SetToken(vaultToken)
	}

	m.vault = client
}

// handleSignals does what you would think, it runs in a loop, blocking and waiting for incoming
// OS signals, and handling them.
func (m *APIManager) handleSignals() {
	for {
		sig := <-m.sigHandle
		switch sig {
		case syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT:
			{
				m.shutdownSubsystems()
				// m.shutdown <- 0 // Exit the sync start
				os.Exit(0)
			}
		case syscall.SIGHUP:
			{
				m.reloadSubsystems()
			}
		default:
			{
			}
		}
	}
}

func (api *APIManager) watchConfig() {
	api.config.WatchConfig()
	api.secrets.WatchConfig()
}

// Return all registered ConfigValues that are set up with this APIManager.
// These values should be READ ONLY!!! Please do not mutate any of these values after
// acquiring a reference to the slice.
func (api *APIManager) GetConfigValues() []*ConfigValue { return api.ckeys }

// Return all registered secret SecretValues that are set up with this APIManager.
// These values should be READ ONLY!!! Please do not mutate any of these values after
// acquiring a reference to the slice.
func (api *APIManager) GetSecretValues() []*SecretValue { return api.skeys }
