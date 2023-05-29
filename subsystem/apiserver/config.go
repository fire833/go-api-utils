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
	apiServerListenPort *manager.ConfigValue[uint16] = manager.NewConfigValue[uint16](
		"apiServerListenPort",
		"Specify the listening port for this instance of APIServer. Should be an unsigned integer between 1 and 65535, but should be above 1024 preferably to avoid needing CAP_SYS_ADMIN or root privileges for the apiAPI process.",
		uint16(8080),
	)

	apiServerListenIp *manager.ConfigKey = &manager.ConfigKey{
		Name:        "apiServerListenIp",
		Description: "Specify the listening IP for apiServer to bind to. Defaults to all available interfaces with 0.0.0.0.",
		TypeOf:      manager.String,
		DefaultVal:  "0.0.0.0",
		IsSecret:    false,
	}

	apiServerConcurrency *manager.ConfigValue[int] = manager.NewConfigValue[int]()
		Name:        "apiServerConcurrency",
		Description: "Specify the amount of concurrent connections to be allowed to the apiServer webserver concurrently.",
		IsSecret:    false,
		DefaultVal:  1000,
		TypeOf:      manager.Int,
	}

	apiServerReadBufferSize *manager.ConfigKey = &manager.ConfigKey{
		Name:        "apiServerReadBufferSize",
		Description: "Specify per-connection buffer size for requests reading. This also limits the maximum header size. Increase this buffer if your clients send multi-KB RequestURIs and/or multi-KB headers (for example, BIG cookies).",
		TypeOf:      manager.Int,
		DefaultVal:  4096,
		IsSecret:    false,
	}

	apiServerWriteBufferSize *manager.ConfigKey = &manager.ConfigKey{
		Name:        "apiServerWriteBufferSize",
		Description: "Per-connection buffer size for responses writing.",
		TypeOf:      manager.Int,
		DefaultVal:  4096,
		IsSecret:    false,
	}

	apiServerReadTimeout *manager.ConfigKey = &manager.ConfigKey{
		Name:        "apiServerReadTimeout",
		Description: "ReadTimeout is the amount of time (in seconds) allowed to read the full request including body. The connection's read deadline is reset when the connection opens, or for keep-alive connections after the first byte has been read.",
		TypeOf:      manager.Int,
		DefaultVal:  120,
		IsSecret:    false,
	}

	apiServerWriteTimeout *manager.ConfigKey = &manager.ConfigKey{
		Name:        "apiServerWriteTimeout",
		Description: "WriteTimeout is the maximum duration (in seconds) before timing out writes of the response. It is reset after the request handler has returned.",
		TypeOf:      manager.Int,
		DefaultVal:  120,
		IsSecret:    false,
	}

	apiServerIdleTimeout *manager.ConfigKey = &manager.ConfigKey{
		Name:        "apiServerIdleTimeout",
		Description: "IdleTimeout is the maximum amount of time (in seconds) to wait for the next request when keep-alive is enabled.",
		TypeOf:      manager.Int,
		DefaultVal:  120,
		IsSecret:    false,
	}

	apiServerPrefix *manager.ConfigKey = &manager.ConfigKey{
		Name:        "apiServerPrefix",
		Description: "Specify a prefix to serve all routes from, logically. Defaults to ''",
		TypeOf:      manager.String,
		DefaultVal:  "",
		IsSecret:    false,
	}
)

func (s *APIServer) Configs() *[]*manager.ConfigKey {
	return &[]*manager.ConfigKey{
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
