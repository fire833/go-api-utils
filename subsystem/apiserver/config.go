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

package apiserver

import manager "github.com/fire833/go-api-utils/mgr"

var (
	apiServerListenPort *manager.ConfigValue = manager.NewConfigValue(
		"apiServerListenPort",
		"Specify the listening port for this instance of APIServer. Should be an unsigned integer between 1 and 65535, but should be above 1024 preferably to avoid needing CAP_SYS_ADMIN or root privileges for the apiAPI process.",
		uint16(8080),
	)

	apiServerListenIp *manager.ConfigValue = manager.NewConfigValue(
		"apiServerListenIp",
		"Specify the listening IP for apiServer to bind to. Defaults to all available interfaces with 0.0.0.0.",
		"0.0.0.0",
	)

	apiServerConcurrency *manager.ConfigValue = manager.NewConfigValue(
		"apiServerConcurrency",
		"Specify the amount of concurrent connections to be allowed to the apiServer webserver concurrently.",
		uint(1000),
	)

	apiServerReadBufferSize *manager.ConfigValue = manager.NewConfigValue(
		"apiServerReadBufferSize",
		"Specify per-connection buffer size for requests reading. This also limits the maximum header size. Increase this buffer if your clients send multi-KB RequestURIs and/or multi-KB headers (for example, BIG cookies).",
		uint(4096),
	)

	apiServerWriteBufferSize *manager.ConfigValue = manager.NewConfigValue(
		"apiServerWriteBufferSize",
		"Per-connection buffer size for responses writing.",
		uint(4096),
	)

	apiServerReadTimeout *manager.ConfigValue = manager.NewConfigValue(
		"apiServerReadTimeout",
		"ReadTimeout is the amount of time (in seconds) allowed to read the full request including body. The connection's read deadline is reset when the connection opens, or for keep-alive connections after the first byte has been read.",
		uint(120),
	)

	apiServerWriteTimeout *manager.ConfigValue = manager.NewConfigValue(
		"apiServerWriteTimeout",
		"WriteTimeout is the maximum duration (in seconds) before timing out writes of the response. It is reset after the request handler has returned.",
		uint(120),
	)

	apiServerIdleTimeout *manager.ConfigValue = manager.NewConfigValue(
		"apiServerIdleTimeout",
		"IdleTimeout is the maximum amount of time (in seconds) to wait for the next request when keep-alive is enabled.",
		uint(120),
	)

	apiServerPrefix *manager.ConfigValue = manager.NewConfigValue(
		"apiServerPrefix",
		"Specify a prefix to serve all routes from, logically. Defaults to ''",
		"",
	)
)

func (s *APIServer) Configs() *[]*manager.ConfigValue {
	return &[]*manager.ConfigValue{
		apiServerListenPort,
		apiServerListenIp,
		apiServerConcurrency,
		apiServerReadBufferSize,
		apiServerWriteBufferSize,
		apiServerReadTimeout,
		apiServerReadTimeout,
		apiServerWriteTimeout,
		apiServerIdleTimeout,
		apiServerPrefix,
	}
}
