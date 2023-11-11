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
	"os"
	"sync"

	"github.com/fasthttp/router"
	"github.com/go-openapi/spec"
	vault "github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
)

// APIManager is the primary object that provides abstraction between process orchestration and startup and
// code that actually performs useful business logic. The easiest way to think about this object is it being
// PID 1 for the application, but WITHIN the actual app process. Some of the functionality that APIManager implements
// includes:
//
// - Registering signal handlers and then "forwarding" those signals to invidual subsystems of the process to handle.
//
//   - Deserializing configuration from multiple locations within the environment (primarily /etc/<app name>) and providing APIs
//     for subsystems to easily access those configuration parameters for their successful use. Some examples of this are
//     the database user/password that is managed by APIManager and accessed by the SQLManager subsystem for connecting to
//     the backing database. Another example could be configuring concurrency on the primary HTTP server. The server subsystem
//     retrieves the integer value through APIManager internal APIs, and is able to go on its merry way with serving up requests.
//
//   - Managing subsystems, and sending signals to subsystems whenever a "config reload" signal is sent to the process, or a
//     shutdown signal is sent to the process, things like the database driver need to gracefully end all transactions and cut off
//     TCP sockets to the DB. Same thing on the frontend, the webserver needs to close out all remaining sockets and requests.
//     For more information on the callbacks to manage subsystems that are invoked by the APIManager, please refer to the Subsystem
//     interface in types.go.
//
//   - Exposing a systems API for introspection into the inner workings of the app process. This includes exposing things such
//     as the swagger docs for the sysAPI itself, exposing all prometheus metrics, exposing readiness/liveness endpoints for the process
//     as a whole as well as individual subsystems, subsystem statistics, and much more. The idea long term is for this API to also
//     be used as a management interface for some kind of controller, that will automatically scale and configure app instances
//     as loads shifts for different API products.
type APIManager struct {
	m sync.RWMutex

	opts *APIManagerOpts

	registrar *SystemRegistrar

	// Count is the current count of subsystems that are registered with the manager.
	// Should not be edited after init.
	count uint
	// Map of all subsystem names mapped to their backing implementations.
	systems map[string]Subsystem

	vault *vault.Client

	// Config contains non-secret key/value data for configuring the process.
	config *viper.Viper

	ckeys []*ConfigValue
	skeys []*SecretValue

	// secret contains secrets credentials for configuring the process.
	// the most prevalent values within this container will be the database user/password.
	secrets *viper.Viper

	shutdown chan uint8

	// registry is the prometheus metrics registry that should have all process metrics registered
	// to it. This registry will then be collected when called upon by the SysAPI within APIManager.
	registry *prometheus.Registry

	// server contains the server data for the app sysAPI. This API is for operators to get detailed
	// insight into the performance, status, and overall health of the process and subsystems of app.
	// The long term goal is develop a
	server *fasthttp.Server

	// router contains the routing logic and handlers for all sysAPI operations and REST calls.
	router *router.Router

	// spec contains the swagger 2.0 docs for the sysAPI, this allows for programmatic access and
	// automated documentation for how the sysAPI is structured and different endpoints available to it.
	spec *spec.Swagger

	// sigHandle is registered to receive all signals for the process. APIManager logic will then
	// perform the task of forwarding signals to the subsystems as needed, and gracefully shutting
	// down the process.
	sigHandle chan os.Signal
}

// General purpose options for creating a new APIManager. These generally are options that
// will be set at compile time and shouldn't need to be worried about at runtime.
type APIManagerOpts struct {
	// Toggle whether the SysAPI should be advertised with this process.
	EnableSysAPI bool

	// Optionally provide a URL to a vault instance to interface with. If this string is
	// empty and the process needs to access vault, it will look for a URL in the VAULT_URL
	// envvar.
	VaultURL string
}

// Subsystem is a component of app that is bootstrapped by the manager upon process startup.
// Each subsystem needs to be registered with an init() method to the APIManager in order for
// its callbacks to be invoked at the proper times on process startup. For an example of a bare-bones
// NOP subsystem, please refer to the manager.DefaultSubsystem structure in defaultsubsystem.go.
type Subsystem interface {
	// All subsystems should be able to implement the collector interface for their metrics to be collected
	// and exported via the APIManager.
	prometheus.Collector

	// Return the name of this subsystem for referencing.
	Name() string

	// Optional hooks to be run before subsystem initialization.
	PreInit()
	PostInit()

	// Configs and Secrets return the ConfigValues and SecretValues for the subsystem for management
	// by the APIMAnager. This includes adding defaults to viper and populating configuration
	// defaulting commands.
	Configs() *[]*ConfigValue
	Secrets() *[]*SecretValue

	// Starts up this subsystem, if it returns an error, will try to reinitalize
	// the subsystem with backoff until an error is no longer returned.
	//
	// The waitgroup should be immediately deferred in the called function.
	Initialize(reg *SystemRegistrar) error

	// If a subsystem needs to be synchronously called by the manager (IE you will call a method
	// that should never return), wrap that method here and it will be held open by the manager
	// as long as the process is alive.
	SyncStart()

	// If the configuration is found to be changed, the manager will call this callback
	// to inform the subsystem to refresh itself given the new configuration changes.
	//
	// The waitgroup should be immediately deferred in the called function.
	Reload()

	// If a kill signal is sent to the process, the manager will inform all subsystems of shutdown
	// through this callback.
	//
	// The waitgroup should be immediately deferred in the called function.
	Shutdown()

	// This callback can be invoked at any point in execution by the manager to determine the
	// current status of the subsystem.
	Status() *SubsystemStatus
}
