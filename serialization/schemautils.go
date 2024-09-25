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

import "github.com/go-openapi/spec"

func NewSchema(name, desc string, properties []spec.Schema) *spec.Schema {
	required := []string{}

	for _, prop := range properties {
		required = append(required, prop.Title)
	}

	return NewSchemaRequired(name, desc, properties, []string{})
}

func NewSchemaRequired(name, desc string, properties []spec.Schema, required []string) *spec.Schema {
	props := spec.SchemaProperties{}

	for _, prop := range properties {
		props[prop.Title] = prop
	}

	return &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Title:       name,
			Description: desc,
			Type:        []string{"object"},
			Required:    required,
			Format:      "",
			Properties:  props,
		},
	}
}

// Generic constructor of a new operation within the API.
func NewOperation(id, desc string, tags []string, params []*spec.Parameter, responses map[uint]*spec.Response) *spec.Operation {
	var (
		ConsumeMIMES []string = []string{"application/json", "application/xml", "application/yaml", "application/API+Protobuf"}
		ProduceMIMES []string = []string{"application/json", "application/xml", "application/yaml", "application/API+Protobuf"}
	)

	op := spec.NewOperation(id).WithTags(tags...).
		WithSummary(desc).WithConsumes(ConsumeMIMES...).WithProduces(ProduceMIMES...).
		WithDescription(desc)

	for _, param := range params {
		op = op.AddParam(param)
	}

	for code, resp := range responses {
		op = op.RespondsWith(int(code), resp)
	}

	return op
}

func NewGetPath(op *spec.Operation) *spec.PathItem {
	return &spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Get: op,
		},
	}
}

func NewPostPath(op *spec.Operation) *spec.PathItem {
	return &spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Post: op,
		},
	}
}

func NewPutPath(op *spec.Operation) *spec.PathItem {
	return &spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Put: op,
		},
	}
}

func NewDeletePath(op *spec.Operation) *spec.PathItem {
	return &spec.PathItem{
		PathItemProps: spec.PathItemProps{
			Delete: op,
		},
	}
}

func NewSchemaStringProperty(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "string", "")
}

func NewSchemaUintProperty(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "uint")
}

func NewSchemaUint64Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "uint64")
}

func NewSchemaUint32Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "uint32")
}

func NewSchemaUint16Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "uint16")
}

func NewSchemaUint8Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "uint8")
}

func NewSchemaIntProperty(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "int")
}

func NewSchemaInt64Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "int64")
}

func NewSchemaInt32Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "int32")
}

func NewSchemaInt16Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "int16")
}

func NewSchemaInt8Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "int8")
}

func NewSchemaFloat32Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "float32")
}

func NewSchemaFloat64Property(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "number", "float64")
}

func NewSchemaObjectProperty(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "object", "")
}

func NewSchemaBooleanProperty(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "boolean", "")
}

func NewSchemaTimestampProperty(name, desc string) spec.Schema {
	return NewSchemaProperty(name, desc, "string", "time-format")
}

func NewSchemaEnumProperty(name, desc, propType, format string, enum []interface{}) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Title:       name,
			Description: desc,
			Type:        []string{propType},
			Format:      format,
			Enum:        enum,
		},
	}
}

func NewSchemaProperty(name, desc, propType, format string) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Title:       name,
			Description: desc,
			Type:        []string{propType},
			Format:      format,
		},
	}
}
