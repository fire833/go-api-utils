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
	"context"
	"reflect"

	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

type genericValue struct {
	key  string
	desc string

	defaultVal interface{}
}

type ConfigValue struct {
	genericValue
}

func NewConfigValue(key, desc string, defVal interface{}) *ConfigValue {
	return &ConfigValue{
		genericValue: genericValue{
			key:        key,
			desc:       desc,
			defaultVal: defVal,
		},
	}
}

func (c *ConfigValue) Get() interface{} {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal)
}

func (c *ConfigValue) GetString() string {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(string)
}

func (c *ConfigValue) GetStringSlice() []string {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).([]string)
}

func (c *ConfigValue) GetBool() bool {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(bool)
}

func (c *ConfigValue) GetInt() int {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(int)
}

func (c *ConfigValue) GetIntSlice() []int {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).([]int)
}

func (c *ConfigValue) GetUint() uint {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(uint)
}

func (c *ConfigValue) GetUint16() uint16 {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(uint16)
}

func (c *ConfigValue) GetUint32() uint32 {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(uint32)
}

func (c *ConfigValue) GetUint64() uint64 {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(uint64)
}

func (c *ConfigValue) GetFloat64() float64 {
	defer panicHandler(c.key)
	return mgrLGet(false, c.key, c.defaultVal).(float64)
}

type SecretValue struct {
	genericValue

	secretmountpath string
	secretpath      string

	vault bool
}

func NewSecretValue(key, desc string, defVal interface{}) *SecretValue {
	return &SecretValue{
		genericValue: genericValue{
			key:        key,
			desc:       desc,
			defaultVal: defVal,
		},
		vault: false,
	}
}

// Create a new secret vault value that will be retrieved from
// <secretmountpath>/<secretpath> within the remote vault instance.
// The provided <key> will be retrieved fro
func NewSecretVaultValue(key, desc string, defVal interface{}, secretmountpath, secretpath string) *SecretValue {
	return &SecretValue{
		genericValue: genericValue{
			key:        key, // This is usually going to be "data", but could be something else.
			desc:       desc,
			defaultVal: defVal,
		},
		secretmountpath: secretmountpath,
		secretpath:      secretpath,
		vault:           true,
	}
}

func (s *SecretValue) Get() interface{} {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal)
	}
}

func (s *SecretValue) GetString() string {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(string)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(string)
	}
}

func (s *SecretValue) GetStringSlice() []string {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).([]string)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).([]string)
	}
}

func (s *SecretValue) GetBool() bool {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(bool)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(bool)
	}
}

func (s *SecretValue) GetInt() int {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(int)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(int)
	}
}

func (s *SecretValue) GetIntSlice() []int {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).([]int)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).([]int)
	}
}

func (s *SecretValue) GetUint() uint {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(uint)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(uint)
	}
}

func (s *SecretValue) GetUint16() uint16 {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(uint16)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(uint16)
	}
}

func (s *SecretValue) GetUint32() uint32 {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(uint32)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(uint32)
	}
}

func (s *SecretValue) GetUint64() uint64 {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(uint64)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(uint64)
	}
}

func (s *SecretValue) GetFloat64() float64 {
	defer panicHandler(s.key)

	if !s.vault {
		return mgrLGet(true, s.key, s.defaultVal).(float64)
	} else {
		return mgrVGet(s.key, s.secretmountpath, s.secretpath, s.defaultVal).(float64)
	}
}

func mgrLGet(secret bool, key string, def interface{}) interface{} {
	if mgr == nil {
		klog.V(3).Infof("using default value for key %s: %v", key, def)
		return def
	}

	var m *viper.Viper
	if secret {
		if mgr.secrets == nil {
			return def
		}

		m = mgr.secrets
	} else {
		if mgr.config == nil {
			return def
		}

		m = mgr.config
	}

	switch reflect.TypeOf(def).Kind() {
	case reflect.String:
		return m.GetString(key)
	case reflect.Bool:
		return m.GetBool(key)
	case reflect.Uint:
		return m.GetUint(key)
	case reflect.Uint16:
		return m.GetUint16(key)
	case reflect.Uint32:
		return m.GetUint32(key)
	case reflect.Uint64:
		return m.GetUint64(key)
	case reflect.Float64:
		return m.GetFloat64(key)
	case reflect.Slice:
		switch reflect.TypeOf(def).Elem().Kind() {
		case reflect.String:
			return m.GetStringSlice(key)
		case reflect.Int:
			return m.GetIntSlice(key)
		default:
			return def
		}
	default:
		return def
	}
}

func mgrVGet(key, mountpath, path string, def interface{}) interface{} {
	if mgr.vault == nil {
		klog.Warningf("vault not enabled in manager, unable to access secret %s, relying on defaults", key)
		return def
	}

	if s, e := mgr.vault.KVv2(mountpath).Get(context.Background(), path); e != nil {
		klog.Warningf("unable to retrieve vault secret (%s/%s): %v, relying on defaults", mountpath, path, e)
		return def
	} else {
		if v, ok := s.Data[key]; ok {
			return v
		} else {
			klog.Warningf("key not found within secret %s, relying on defaults", key)
			return def
		}
	}
}

func panicHandler(name string) {
	if r := recover(); r != nil {
		klog.Errorf("unable to cast secret %s to desired type: %v", name, r)
	}
}
