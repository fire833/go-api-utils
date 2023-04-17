package apiserver

import (
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/fasthttp/router"
	"github.com/fire833/go-api-utils/manager"
	"github.com/fire833/go-api-utils/serialization"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"k8s.io/klog/v2"
)

var (
	SERVER *APIServer
)

// APIServer is the primary object by which content is served RESTfully to clients on the internet.
// This object will have all Object managers and/or APIGroups registered to it that the operator
// wishes to be served within this individual app process. Given that app is supposed to be a
// distributed API, with APIGroups and Objects to be served across multiple processes and being
// multiplexed by a Layer 7 load balander in front of it, the registration of objects to this server
// vary based on the specific app binary you are building and could even be non-static at runtime.
//
// Otherwise, the api is another subsystem that is registered to the APIManager, and gets started
// up in the correct order along with the other subsystems of the process (ie SQL driver, elastic client,
// SYS API, Auth API client, etc.)
type APIServer struct {
	manager.DefaultSubsystem

	servername string

	// server contains state for serving requests over HTTP to clients on the internet.
	server *fasthttp.Server

	// router contains the route state for this api.
	router *router.Router

	// spec contains the generated Swagger 2.0 spec that will be implemented by this api.
	// spec2 *spec.Swagger

	// spec contains the generated OpenAPI 3.0 spec that will be implemented by this api.
	// spec3 openapi.OpenAPI3Spec
	requestCount prometheus.Counter
}

func (s *APIServer) Name() string { return "api" }

func (s *APIServer) Initialize(wg *sync.WaitGroup, reg *manager.SystemRegistrar) error {
	defer wg.Done()

	s.servername = reg.AppName

	s.router = router.New()

	s.router.NotFound = func(ctx *fasthttp.RequestCtx) {
		serialization.GenericNotFoundResponseHandler(ctx)
	}

	s.router.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		serialization.GenericMethodNotAllowedResponseHandler(ctx)
	}

	s.router.PanicHandler = func(rc *fasthttp.RequestCtx, i interface{}) {
		// Since we should really NEVER panic inside of a handler (unless we have a bug),
		// log this error as a global error.
		klog.Errorf("request handler panicked: %s - %v", rc.Request.RequestURI(), i)

		serialization.GenericInternalErrorResponseHandler(rc)
	}

	s.requestCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: reg.AppName,
		Subsystem: s.Name(),
		Name:      "total_requests",
		Help:      "Metrics on the total number of requests made to this instance, successful or otherwise",
	})

	// Register from the global object.
	reg.Registration.RegisterEndpoints(apiServerPrefix.RetriveValue().(string), s.router)

	s.server = &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			defer s.requestCount.Inc()
			// TODO add auth middleware here to be done before traversing the route tree.
			s.router.Handler(ctx)
		},

		// overwrite the server name for a bit more obfuscation.
		Name:        "null",
		Concurrency: apiServerConcurrency.RetriveValue().(int),

		// Just enable this to always true, we shouldn't ever need information leaked.
		SecureErrorLogMessage: true,

		// Even better, just disable the server header entirely.
		NoDefaultServerHeader: true,

		// Set handler for dealing with lower level failures.
		ErrorHandler: func(ctx *fasthttp.RequestCtx, e error) {
			// May want to update this in the future with more fine-grained checks,
			// But return a not-acceptable error if there is a low level error with
			// reading in a request/response.
			serialization.NotAcceptableResponseHandler(ctx, e.Error())
		},

		ReadBufferSize:  apiServerReadBufferSize.RetriveValue().(int),
		WriteBufferSize: apiServerWriteBufferSize.RetriveValue().(int),
		ReadTimeout:     time.Duration(time.Second * time.Duration(apiServerReadTimeout.RetriveValue().(int))),
		WriteTimeout:    time.Duration(time.Second * time.Duration(apiServerWriteTimeout.RetriveValue().(int))),
		IdleTimeout:     time.Duration(time.Second * time.Duration(apiServerIdleTimeout.RetriveValue().(int))),
	}

	// For now, we don't need the swagger specification to be in memory with the process,
	// that's just extra overhead that currently won't do anything.
	// s.spec2, _ = app.GenerateSpec(&app.SwaggerGenOpts{})

	s.IsInitialized = true

	return nil
}

func (s *APIServer) SyncStart() {
	if e := s.server.ListenAndServe(apiServerListenIp.RetriveValue().(string) + ":" + strconv.Itoa(int(apiServerListenPort.RetriveValue().(uint16)))); e != nil {
		klog.Errorf("unable to start api: %s", e.Error())
		os.Exit(1) // TODO perhaps make a better exit strategy here.
	}
}

func (s *APIServer) Reload(wg *sync.WaitGroup) { wg.Done() }

func (s *APIServer) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
	if e := s.server.Shutdown(); e != nil {
		klog.Errorf("unable to gracefully shutdown apiserver subsystem: %v", e)
	}

	s.IsShutdown = true
}

func (s *APIServer) Status() *manager.SubsystemStatus {
	return &manager.SubsystemStatus{
		Name:          s.Name(),
		IsInitialized: s.IsInitialized,
		IsShutdown:    s.IsShutdown,
	}
}

func (s *APIServer) Collect(ch chan<- prometheus.Metric) {
	if s.server != nil {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(s.servername, s.Name(), "open_connections"),
				"Metrics on open connections to main HTTP server",
				nil,
				prometheus.Labels{},
			),
			prometheus.GaugeValue,
			float64(s.server.GetOpenConnectionsCount()))

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(s.servername, s.Name(), "served_connections"),
				"Metrics on currently served connections to main HTTP server",
				nil,
				prometheus.Labels{},
			),
			prometheus.GaugeValue,
			float64(s.server.GetCurrentConcurrency()))

		ch <- s.requestCount
	}
}
