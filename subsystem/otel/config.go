package otel

import "github.com/fire833/go-api-utils/manager"

var (
	otelDefaultTracingStatus *manager.ConfigKey = &manager.ConfigKey{
		Name:        "otelDefaultTracingStatus",
		Description: "Toggle whether or not tracing is default enabled for all types within the running instance or not. This will be the initial boolean value applied to the internal hashmap.",
		DefaultVal:  false,
		IsSecret:    false,
		TypeOf:      manager.Bool,
	}

	otelExporterType *manager.ConfigKey = &manager.ConfigKey{
		Name:        "otelExporterType",
		Description: "Specify the span exporter you wish to use to export traces from opentelemetry tracing within this app instance.",
		DefaultVal:  "otlphttp",
		IsSecret:    false,
		TypeOf:      manager.String,
	}

	otelHTTPExportHost *manager.ConfigKey = &manager.ConfigKey{
		Name:        "otelHTTPExportEndpoint",
		Description: "otelHTTPExportEndpoint allows one to set the address of the collector endpoint that the driver will use to send spans. If unset, it will instead try to use the default endpoint (localhost:4318). Note that the endpoint must not contain any URL path.",
		DefaultVal:  "localhost:4318",
		IsSecret:    false,
		TypeOf:      manager.String,
	}

	otelHTTPExportPath *manager.ConfigKey = &manager.ConfigKey{
		Name:        "otelHTTPExportPath",
		Description: "otelHTTPExportPath allows one to override the default URL path used for sending traces. If unset, default ('/v1/traces') will be used.",
		DefaultVal:  "/v1/traces",
		IsSecret:    false,
		TypeOf:      manager.String,
	}

	jaegerExportHost *manager.ConfigKey = &manager.ConfigKey{
		Name:        "jaegerExportHost",
		Description: "jaegerExportHost sets a host to be used in the Jaeger agent client endpoint. This option overrides any value set for the OTEL_EXPORTER_JAEGER_AGENT_HOST environment variable. If this option is not passed and the env var is not set, 'localhost' will be used by default.",
		DefaultVal:  "localhost",
		IsSecret:    false,
		TypeOf:      manager.String,
	}

	jaegerExportPort *manager.ConfigKey = &manager.ConfigKey{
		Name:        "jaegerExportPort",
		Description: "jaegerExportPort sets a port to be used in the Jaeger agent client endpoint. This option overrides any value set for the OTEL_EXPORTER_JAEGER_AGENT_PORT environment variable. If this option is not passed and the env var is not set, '6831' will be used by default.",
		DefaultVal:  "6831",
		IsSecret:    false,
		TypeOf:      manager.String,
	}
)

func (otel *OTELManager) Configs() *[]*manager.ConfigKey {
	return &[]*manager.ConfigKey{
		otelDefaultTracingStatus,
		otelExporterType,
		otelHTTPExportHost,
		otelHTTPExportPath,
		jaegerExportHost,
		jaegerExportPort,
	}
}
