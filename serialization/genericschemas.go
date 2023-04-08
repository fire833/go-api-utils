package serialization

import (
	"github.com/go-openapi/spec"
)

var (
	// GenericErrorResponseSchema represents the schema for a Generic Error Message to be provided
	// by app endpoints. The idea behind this schema is for it to be sent with every non-200 status code
	// the API may return, and customize the fields depending on root cause. This allows for easier error
	// handling client-side.
	GenericErrorResponseSchema *spec.Schema = &spec.Schema{
		SchemaProps: spec.SchemaProps{
			// ID:       "GenericErrorResponse",
			Title:    "GenericErrorResponse",
			Type:     []string{"object"},
			Required: []string{"error", "description", "code", "timestamp"},
			Properties: spec.SchemaProperties{
				"error": {
					SchemaProps: spec.SchemaProps{
						// ID:          "error",
						Title:       "error",
						Type:        []string{"string"},
						Format:      "",
						Description: "The http error string for this error.",
					},
				},
				"description": {
					SchemaProps: spec.SchemaProps{
						// ID:          "description",
						Title:       "description",
						Type:        []string{"string"},
						Format:      "",
						Description: "A more verbose description of this error, sometimes with a subsystem reference.",
					},
				},
				"code": {
					SchemaProps: spec.SchemaProps{
						// ID:          "code",
						Title:       "code",
						Type:        []string{"integer"},
						Format:      "uint32",
						Description: "The integer http status code corresponding to error string.",
					},
				},
				"timestamp": {
					SchemaProps: spec.SchemaProps{
						// ID:          "timestamp",
						Title:       "timestamp",
						Type:        []string{"string"},
						Format:      "date-time",
						Description: "The timestamp of this response being returned.",
					},
				},
			},
		},
	}

	OKResponseSchema *spec.Schema = &spec.Schema{
		SchemaProps: spec.SchemaProps{
			// ID:       "OKResponse",
			Title:    "OKResponse",
			Type:     []string{"object"},
			Required: []string{"message", "timestamp", "code"},
			Properties: spec.SchemaProperties{
				"message": {
					SchemaProps: spec.SchemaProps{
						// ID:          "message",
						Title:       "message",
						Type:        []string{"string"},
						Format:      "",
						Description: "Human-readable message for the successful response.",
					},
				},
				"timestamp": {
					SchemaProps: spec.SchemaProps{
						// ID:          "timestamp",
						Title:       "timestamp",
						Type:        []string{"string"},
						Format:      "date-time",
						Description: "The timestamp of this response being returned.",
					},
				},
				"code": {
					SchemaProps: spec.SchemaProps{
						// ID:          "code",
						Title:       "code",
						Type:        []string{"integer"},
						Format:      "uint32",
						Description: "The integer http status code corresponding to the response.",
					},
				},
			},
		},
	}
)

// AddResponseBoilerplate adds generic handlers for different status codes on any one API operation.
// This includes setting the MIME types as well as error status code handlers.
func AddResponseBoilerplate(resp *spec.Operation) *spec.Operation {
	authParam := spec.ParamRef("#/parameters/Authorization")
	acceptParam := spec.ParamRef("#/parameters/Accept")
	contentTypeParam := spec.ParamRef("#/parameters/Content-Type")

	resp.Parameters = append(resp.Parameters, *authParam, *acceptParam, *contentTypeParam)

	return resp.
		WithConsumes("application/json", "application/xml", "application/yaml", "application/VTAPI+Protobuf").
		WithProduces("application/json", "application/xml", "application/yaml", "application/VTAPI+Protobuf").
		RespondsWith(400, spec.ResponseRef("#/responses/IncorrectResponse")).
		RespondsWith(401, spec.ResponseRef("#/responses/UnauthorizedResponse")).
		RespondsWith(404, spec.ResponseRef("#/responses/NotFoundResponse")).
		RespondsWith(406, spec.ResponseRef("#/responses/UnnaceptableResponse")).
		RespondsWith(429, spec.ResponseRef("#/responses/RateLimitResponse")).
		RespondsWith(500, spec.ResponseRef("#/responses/InternalErrorResponse"))
}
