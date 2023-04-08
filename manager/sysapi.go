package manager

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/fire833/go-api-utils/serialization"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type ConfigInfo struct {
	Meta  *ConfigKey  `json:"meta" yaml:"meta" xml:"meta"`
	Value interface{} `json:"value" yaml:"value" xml:"value"`
}

var (
	sysAPIListenPort *ConfigKey = &ConfigKey{
		Name:        "sysAPIListenPort",
		Description: "Specify the listening port of the sysAPI. Should be an unsigned integer between 1 and 65535, but should be above 1024 preferably to avoid needing CAP_SYS_ADMIN or root privileges for the app process.",
		TypeOf:      Uint16,
		DefaultVal:  8081,
		IsSecret:    false,
	}

	sysAPIConcurrency *ConfigKey = &ConfigKey{
		Name:        "sysAPIConcurrency",
		Description: "Specify the amount of concurrent connections to be allowed to the SysAPI webserver concurrently.",
		IsSecret:    false,
		DefaultVal:  1000,
		TypeOf:      Int,
	}

	sysAPIReadBufferSize *ConfigKey = &ConfigKey{
		Name:        "sysAPIReadBufferSize",
		Description: "Specify per-connection buffer size for requests reading. This also limits the maximum header size. Increase this buffer if your clients send multi-KB RequestURIs and/or multi-KB headers (for example, BIG cookies).",
		TypeOf:      Int,
		DefaultVal:  4096,
		IsSecret:    false,
	}

	sysAPIWriteBufferSize *ConfigKey = &ConfigKey{
		Name:        "sysAPIWriteBufferSize",
		Description: "Per-connection buffer size for responses writing.",
		TypeOf:      Int,
		DefaultVal:  4096,
		IsSecret:    false,
	}

	sysAPIReadTimeout *ConfigKey = &ConfigKey{
		Name:        "sysAPIReadTimeout",
		Description: "ReadTimeout is the amount of time (in seconds) allowed to read the full request including body. The connection's read deadline is reset when the connection opens, or for keep-alive connections after the first byte has been read.",
		TypeOf:      Int,
		DefaultVal:  120,
		IsSecret:    false,
	}

	sysAPIWriteTimeout *ConfigKey = &ConfigKey{
		Name:        "sysAPIWriteTimeout",
		Description: "WriteTimeout is the maximum duration (in seconds) before timing out writes of the response. It is reset after the request handler has returned.",
		TypeOf:      Int,
		DefaultVal:  120,
		IsSecret:    false,
	}

	sysAPIIdleTimeout *ConfigKey = &ConfigKey{
		Name:        "sysAPIIdleTimeout",
		Description: "IdleTimeout is the maximum amount of time (in seconds) to wait for the next request when keep-alive is enabled.",
		TypeOf:      Int,
		DefaultVal:  120,
		IsSecret:    false,
	}
)

func init() {
	Manager.registerConfigKey(
		sysAPIListenPort,
		sysAPIConcurrency,
		sysAPIReadBufferSize,
		sysAPIWriteBufferSize,
		sysAPIReadTimeout,
		sysAPIWriteTimeout,
		sysAPIIdleTimeout,
	)
}

