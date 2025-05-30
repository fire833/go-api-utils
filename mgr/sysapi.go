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
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/fire833/go-api-utils/serialization"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"k8s.io/klog/v2"
)

type ConfigInfo struct {
	Meta  *ConfigValue `json:"meta" yaml:"meta" xml:"meta"`
	Value interface{}  `json:"value" yaml:"value" xml:"value"`
}

type SecretInfo struct {
	Meta  *SecretValue `json:"meta" yaml:"meta" xml:"meta"`
	Value interface{}  `json:"value" yaml:"value" xml:"value"`
}

var (
	sysAPIListenAddress *ConfigValue = NewConfigValue(
		"sysAPIListenAddress",
		"Specify the listening address of sysAPI.",
		"0.0.0.0",
	)

	sysAPIListenPort *ConfigValue = NewConfigValue(
		"sysAPIListenPort",
		"Specify the listening port of sysAPI. Should be an unsigned integer between 1 and 65535, but should be above 1024 preferably to avoid needing CAP_SYS_ADMIN or root privileges for the app process.",
		uint16(8081),
	)

	sysAPIConcurrency *ConfigValue = NewConfigValue(
		"sysAPIConcurrency",
		"Specify the amount of concurrent connections to be allowed to the SysAPI webserver concurrently.",
		uint(1000),
	)

	sysAPIReadBufferSize *ConfigValue = NewConfigValue(
		"sysAPIReadBufferSize",
		"Specify per-connection buffer size for requests reading. This also limits the maximum header size. Increase this buffer if your clients send multi-KB RequestURIs and/or multi-KB headers (for example, BIG cookies).",
		uint(4096),
	)

	sysAPIWriteBufferSize *ConfigValue = NewConfigValue(
		"sysAPIWriteBufferSize",
		"Per-connection buffer size for responses writing.",
		uint(4096),
	)

	sysAPIReadTimeout *ConfigValue = NewConfigValue(
		"sysAPIReadTimeout",
		"ReadTimeout is the amount of time (in seconds) allowed to read the full request including body. The connection's read deadline is reset when the connection opens, or for keep-alive connections after the first byte has been read.",
		uint(120),
	)

	sysAPIWriteTimeout *ConfigValue = NewConfigValue(
		"sysAPIWriteTimeout",
		"WriteTimeout is the maximum duration (in seconds) before timing out writes of the response. It is reset after the request handler has returned.",
		uint(120),
	)

	sysAPIIdleTimeout *ConfigValue = NewConfigValue(
		"sysAPIIdleTimeout",
		"IdleTimeout is the maximum amount of time (in seconds) to wait for the next request when keep-alive is enabled.",
		uint(120),
	)
)

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

		Concurrency:     int(sysAPIConcurrency.GetUint()),
		ReadBufferSize:  int(sysAPIReadBufferSize.GetUint()),
		WriteBufferSize: int(sysAPIWriteBufferSize.GetUint()),
		ReadTimeout:     time.Duration(time.Second * time.Duration(sysAPIReadTimeout.GetUint())),
		WriteTimeout:    time.Duration(time.Second * time.Duration(sysAPIWriteTimeout.GetUint())),
		IdleTimeout:     time.Duration(time.Second * time.Duration(sysAPIIdleTimeout.GetUint())),

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

		for _, key := range m.ckeys {
			k := key.Get()
			values = append(values, ConfigInfo{
				Meta:  key,
				Value: k,
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
		var values []SecretInfo

		for _, key := range m.skeys {
			values = append(values, SecretInfo{
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
	bind := fmt.Sprintf("%s:%d", sysAPIListenAddress.GetString(), sysAPIListenPort.GetUint16())
	klog.V(5).Infof("starting sysAPI on %s", bind)
	m.server.ListenAndServe(bind)
}
