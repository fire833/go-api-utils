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

package otel

import manager "github.com/fire833/go-api-utils/mgr"

var (
	otelDefaultTracingStatus *manager.ConfigValue = manager.NewConfigValue(
		"otelDefaultTracingStatus",
		"Toggle whether or not tracing is default enabled for all types within the running instance or not. This will be the initial boolean value applied to the internal hashmap.",
		false,
	)

	otelExporterType *manager.ConfigValue = manager.NewConfigValue(
		"otelExporterType",
		"Specify the span exporter you wish to use to export traces from opentelemetry tracing within this app instance.",
		"otlphttp",
	)

	otelHTTPExportHost *manager.ConfigValue = manager.NewConfigValue(
		"otelHTTPExportEndpoint",
		"otelHTTPExportEndpoint allows one to set the address of the collector endpoint that the driver will use to send spans. If unset, it will instead try to use the default endpoint (localhost:4318). Note that the endpoint must not contain any URL path.",
		"localhost:4318",
	)

	otelHTTPExportPath *manager.ConfigValue = manager.NewConfigValue(
		"otelHTTPExportPath",
		"otelHTTPExportPath allows one to override the default URL path used for sending traces. If unset, default ('/v1/traces') will be used.",
		"/v1/traces",
	)

	jaegerExportHost *manager.ConfigValue = manager.NewConfigValue(
		"jaegerExportHost",
		"jaegerExportHost sets a host to be used in the Jaeger agent client endpoint. This option overrides any value set for the OTEL_EXPORTER_JAEGER_AGENT_HOST environment variable. If this option is not passed and the env var is not set, 'localhost' will be used by default.",
		"localhost",
	)

	jaegerExportPort *manager.ConfigValue = manager.NewConfigValue(
		"jaegerExportPort",
		"jaegerExportPort sets a port to be used in the Jaeger agent client endpoint. This option overrides any value set for the OTEL_EXPORTER_JAEGER_AGENT_PORT environment variable. If this option is not passed and the env var is not set, '6831' will be used by default.",
		"6831",
	)
)

func (otel *OTELManager) Configs() *[]*manager.ConfigValue {
	return &[]*manager.ConfigValue{
		otelDefaultTracingStatus,
		otelExporterType,
		otelHTTPExportHost,
		otelHTTPExportPath,
		jaegerExportHost,
		jaegerExportPort,
	}
}
