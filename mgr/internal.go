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

package mgr

import (
	"github.com/fire833/go-api-utils/serialization"
	"github.com/go-openapi/spec"
)

var (
	subsystemStatusSchema *spec.Schema = &spec.Schema{
		SchemaProps: spec.SchemaProps{
			Title:       "SubsystemStatus",
			Description: "Serialized object describing the status and additional metadata about a particular subsystem.",
			Type:        []string{"object"},
			Format:      "",
			Properties: spec.SchemaProperties{
				"name": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "The name of this subsystem from the source code.",
						Title:       "name",
						Type:        []string{"string"},
						Format:      "",
					},
				},
				"isInitialized": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "Boolean value of whether this subsystem is initialized.",
						Title:       "isInitialized",
						Type:        []string{"boolean"},
						Format:      "",
					},
				},
				"isShutdown": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "Boolean value of whether this subsystem is shutdown.",
						Title:       "isShutdown",
						Type:        []string{"boolean"},
						Format:      "",
					},
				},
				"meta": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "Arbitrary metadata object emitted by this subsystem.",
						Title:       "meta",
						Type:        []string{"object"},
						Format:      "",
					},
				},
			},
		},
	}

	buildInfoSchema *spec.Schema = &spec.Schema{
		SwaggerSchemaProps: spec.SwaggerSchemaProps{
			Example: BuildInfo{
				Version:   "1.0.0",
				Commit:    "3c03823782098c24e57cf779643a5a2d6883e1b6",
				BuildTime: "Sun Jan 1 00:00:01 CDT 2022",
				Os:        "linux",
				Arch:      "amd64",
			},
		},
		SchemaProps: spec.SchemaProps{
			Title:       "BuildInfo",
			Description: "Serialized object describing the build information used for creating this application instance.",
			Type:        []string{"object"},
			Format:      "",
			Properties: spec.SchemaProperties{
				"version": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "The specific version of application.",
						Title:       "version",
						Type:        []string{"string"},
						Format:      "",
					},
				},
				"commit": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "The git commit from which this application instance is derived.",
						Title:       "commit",
						Type:        []string{"string"},
						Format:      "",
					},
				},
				"buildTime": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "The time at which this instance of application was compiled.",
						Title:       "buildTime",
						Type:        []string{"string"},
						Format:      "",
					},
				},
				"os": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "The OS this binary is meant for.",
						Title:       "os",
						Type:        []string{"string"},
						Format:      "",
					},
				},
				"arch": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Description: "The platform this binary is meant for.",
						Title:       "arch",
						Type:        []string{"string"},
						Format:      "",
					},
				},
			},
		},
	}

	configKeySchema *spec.Schema = serialization.NewSchema("ConfigKey", "Serialized object describing the value of a config/secret key/value within the current process.", []spec.Schema{
		serialization.NewSchemaObjectProperty("value", "The current value of that config key in memory."),
		// Meta sub-object
		*serialization.NewSchema("meta", "Metadata associated with this config key.", []spec.Schema{
			serialization.NewSchemaStringProperty("name", "Specify the actual key name for this property. This can be something like 'serverConcurrency', 'sqlDbUser', 'sqlDbPass', etc."),
			serialization.NewSchemaStringProperty("description", "Description of this key/value pair, what its used for, and any edge case information about it."),
			serialization.NewSchemaEnumProperty("typeOf", "The type of this value", "string", "",
				[]interface{}{"String", "StringSlice", "Bool", "Int", "IntSlice", "Uint", "Uint16", "Uint32", "Uint64", "Float64", "Time"}),
			serialization.NewSchemaObjectProperty("defaultVal", "Default value for this config key."),
			serialization.NewSchemaBooleanProperty("isSecret", "Whether or not this configkey value is to be regarded as a secret."),
		}),
	})
)
