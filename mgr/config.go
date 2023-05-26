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
	"errors"
	"reflect"

	"github.com/spf13/viper"
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

func (c *ConfigValue[T]) Get() (interface{}, error) {
	return mgrLGet(false, c.key, c.defaultVal)
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

func (s *SecretValue[T]) Get() (interface{}, error) {
	if !s.vault {
		return mgrLGet[T](true, s.key, s.defaultVal)
	} else {
		return mgrVGet[T](s.key, s.mountpath, s.secretpath, s.defaultVal)
	}
}

func mgrLGet[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64](secret bool, key string, def T) (interface{}, error) {
	if mgr == nil {
		return def, errors.New("manager object not initialized")
	}

	var m *viper.Viper
	if secret {
		if mgr.secrets == nil {
			return nil, errors.New("secrets not initialized")
		}

		m = mgr.secrets
	} else {
		if mgr.config == nil {
			return nil, errors.New("configs not initialized")
		}

		m = mgr.config
	}

	switch reflect.TypeOf(def).Kind() {
	case reflect.String:
		return m.GetString(key), nil
	case reflect.Bool:
		return m.GetBool(key), nil
	case reflect.Uint:
		return m.GetUint(key), nil
	case reflect.Uint16:
		return m.GetUint16(key), nil
	case reflect.Uint32:
		return m.GetUint32(key), nil
	case reflect.Uint64:
		return m.GetUint64(key), nil
	case reflect.Float64:
		return m.GetFloat64(key), nil
	case reflect.Slice:
		switch reflect.TypeOf(def).Elem().Kind() {
		case reflect.String:
			return m.GetStringSlice(key), nil
		case reflect.Int:
			return m.GetIntSlice(key), nil
		default:
			return nil, errors.New("slice type must be either string or int")
		}
	default:
		return nil, errors.New("type not defined")
	}
}

func mgrVGet[T string | []string | bool | int | []int | uint | uint16 | uint32 | uint64 | float64](key, p1, p2 string, def T) (interface{}, error) {
	if mgr.vault == nil {
		return def, errors.New("vault client not initialized")
	}

	kv := mgr.vault.KVv2(p1)

	if s, e := kv.Get(context.Background(), p2); e != nil {
		return nil, e
	} else {
		if v, ok := s.Data["key"]; ok {
			return v, nil
		} else {
			return nil, errors.New("secret not found in location in vault")
		}
	}
}
