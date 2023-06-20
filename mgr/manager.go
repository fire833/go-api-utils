/*
*	Copyright (C) 2023 Kendall Tauser
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
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fasthttp/router"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"k8s.io/klog/v2"
)

// Return a new, uninitialized APIManager object. Please do not allocate more than one of
// these objects per process, they are desined to the the PID 1 of a process, and manage the
// routines of subsystems that perform actual business logic.
func New(opts *APIManagerOpts) *APIManager {
	m := &APIManager{
		count:     0,
		opts:      opts,
		systems:   make(map[string]Subsystem),
		shutdown:  make(chan uint8),
		config:    viper.New(),
		secrets:   viper.New(),
		registry:  prometheus.NewRegistry(),
		vault:     nil,
		router:    router.New(),
		spec:      nil, // Start with null, the spec should be generated on Initialize().
		server:    nil, // Start with null, the server should be started on Initialize().
		sigHandle: make(chan os.Signal, 5),
	}

	if mgr == nil {
		mgr = m
	} else {
		panic("manager already initialized")
	}

	return m
}

// For subsystems which need to have viewports/logic exposed through SysAPI,
// subsystems can call this method in order to register their handler functions
// to be called by SysAPI. This can include things such as toggling tracing
// with the otel subsystem depending on type, exposing control operations/tweaking
// capabilites with controllers, etc.
//
// PLEASE USE THIS WITH CAUTION!!!
// This method could inadvertently expose logic that should not be callable from
// outside sources, which obviously can be a great security risk. So, please be
// careful with what handlers you are exposing with this method.
// You have been warned.
//
// Also, please ensure you are providing ACCURATE swagger specification objects
// along with this handler registration request, as these will be integrated with the server's
// swagger spec. These objects will then be served by /swagger.json, and it will make integration
// MUCH easier if the actual behavior of the endpoint is reflected in the swagger documentation.
func RegisterSysAPIHandler(method, path string, handler fasthttp.RequestHandler, swaggerdoc spec.PathItem, schemas ...*spec.Schema) error {
	m := mgr
	if m == nil {
		return errors.New("global APIManager not initialized")
	}

	m.m.Lock()
	defer m.m.Unlock()

	m.router.Handle(method, path, handler)
	if err := recover(); err != nil {
		return fmt.Errorf("handler registration error: %v", err)
	}

	if _, exists := m.spec.Paths.Paths[path]; exists {
		return errors.New("path already registered with SysAPI server")
	}

	m.spec.Paths.Paths[path] = swaggerdoc

	for _, schema := range schemas {
		if _, exists := m.spec.Definitions[schema.Title]; exists {
			return errors.New("schema exists already within sysAPI spec")
		} else {
			m.spec.Definitions[schema.Title] = *schema
		}
	}

	return nil
}

// loads in configuration/secrets to override default values with the given application.
func (m *APIManager) initConfigs() {
	// Configure config file initialization first.
	m.config.AddConfigPath("/etc/" + m.registrar.AppName)
	m.config.AddConfigPath("test")
	m.config.SetConfigName("config")

	if e := m.config.ReadInConfig(); e != nil {
		klog.Errorf("ALERT: unable to read in configuration file! Relying on system defaults. Error: %v", e)
	}

	m.secrets.AddConfigPath("/etc/" + m.registrar.AppName + "/secrets")
	m.secrets.AddConfigPath("test")
	m.secrets.SetConfigName("secrets")

	if e := m.secrets.ReadInConfig(); e != nil {
		klog.Errorf("ALERT: unable to read in secrets file! Relying on system defaults. Error: %v", e)
	}
}

// Global method used for initializing an APIManager instance. This includes registering
// signal handlers, creating SysAPI, and starting up all the required subsystems as according to the
// provided registrar.
func (m *APIManager) Initialize(registrar *SystemRegistrar) {
	// Tell the runtime to forward signals from the OS to this channel for downstream processing.
	signal.Notify(m.sigHandle)

	m.registrar = registrar

	for name, sys := range m.systems {
		configs := *sys.Configs()
		klog.V(5).Infof("registering %d config keys for subsystem %s", len(configs), name)

		m.ckeys = append(m.ckeys, configs...)

		secrets := *sys.Secrets()
		klog.V(5).Infof("registering %d secret keys for subsystem %s", len(secrets), name)

		m.skeys = append(m.skeys, secrets...)
	}

	// read in configuration and secrets before booting further, or at least attempt to.
	m.initConfigs()

	// Set up the sysAPI and all its handlers.
	if m.opts.EnableSysAPI {
		m.initSysAPI()
	}

	m.initializeSubsystems(registrar)

	m.setGlobals()
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

// Return all registered ConfigValues that are set up with this APIManager.
// These values should be READ ONLY!!! Please do not mutate any of these values after
// acquiring a reference to the slice.
func (api *APIManager) GetConfigValues() []*ConfigValue { return api.ckeys }

// Return all registered secret SecretValues that are set up with this APIManager.
// These values should be READ ONLY!!! Please do not mutate any of these values after
// acquiring a reference to the slice.
func (api *APIManager) GetSecretValues() []*SecretValue { return api.skeys }
