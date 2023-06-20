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

package otel

import (
	"context"
	"errors"
	"os"
	"strings"
	"sync"

	manager "github.com/fire833/go-api-utils/mgr"
	"github.com/go-openapi/spec"
	"github.com/valyala/fasthttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"k8s.io/klog/v2"
)

const (
	OTelManagerSubsystem string = "otel"
	SamplerDesc          string = "appSampler"
)

var OTEL *OTELManager

type OTELManager struct {
	manager.DefaultSubsystem
	sdktrace.Sampler

	m sync.RWMutex

	// Hashmap of all possible combinations of API group, Object, and trace operation associated with that object,
	// and the corresponding boolean on whether that trace should be recorded or dropped.
	sampleToggle map[string]bool
	sampleLock   sync.Mutex // Lock for sampleToggle map.

	exporter trace.SpanExporter

	tracer *sdktrace.TracerProvider
}

func (o *OTELManager) Name() string { return OTelManagerSubsystem }

func (o *OTELManager) Initialize(wg *sync.WaitGroup, reg manager.AppRegistration) error {
	defer wg.Done()
	o.m.Lock()
	defer o.m.Unlock()

	switch otelExporterType.GetString() {
	case "otlp":
		{
			return errors.New("otlp tracing exposition format not implemented")
			// if exp, e := otlptrace.New(context.Background(), otlp); e != nil {
			// 	return e
			// } else {
			// 	o.exporter = exp
			// }
		}
	case "otlphttp", "http", "h":
		{
			if exp, e := otlptracehttp.New(
				context.Background(),
				otlptracehttp.WithEndpoint(otelHTTPExportHost.GetString()),
				otlptracehttp.WithURLPath(otelHTTPExportPath.GetString()),
				otlptracehttp.WithInsecure(), // TODO: refactor this so insecure can be configured.
			); e != nil {
				return e
			} else {
				o.exporter = exp
			}
		}
	case "jaeger", "j":
		{
			if exp, e := jaeger.New(
				jaeger.WithAgentEndpoint(
					jaeger.WithAgentHost(jaegerExportHost.GetString()),
					jaeger.WithAgentPort(jaegerExportPort.GetString()),
				),
			); e != nil {
				return e
			} else {
				o.exporter = exp
			}
		}
		// Default to exporting with prometheus.
	default:
		{
			if exp, e := stdouttrace.New(
				stdouttrace.WithWriter(os.Stdout),
				stdouttrace.WithPrettyPrint(),
			); e != nil {
				return e
			} else {
				o.exporter = exp
			}
		}
	}

	// Get resource for this process and the metrics it exports.
	r, e1 := resource.New(
		context.Background(),
		resource.WithContainer(),
		resource.WithProcess(),
		resource.WithOSType(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("app"),
			// semconv.ServiceVersionKey.String(pkg.Version),
		),
	)
	if e1 != nil {
		return e1
	}

	// Create the trace provider for global consumption.
	o.tracer = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(o.exporter),
		sdktrace.WithResource(r),
		sdktrace.WithSampler(o),
		sdktrace.WithRawSpanLimits(sdktrace.SpanLimits{}),
	)

	otel.SetTracerProvider(o.tracer)
	otel.SetLogger(klog.NewKlogr())

	// set up the sampler code, including SysAPI routes for modifying which operations will be traced.
	o.sampleToggle = make(map[string]bool)

	val := otelDefaultTracingStatus.GetBool()

	for _, trace := range reg.RegisterOTELTraces() {
		o.sampleToggle[trace] = val
	}

	manager.RegisterSysAPIHandler(fasthttp.MethodGet, "/trace/status", o.statusHandler, spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Get: spec.NewOperation("traceStatus").WithDescription("Get the status of all operations enabled to be traced within the running instance.").
				WithTags("sys", "otel").
				RespondsWith(200, spec.NewResponse().WithDescription("List of all operations and whether they are enabled.").WithSchema(spec.RefSchema("#/definitions/SamplerStatusList"))),
		},
	}, samplerStatus)

	manager.RegisterSysAPIHandler(fasthttp.MethodPut, "/trace/enable/{NAME}", o.enable, spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Put: spec.NewOperation("enableTrace").WithDescription("Enable one or all traces for all internal Objects being served by this instance.").
				WithTags("sys", "otel").
				AddParam(spec.PathParam("NAME").Typed("string", "").WithDescription("Specify the object name you wish to enable. NOTE: if you specify 'all' as the ID value, then all traces that have not already been enabled within the instance will be enabled. For a reference of all possible names, call /trace/status for a full list and current status.")).
				RespondsWith(200, spec.NewResponse().WithSchema(spec.RefSchema("#/definitions/OKResponse"))).
				RespondsWith(404, spec.NewResponse().WithSchema(spec.RefSchema("#/definitions/GenericErrorResponse")).WithDescription("Returns a 404 if the named object is not found on the current instance.")),
		},
	})

	manager.RegisterSysAPIHandler(fasthttp.MethodPut, "/trace/disable/{NAME}", o.disable, spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Put: spec.NewOperation("disableTrace").WithDescription("Disable one or all traces for all internal Objects being served by this instance.").
				WithTags("sys", "otel").
				AddParam(spec.PathParam("NAME").Typed("string", "").WithDescription("Specify the object name you wish to disable. NOTE: if you specify 'all' as the ID value, then all traces that have not already been disabled within the instance will be disabled. For a reference of all possible names, call /trace/status for a full list and current status.")).
				RespondsWith(200, spec.NewResponse().WithSchema(spec.RefSchema("#/definitions/OKResponse"))).
				RespondsWith(404, spec.NewResponse().WithSchema(spec.RefSchema("#/definitions/GenericErrorResponse")).WithDescription("Returns a 404 if the named object is not found on the current instance.")),
		},
	})

	o.IsInitialized = true

	return nil
}