func (m *APIManager) initSysAPI() {
	ser := &fasthttp.Server{
		// overwrite the server name for a bit more obfuscation.
		Name: "null",

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

		Concurrency:     sysAPIConcurrency.RetriveValue().(int),
		ReadBufferSize:  sysAPIReadBufferSize.RetriveValue().(int),
		WriteBufferSize: sysAPIWriteBufferSize.RetriveValue().(int),
		ReadTimeout:     time.Duration(time.Second * time.Duration(sysAPIReadTimeout.RetriveValue().(int))),
		WriteTimeout:    time.Duration(time.Second * time.Duration(sysAPIWriteTimeout.RetriveValue().(int))),
		IdleTimeout:     time.Duration(time.Second * time.Duration(sysAPIIdleTimeout.RetriveValue().(int))),

		Handler: m.router.Handler,
	}

	// The primary spec object for sysAPI. Can have other stuff registered to it through RegisterSysAPIHandler().
	spec := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Consumes: []string{"application/json", "application/xml", "application/yaml", "application/API+Protobuf"},
			Produces: []string{"application/json", "application/xml", "application/yaml", "application/API+Protobuf"},
			Swagger:  "2.0",
			Info: &spec.Info{
				InfoProps: spec.InfoProps{
					Description: m.registrar.AppName + " Systems API.",
					Title:       m.registrar.AppName + " SysAPI",
					Version:     "0.0.4",
				},
			},
			Host:     "",
			BasePath: "/",
			Paths: &spec.Paths{
				Paths: map[string]spec.PathItem{
					"/metrics": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getMetrics").
								WithTags("sys").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns all available prometheus metrics for the system. This will be a text document in prometheus exporter format.")),
						},
					},
					"/swagger.json": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getSwagger").
								WithTags("sys").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns the current swagger specification file for sysAPI on the current app instance being queried. This will always be returned in JSON format.")),
						},
					},
					"/readyz": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getReady").
								WithTags("sys").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns the readiness of the app process, acts as a overarching readiness healthcheck for the system.").
									WithSchema(spec.RefSchema("#/definitions/OKResponse"))),
						},
					},
					"/livez": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getLive").
								WithTags("sys").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns the liveness of the app process, acts as a overarching liveness healthcheck for the system.").
									WithSchema(spec.RefSchema("#/definitions/OKResponse"))),
						},
					},
					"/configuration": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getConfiguration").
								WithProduces("application/json", "application/yaml", "application/xml").
								WithTags("sys").
								WithDescription("Returns the current configuration parameters of the app process.").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns the current configuration parameters of the app process.").
									WithSchema(spec.ArrayProperty(spec.RefSchema("#/definitions/ConfigKeyValue")))),
						},
					},
					"/secrets": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getSecrets").
								WithProduces("application/json", "application/yaml", "application/xml").
								WithTags("sys").
								WithDescription("Returns the current secret keys configured with the app process. Please note the actual values will be hidden.").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns the current secret keys configured with the app process. Please note the actual values will be hidden.").
									WithSchema(spec.ArrayProperty(spec.RefSchema("#/definitions/ConfigKeyValue")))),
						},
					},
					"/status": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getStatus").
								WithTags("sys").
								WithDescription("Return the status of all subsystems within the running app instance.").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Return the status of all subsystems within the running app instance.").
									WithSchema(spec.ArrayProperty(spec.RefSchema("#/definitions/SubsystemStatus")))),
						},
					},
					"/status/{SUBSYSTEM}": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getStatusSubsystem").
								WithTags("sys").
								AddParam(spec.PathParam("SUBSYSTEM").Typed("string", "")).
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns the status of a subsystem in the current process.").
									WithSchema(spec.RefSchema("#/definitions/SubsystemStatus"))).
								RespondsWith(400, spec.NewResponse().
									WithDescription("Returns that the subsystem was not found on the process.").
									WithSchema(spec.RefSchema("#/definitions/GenericErrorResponse"))),
						},
					},
					"/buildinfo": {
						PathItemProps: spec.PathItemProps{
							Get: spec.NewOperation("getBuildInfo").
								WithTags("sys").
								RespondsWith(200, spec.NewResponse().
									WithDescription("Returns an object with the current build info.").
									WithSchema(spec.RefSchema("#/definitions/BuildInfo"))),
						},
					},
				},
			},
			Definitions: spec.Definitions{
				"SubsystemStatus":      *subsystemStatusSchema,
				"BuildInfo":            *buildInfoSchema,
				"OKResponse":           *serialization.OKResponseSchema,
				"GenericErrorResponse": *serialization.GenericErrorResponseSchema,
				"ConfigKeyValue":       *configKeySchema,
				"ConfigValueType":      *configValueType,
			},
		},
	}

	// Load intial collectors to the registry subsystem.
	m.registry.Register(collectors.NewBuildInfoCollector())
	m.registry.Register(collectors.NewGoCollector())

	m.router.NotFound = func(ctx *fasthttp.RequestCtx) {
		serialization.GenericNotFoundResponseHandler(ctx)
	}

	m.router.MethodNotAllowed = func(ctx *fasthttp.RequestCtx) {
		serialization.GenericMethodNotAllowedResponseHandler(ctx)
	}

	m.router.PanicHandler = func(rc *fasthttp.RequestCtx, i interface{}) {
		serialization.GenericInternalErrorResponseHandler(rc)
	}

	// Get HTTP handler for this registry, and register it with the sysapi http server.
	m.router.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(
		promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})))

	m.router.GET("/readyz", func(ctx *fasthttp.RequestCtx) {
		serialization.GenericOKResponseHandler(ctx)
	})

	m.router.GET("/status", func(ctx *fasthttp.RequestCtx) {
		statuses := []*SubsystemStatus{}
		for _, sys := range m.systems {
			statuses = append(statuses, sys.Status())
		}

		serialization.MarshalBodyByAcceptHeader(ctx, &SubsystemStatusList{Items: statuses})
	})

	m.router.GET("/status/{SUBSYSTEM}", func(ctx *fasthttp.RequestCtx) {
		sub := ctx.UserValue("SUBSYSTEM").(string)
		if val, ok := m.systems[sub]; ok {
			serialization.MarshalBodyByAcceptHeader(ctx, val.Status())
			return
		}

		serialization.BadRequestResponseHandler(ctx, "subsystem not found in process")
	})

	m.router.GET("/livez", func(ctx *fasthttp.RequestCtx) {
		serialization.GenericOKResponseHandler(ctx)
	})

	m.router.GET("/swagger.json", func(ctx *fasthttp.RequestCtx) {
		data, e := json.MarshalIndent(m.spec, "", "   ")
		if e != nil {
			serialization.InternalErrorResponseHandler(ctx, e.Error())
			return
		}

		ctx.Response.SetBody(data)
		ctx.Response.SetStatusCode(http.StatusOK)
	})

	m.router.GET("/configuration", func(ctx *fasthttp.RequestCtx) {
		var values []ConfigInfo

		for _, key := range m.keys {
			values = append(values, ConfigInfo{
				Meta:  key,
				Value: key.RetriveValue(),
			})
		}

		data, e := json.MarshalIndent(values, "", "   ")
		if e != nil {
			serialization.InternalErrorResponseHandler(ctx, e.Error())
		}

		ctx.Response.SetBody(data)
		ctx.Response.SetStatusCode(http.StatusOK)
	})

	m.router.GET("/secrets", func(ctx *fasthttp.RequestCtx) {
		var values []ConfigInfo

		for _, key := range m.secretKeys {
			values = append(values, ConfigInfo{
				Meta:  key,
				Value: "*****",
			})
		}

		data, e := json.MarshalIndent(values, "", "   ")
		if e != nil {
			serialization.InternalErrorResponseHandler(ctx, e.Error())
		}

		ctx.Response.SetBody(data)
		ctx.Response.SetStatusCode(http.StatusOK)
	})

	m.router.GET("/buildinfo", func(ctx *fasthttp.RequestCtx) {
		bi := &BuildInfo{
			// Version:   pkg.Version,
			// Commit:    pkg.Commit,
			// BuildTime: pkg.BuildTime,
			// Os:        pkg.Os,
			// Arch:      pkg.Arch,
		}

		serialization.MarshalBodyByAcceptHeader(ctx, bi)
	})

	m.server = ser
	m.spec = spec
}

func (m *APIManager) startSysAPI() {
	m.server.ListenAndServe(":" + strconv.Itoa(int(sysAPIListenPort.RetriveValue().(uint16))))
}
