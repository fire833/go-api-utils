package apiserver

import (
	"fmt"
	"time"

	"github.com/fasthttp/router"
	manager "github.com/fire833/go-api-utils/mgr"
	"github.com/fire833/go-api-utils/serialization"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/valyala/fasthttp"
	"k8s.io/klog/v2"
)

var SERVER *APIServer

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

func New() *APIServer {
	return &APIServer{
		server:     nil,
		router:     nil,
		servername: "unknown",
	}
}

func (s *APIServer) Name() string { return "api" }

func (s *APIServer) Initialize(reg *manager.SystemRegistrar) error {
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
	klog.V(5).Info("api: registering api endpoints to router")
	reg.Registration.RegisterEndpoints(apiServerPrefix.GetString(), s.router)

	klog.V(5).Info("api: initializing fasthttp server")
	s.server = &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			defer s.requestCount.Inc()
			// TODO add auth middleware here to be done before traversing the route tree.
			s.router.Handler(ctx)
		},

		// overwrite the server name for a bit more obfuscation.
		Name:        "null",
		Concurrency: int(apiServerConcurrency.GetUint()),

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

		ReadBufferSize:  int(apiServerReadBufferSize.GetUint()),
		WriteBufferSize: int(apiServerWriteBufferSize.GetUint()),
		ReadTimeout:     time.Duration(time.Second * time.Duration(apiServerReadTimeout.GetUint())),
		WriteTimeout:    time.Duration(time.Second * time.Duration(apiServerWriteTimeout.GetUint())),
		IdleTimeout:     time.Duration(time.Second * time.Duration(apiServerIdleTimeout.GetUint())),
	}

	// For now, we don't need the swagger specification to be in memory with the process,
	// that's just extra overhead that currently won't do anything.
	// s.spec2, _ = app.GenerateSpec(&app.SwaggerGenOpts{})

	s.IsInitialized = true

	return nil
}

func (s *APIServer) SyncStart() {
	klog.V(2).Infof("serving apiserver on %s:%d", apiServerListenIp.GetString(), apiServerListenPort.GetUint16())
	if e := s.server.ListenAndServe(fmt.Sprintf("%s:%d", apiServerListenIp.GetString(), apiServerListenPort.GetUint16())); e != nil {
		klog.Errorf("unable to start api: %s", e.Error())
		// os.Exit(1) // TODO perhaps make a better exit strategy here.
	}
}

func (s *APIServer) Shutdown() {
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
