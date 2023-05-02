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

package mgr

import (
	"github.com/fasthttp/router"
	"github.com/go-openapi/spec"
	"github.com/prometheus/client_golang/prometheus"
)

// SystemRegistrar allows for multiple different applications within an API
// project to utilize one manager interface for multiple subsystem-based applications.
// Developers can specify the application configuration to specify the desired subsystem
// behavior for a specific implementation of the APIManager, Ie you can specify only
// certain subsystems to be included with application A, but others can be enabled/disabled
// exclusively for application B, whch also initializes itself with the APIManager interface
// code in this package.
type SystemRegistrar struct {
	// The name of this application, will be used for logging and
	// other times when a specific name for the app is needed.
	AppName string

	// The subsystems you want to begin with this application,
	// IN THE DESIRED BOOT ORDER. Ie, if the application relies
	// on a database system to connect, you should run that before
	// you begin serving your HTTP server.
	Systems []Subsystem

	// Registration is the data access that you want to start with
	// this application. For example, in the case of app, you will
	// probably wish to register the appdata.app object here, and
	// this will be passed to subsystems for use.
	Registration AppRegistration
}

func NewRegistrar(name string, dataRegistration AppRegistration, systems ...Subsystem) *SystemRegistrar {
	return &SystemRegistrar{
		AppName:      name,
		Registration: dataRegistration,
		Systems:      systems,
	}
}

// AppRegistration is an interface for applications to register their endpoints, collectors,
// and other operations with the greater manager context and subsystems registered with this
// application.
type AppRegistration interface {
	// RegisterEndpoints should register all endpoints for an application for serving
	// with the server subsystem.
	RegisterEndpoints(prefix string, router *router.Router)

	// RegisterCollectors should register all collectors with the global prometheus registry for
	// serving metrics via SysAPI.
	RegisterCollectors(registry *prometheus.Registry)

	// RegisterSwagger2 should register all definitions and endpoints with a openAPI 2 object
	// for exporting API documentation via multiple interfaces.
	RegisterSwagger2(prefix string, paths *spec.Paths, definitions spec.Definitions)

	// RegisterOTELTraces should register all sub-operation labels that can be found in spans from
	// this application. This allows for admins to toggle which spans are collected at runtime.
	RegisterOTELTraces() []string
}

// GetConfigKeys will retrieve all the configkeys for this registrar.
func (r *SystemRegistrar) GetConfigKeys() []*ConfigKey {
	keys := []*ConfigKey{}

	for _, sys := range r.Systems {
		for _, key := range *sys.Configs() {
			if !key.IsSecret {
				keys = append(keys, key)
			}
		}
	}

	return keys
}

func (r *SystemRegistrar) GetSecretConfigKeys() []*ConfigKey {
	keys := []*ConfigKey{}

	for _, sys := range r.Systems {
		for _, key := range *sys.Configs() {
			if key.IsSecret {
				keys = append(keys, key)
			}
		}
	}

	return keys
}
