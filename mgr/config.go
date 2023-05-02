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

// ConfigKey is a wrapper around a key that can be used
// within an application config file to configure the operation
// of an app subsystem in some way.
type ConfigKey struct {

	// Specify the actual key name for this property.
	// This can be something like "serverConcurrency",
	// "sqlDbUser", "sqlDbPass", etc.
	Name string `json:"name" yaml:"name" xml:"name"`
	// Detailed description of this configuration parameter,
	// where it is currently used within the system, a range
	// of acceptable values for this parameter (if applicable),
	// and other operational details for what that data is
	// used for and how to avoid configuring it improperly.
	Description string `json:"description" yaml:"description" xml:"description"`

	// Specify the type of this config parameter. This should
	// generally be one of the primitive types, such as string,
	// boolean, integer, float, or slice.
	TypeOf ConfigValueType `json:"typeOf" yaml:"typeOf" xml:"typeOf"`

	// If applicable, specify the default value for this
	// configuration parameter.
	DefaultVal interface{} `json:"defaultVal" yaml:"defaultVal" xml:"defaultVal"`

	// Toggle whether this configuration value should be found
	// within the secret configuration file provided by vault,
	// or in the publicly accessible config file. Generally,
	// sensitive secrets like credentials for databases should
	// be stored as a secret, while less secret configuration
	// such as ports for the server to bind to, concurrency
	// configurations, i/o timeouts, etc. should be stored in
	// the regular config file.
	IsSecret bool `json:"isSecret" yaml:"isSecret" xml:"isSecret"`
}

// ConfigValueType is an enum that describes the different types of
// Config value types that can be accessed by this ConfigKey. These
// include the basic types such as strings, bool, int, and a few slices,
// as well as a more complex time.Time type.
type ConfigValueType string

const (
	String      ConfigValueType = "String"
	StringSlice ConfigValueType = "StringSlice"
	Bool        ConfigValueType = "Bool"
	Int         ConfigValueType = "Int"
	IntSlice    ConfigValueType = "IntSlice"
	Uint        ConfigValueType = "Uint"
	Uint16      ConfigValueType = "Uint16"
	Uint32      ConfigValueType = "Uint32"
	Uint64      ConfigValueType = "Uint64"
	Float64     ConfigValueType = "Float64"
	Time        ConfigValueType = "Time"
)

// Wrapper to create a new ConfigKey to be registered with an APIManager.
func NewConfigKey(name, desc string, typeof ConfigValueType, defval interface{}, issecret bool) *ConfigKey {
	return &ConfigKey{
		Name:        name,
		Description: desc,
		TypeOf:      typeof,
		DefaultVal:  defval,
		IsSecret:    issecret,
	}
}

// WIP: Standard interface for retrieving a value from a ConfigKey.
// You will need to cast the provided interface to the appropriate
// type if you wish to use the value as its original type.
func (c *ConfigKey) RetriveValue() interface{} {
	if c.IsSecret {
		return c.retrieveValue(c.Name, false)
	} else {
		return c.retrieveValue(c.Name, true)
	}
}

func (c *ConfigKey) retrieveValue(key string, config bool) interface{} {
	if config {
		switch c.TypeOf {
		case String:
			return Manager.GetString(key)
		case StringSlice:
			return Manager.GetStringSlice(key)
		case Bool:
			return Manager.GetBool(key)
		case Int:
			return Manager.GetInt(key)
		case IntSlice:
			return Manager.GetIntSlice(key)
		case Uint:
			return Manager.GetUint(key)
		case Uint16:
			return Manager.GetUint16(key)
		case Uint32:
			return Manager.GetUint32(key)
		case Uint64:
			return Manager.GetUint64(key)
		case Float64:
			return Manager.GetFloat64(key)
		case Time:
			return Manager.GetTime(key)
		default:
			panic("typeOf not found")
		}
	} else {
		switch c.TypeOf {
		case String:
			return Manager.GetSecretString(key)
		case StringSlice:
			return Manager.GetSecretStringSlice(key)
		case Bool:
			return Manager.GetSecretBool(key)
		case Int:
			return Manager.GetSecretInt(key)
		case IntSlice:
			return Manager.GetSecretIntSlice(key)
		case Uint:
			return Manager.GetSecretUint(key)
		case Uint16:
			return Manager.GetSecretUint16(key)
		case Uint32:
			return Manager.GetSecretUint32(key)
		case Uint64:
			return Manager.GetSecretUint64(key)
		case Float64:
			return Manager.GetSecretFloat64(key)
		case Time:
			return Manager.GetSecretTime(key)
		default:
			panic("typeOf secret not found")
		}
	}
}
