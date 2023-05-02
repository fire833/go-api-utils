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
func New() *APIManager {
	return &APIManager{
		count:     0,
		systems:   make(map[string]Subsystem),
		shutdown:  make(chan uint8),
		keys:      []*ConfigKey{},
		config:    viper.New(),
		secrets:   viper.New(),
		registry:  prometheus.NewRegistry(),
		router:    router.New(),
		spec:      nil, // Start with null, the spec should be generated on Initialize().
		server:    nil, // Start with null, the server should be started on Initialize().
		sigHandle: make(chan os.Signal, 5),
	}
}

// Register subsystem should be called in an init function by all subsystems
// that want to be registered and run by the application at runtime.
// Please make sure that your subsystem defers WaitGroups provided to it
// IMMEDIATELY. If you do not do this, then the process will deadlock on startup
// and you will have a broken process, and then everyone has a bad day.
func (m *APIManager) registerSubsystem(sub Subsystem) {
	m.m.Lock()
	m.systems[sub.Name()] = sub
	m.registry.Register(sub)
	m.count++
	m.m.Unlock()

	if c := sub.Configs(); c != nil {
		m.registerConfigKey(*c...)
	}
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
func (m *APIManager) RegisterSysAPIHandler(method, path string, handler fasthttp.RequestHandler, swaggerdoc spec.PathItem, schemas ...*spec.Schema) error {
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

func (m *APIManager) registerConfigKey(keys ...*ConfigKey) {
	m.m.Lock()
	defer m.m.Unlock()

	for _, key := range keys {
		if key == nil {
			continue
		}

		if key.IsSecret {
			if key.DefaultVal != nil {
				m.secrets.SetDefault(key.Name, key.DefaultVal)
			}

			m.secretKeys = append(m.secretKeys, key)

		} else {
			if key.DefaultVal != nil {
				m.config.SetDefault(key.Name, key.DefaultVal)
			}

			m.keys = append(m.keys, key)
		}
	}
}

// loads in configuration/secrets to override default values with the given application.
func (m *APIManager) initConfigs() {
	// Register defaults first.
	for _, sub := range m.systems {
		m.registerConfigKey(*sub.Configs()...)
	}

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
	for _, sys := range registrar.Systems {
		m.registerSubsystem(sys)
	}

	// Tell the runtime to forward signals from the OS to this channel for downstream processing.
	signal.Notify(m.sigHandle)

	m.registrar = registrar

	// read in configuration and secrets before booting further, or at least attempt to.
	m.initConfigs()

	// Set up the sysAPI and all its handlers.
	m.initSysAPI()

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

// Return all registered ConfigKeys that are set up with this APIManager.
// These values should be READ ONLY!!! Please do not mutate any of these values after
// acquiring a reference to the slice.
func (api *APIManager) GetConfigKeys() []*ConfigKey { return api.keys }

// Return all registered secret ConfigKeys that are set up with this APIManager.
// These values should be READ ONLY!!! Please do not mutate any of these values after
// acquiring a reference to the slice.
func (api *APIManager) GetSecretConfigKeys() []*ConfigKey { return api.secretKeys }

// func (api *APIManager) GetSecretConfigKeys() []*ConfigKey { return api.secretKeys }