/*
*	Copyright (C) 2024 Kendall Tauser
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

package serialization

import (
	"github.com/go-openapi/spec"
)

var (
	// GenericErrorResponseSchema represents the schema for a Generic Error Message to be provided
	// by app endpoints. The idea behind this schema is for it to be sent with every non-200 status code
	// the API may return, and customize the fields depending on root cause. This allows for easier error
	// handling client-side.
	GenericErrorResponseSchema *spec.Schema = NewSchema("GenericErrorResponse", "Generic response to an error.", []spec.Schema{
		NewSchemaStringProperty("error", "The http error string for this error."),
		NewSchemaStringProperty("description", "A more verbose description of this error, sometimes with a subsystem reference."),
		NewSchemaUint32Property("code", "The integer http status code corresponding to error string."),
		NewSchemaTimestampProperty("timestamp", "The timestamp of this response being returned."),
	})

	OKResponseSchema *spec.Schema = NewSchema("OKResponse", "Generic OK response.", []spec.Schema{
		NewSchemaStringProperty("message", "Human-readable message for the successful response."),
		NewSchemaTimestampProperty("timestamp", "The timestamp of this response being returned."),
		NewSchemaUint32Property("code", "The integer http status code corresponding to the response."),
	})
)

// AddResponseBoilerplate adds generic handlers for different status codes on any one API operation.
// This includes setting the MIME types as well as error status code handlers.
func AddResponseBoilerplate(resp *spec.Operation) *spec.Operation {
	authParam := spec.ParamRef("#/parameters/Authorization")
	acceptParam := spec.ParamRef("#/parameters/Accept")
	contentTypeParam := spec.ParamRef("#/parameters/Content-Type")

	resp.Parameters = append(resp.Parameters, *authParam, *acceptParam, *contentTypeParam)

	return resp.
		WithConsumes("application/json", "application/xml", "application/yaml", "application/protobuf").
		WithProduces("application/json", "application/xml", "application/yaml", "application/protobuf").
		RespondsWith(400, spec.ResponseRef("#/responses/IncorrectResponse")).
		RespondsWith(401, spec.ResponseRef("#/responses/UnauthorizedResponse")).
		RespondsWith(404, spec.ResponseRef("#/responses/NotFoundResponse")).
		RespondsWith(406, spec.ResponseRef("#/responses/UnnaceptableResponse")).
		RespondsWith(429, spec.ResponseRef("#/responses/RateLimitResponse")).
		RespondsWith(500, spec.ResponseRef("#/responses/InternalErrorResponse"))
}
