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
	"fmt"

	"github.com/fire833/go-api-utils/serialization"
	"github.com/go-openapi/spec"
	"github.com/valyala/fasthttp"
)

//go:generate protoc --go_out=. --go_opt=Motel.proto=../otel otel.proto
//go:generate protoc --go_out=. --go_opt=Motel.proto=../otel otel_list.proto

var samplerStatus *spec.Schema = &spec.Schema{
	SchemaProps: spec.SchemaProps{
		Title:       "SamplerStatusList",
		Description: "Serialized object containing the enablement status of all trace operations within the current running instance.",
		Type:        []string{"object"},
		Format:      "",
		Properties: spec.SchemaProperties{
			"items": *spec.ArrayProperty(&spec.Schema{
				SchemaProps: spec.SchemaProps{
					Title:       "SamplerStatus",
					Description: "The name and status of a trace operation in the current instance.",
					Type:        []string{"object"},
					Format:      "",
					Properties: spec.SchemaProperties{
						"name": {
							SchemaProps: spec.SchemaProps{
								Title:  "name",
								Type:   []string{"string"},
								Format: "",
							},
						},
						"enabled": {
							SchemaProps: spec.SchemaProps{
								Title:  "enabled",
								Type:   []string{"boolean"},
								Format: "",
							},
						},
					},
				},
			}),
		},
	},
}

// Function to return the status of all trace operations and whether they are enabled.
func (o *OTELManager) statusHandler(ctx *fasthttp.RequestCtx) {
	list := []*SamplerStatus{}

	for name, value := range o.sampleToggle {
		list = append(list, &SamplerStatus{
			Name:    name,
			Enabled: value,
		})
	}

	serialization.MarshalBodyByAcceptHeader(ctx, &SamplerStatusList{
		Items: list,
	})
}

// method for enabling tracing for certain objects within a running app instance.
func (o *OTELManager) enable(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("NAME").(string)
	o.toggleTrace(ctx, id, true)
}

// method for diabling tracing for certain objects within a running app instance.
func (o *OTELManager) disable(ctx *fasthttp.RequestCtx) {
	id := ctx.UserValue("NAME").(string)
	o.toggleTrace(ctx, id, false)
}

func (o *OTELManager) toggleTrace(ctx *fasthttp.RequestCtx, id string, toggleOn bool) {
	e := ""
	if toggleOn {
		e = "enabled"
	} else {
		e = "disabled"
	}

	if id == "all" {
		o.sampleLock.Lock()
		for name, enabled := range o.sampleToggle {
			if enabled != toggleOn { // Only update the elements that are not enabled already.
				o.sampleToggle[name] = toggleOn
			}
		}
		o.sampleLock.Unlock()

		serialization.OKResponseHandler(ctx, 200, fmt.Sprintf("successfully %s tracing on all objects", e))
		return
	} else {
		if enabled, ok := o.sampleToggle[id]; ok && enabled != toggleOn {
			o.sampleLock.Lock()
			o.sampleToggle[id] = toggleOn
			o.sampleLock.Unlock()

			serialization.OKResponseHandler(ctx, 200, fmt.Sprintf("%s object tracing successfully %s", id, e))
		} else if !ok {
			serialization.NotFoundResponseHandler(ctx, "object not found within list of objects to toggle")
			return
		} else if enabled == toggleOn {
			serialization.OKResponseHandler(ctx, 200, fmt.Sprintf("%s object tracing already %s, ignoring", id, e))
			return
		} else {
			serialization.GenericOKResponseHandler(ctx)
			return
		}
	}
}
