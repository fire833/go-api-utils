package manager

import (
	"os"
	"sync"

	"github.com/fasthttp/router"
	"github.com/go-openapi/spec"
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
// - Deserializing configuration from multiple locations within the environment (primarily /etc/<app name>) and providing APIs
//   for subsystems to easily access those configuration parameters for their successful use. Some examples of this are
//   the database user/password that is managed by APIManager and accessed by the SQLManager subsystem for connecting to
//   the backing database. Another example could be configuring concurrency on the primary HTTP server. The server subsystem
//   retrieves the integer value through APIManager internal APIs, and is able to go on its merry way with serving up requests.
//
// - Managing subsystems, and sending signals to subsystems whenever a "config reload" signal is sent to the process, or a
//   shutdown signal is sent to the process, things like the database driver need to gracefully end all transactions and cut off
//   TCP sockets to the DB. Same thing on the frontend, the webserver needs to close out all remaining sockets and requests.
//   For more information on the callbacks to manage subsystems that are invoked by the APIManager, please refer to the Subsystem
//   interface in types.go.
//
// - Exposing a systems API for introspection into the inner workings of the app process. This includes exposing things such
//   as the swagger docs for the sysAPI itself, exposing all prometheus metrics, exposing readiness/liveness endpoints for the process
//   as a whole as well as individual subsystems, subsystem statistics, and much more. The idea long term is for this API to also
//   be used as a management interface for some kind of controller, that will automatically scale and configure app instances
//   as loads shifts for different API products.
//
type APIManager struct {
	m sync.RWMutex

	registrar *SystemRegistrar

	// Count is the current count of subsystems that are registered with the manager.
	// Should not be edited after init.
	count uint
	// Map of all subsystem names mapped to their backing implementations.
	systems map[string]Subsystem

	// Set of keys that are to be managed by APIManager, and will be called upon by one
	// or more subsystems that will depend on that value.
	keys []*ConfigKey

	// Set of Keys that are secret, same as keys field.
	secretKeys []*ConfigKey

	// Config contains non-secret key/value data for configuring the process.
	config *viper.Viper

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

	// Return the set of config keys that should be tracked/registered with the APImanager.
	// This should return a copy of the reference slice of config keys, since they will need to be
	// called by the registering subsystem at initialization/runtime.
	Configs() *[]*ConfigKey

	// Optionally, can set a global variable to point to the initialized subsystem object.
	// This can be useful for global functions to point to a central state location, and then check and
	// fail out if the pointer is nil, otherwise can be followed to where the initialized and running
	// subsystem is located.
	SetGlobal()

	// Starts up this subsystem, if it returns an error, will try to reinitalize
	// the subsystem with backoff until an error is no longer returned.
	//
	// The waitgroup should be immediately deferred in the called function.
	Initialize(wg *sync.WaitGroup, reg AppRegistration) error

	// If a subsystem needs to be synchronously called by the manager (IE you will call a method
	// that should never return), wrap that method here and it will be held open by the manager
	// as long as the process is alive.
	SyncStart()

	// If the configuration is found to be changed, the manager will call this callback
	// to inform the subsystem to refresh itself given the new configuration changes.
	//
	// The waitgroup should be immediately deferred in the called function.
	Reload(wg *sync.WaitGroup)

	// If a kill signal is sent to the process, the manager will inform all subsystems of shutdown
	// through this callback.
	//
	// The waitgroup should be immediately deferred in the called function.
	Shutdown(wg *sync.WaitGroup)

	// This callback can be invoked at any point in execution by the manager to determine the
	// current status of the subsystem.
	Status() *SubsystemStatus
}