// Free all resources from the exporter and shutdown.
func (o *OTELManager) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()

	klog.V(4).Infoln("otel: shutting down tracer")
	if e := o.tracer.Shutdown(context.Background()); e != nil {
		klog.Errorf("unable to gracefully shutdown otel tracer subsystem: %v", e)
	}

	klog.V(4).Infoln("otel: shutting down span exporter")
	if e := o.exporter.Shutdown(context.Background()); e != nil {
		klog.Errorf("unable to gracefully shutdown otel exporter subsystem: %v", e)
	}

	o.IsShutdown = true
}

func (e *OTELManager) Status() *manager.SubsystemStatus {
	return &manager.SubsystemStatus{
		Name:          e.Name(),
		IsInitialized: e.IsInitialized,
		IsShutdown:    e.IsShutdown,
		Meta:          nil,
	}
}

// Implementation for using the manager as a sampler for the global trace provider.
func (o *OTELManager) Description() string { return SamplerDesc }

// ShouldSample returns a SamplingResult based on a decision made from the passed parameters.
func (o *OTELManager) ShouldSample(parameters sdktrace.SamplingParameters) sdktrace.SamplingResult {
	// TODO: we may be able to avoid doing this split in the hot path, need to look into this in the future.
	params := strings.SplitN(parameters.Name, "/", 2) // For now, we just want to filter based on object,
	// we don't need to worry about filtering on a per-operation basis just yet.

	// Look up this specific operation within the hash table, trace if the operation is set to true.
	if yes, ok := o.sampleToggle[params[0]]; yes && ok {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.RecordAndSample,
			Tracestate: oteltrace.SpanContextFromContext(parameters.ParentContext).TraceState(),
		}
	}

	// Default to dropping everything.
	return sdktrace.SamplingResult{
		Decision:   sdktrace.Drop,
		Tracestate: oteltrace.SpanContextFromContext(parameters.ParentContext).TraceState(),
	}
}
