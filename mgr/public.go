package mgr

import (
	"errors"
	"fmt"
	"os"

	"github.com/fasthttp/router"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
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
