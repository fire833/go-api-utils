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
	"context"
	"reflect"

	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

type Value[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64] interface {
}

type genericValue[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64] struct {
	Value[T]

	key  string
	desc string

	defaultVal T
}

type ConfigValue[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64] struct {
	genericValue[T]
}

func NewConfigValue[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64](key, desc string, defVal T) *ConfigValue[T] {
	return &ConfigValue[T]{
		genericValue: genericValue[T]{
			key:        key,
			desc:       desc,
			defaultVal: defVal,
		},
	}
}

func (c *ConfigValue[T]) GetString() string {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(string)
}

func (c *ConfigValue[T]) GetStringSlice() []string {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).([]string)
}

func (c *ConfigValue[T]) GetBool() bool {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(bool)
}

func (c *ConfigValue[T]) GetInt() int {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(int)
}

func (c *ConfigValue[T]) GetIntSlice() []int {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).([]int)
}

func (c *ConfigValue[T]) GetUint() uint {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(uint)
}

func (c *ConfigValue[T]) GetUint16() uint16 {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(uint16)
}

func (c *ConfigValue[T]) GetUint32() uint32 {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(uint32)
}

func (c *ConfigValue[T]) GetUint64() uint64 {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(uint64)
}

func (c *ConfigValue[T]) GetFloat64() float64 {
	defer panicHandler()
	return mgrLGet(false, c.key, c.defaultVal).(float64)
}

type SecretValue[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64] struct {
	genericValue[T]

	mountpath  string
	secretpath string

	vault bool
}

func NewSecretValue[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64](key, desc string, defVal T) *SecretValue[T] {
	return &SecretValue[T]{
		genericValue: genericValue[T]{
			key:        key,
			desc:       desc,
			defaultVal: defVal,
		},
		vault: false,
	}
}

func NewSecretVaultValue[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64](key, desc string, defVal T, mountpath, secretpath string) *SecretValue[T] {
	return &SecretValue[T]{
		genericValue: genericValue[T]{
			key:        key,
			desc:       desc,
			defaultVal: defVal,
		},
		mountpath:  mountpath,
		secretpath: secretpath,
		vault:      true,
	}
}

func (s *SecretValue[T]) GetString() string {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(string)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(string)
	}
}

func (s *SecretValue[T]) GetStringSlice() []string {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).([]string)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).([]string)
	}
}

func (s *SecretValue[T]) GetBool() bool {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(bool)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(bool)
	}
}

func (s *SecretValue[T]) GetInt() int {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(int)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(int)
	}
}

func (s *SecretValue[T]) GetIntSlice() []int {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).([]int)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).([]int)
	}
}

func (s *SecretValue[T]) GetUint() uint {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(uint)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(uint)
	}
}

func (s *SecretValue[T]) GetUint16() uint16 {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(uint16)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(uint16)
	}
}

func (s *SecretValue[T]) GetUint32() uint32 {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(uint32)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(uint32)
	}
}

func (s *SecretValue[T]) Get() uint64 {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(uint64)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(uint64)
	}
}

func (s *SecretValue[T]) GetFloat64() float64 {
	defer panicHandler()

	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal).(float64)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal).(float64)
	}
}

func mgrLGet[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64](secret bool, key string, def T) interface{} {
	if mgr == nil {
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

func mgrVGet[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64](key, p1, p2 string, def T) interface{} {
	if mgr.vault == nil {
		return def
	}

	kv := mgr.vault.KVv2(p1)

	if s, e := kv.Get(context.Background(), p2); e != nil {
		return nil
	} else {
		if v, ok := s.Data["key"]; ok {
			return v
		} else {
			return def
		}
	}
}

func panicHandler() {
	if r := recover(); r != nil {
		klog.Errorf("unable to cast secret to desired type: %v", r)
	}
}
