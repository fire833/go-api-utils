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

import "time"

func (api *APIManager) RegisterDefault(key string, value interface{}) {
	api.config.SetDefault(key, value)
}

func (api *APIManager) InConfig(key string) bool {
	return api.config.InConfig(key)
}

func (api *APIManager) IsSet(key string) bool {
	return api.config.IsSet(key)
}

func (api *APIManager) WatchConfig() {
	api.config.WatchConfig()
	api.secrets.WatchConfig()
}

// Wrapper functions around the viper config container.

func (api *APIManager) GetString(key string) string {
	return api.config.GetString(key)
}

func (api *APIManager) GetStringSlice(key string) []string {
	return api.config.GetStringSlice(key)
}

func (api *APIManager) GetBool(key string) bool {
	return api.config.GetBool(key)
}

func (api *APIManager) GetInt(key string) int {
	return api.config.GetInt(key)
}

func (api *APIManager) GetIntSlice(key string) []int {
	return api.config.GetIntSlice(key)
}

func (api *APIManager) GetUint(key string) uint {
	return api.config.GetUint(key)
}

func (api *APIManager) GetUint16(key string) uint16 {
	return api.config.GetUint16(key)
}

func (api *APIManager) GetUint32(key string) uint32 {
	return api.config.GetUint32(key)
}

func (api *APIManager) GetUint64(key string) uint64 {
	return api.config.GetUint64(key)
}

func (api *APIManager) GetTime(key string) time.Time {
	return api.config.GetTime(key)
}

func (api *APIManager) GetFloat64(key string) float64 {
	return api.config.GetFloat64(key)
}

// Wrapper functions around the viper secrets container.

func (api *APIManager) GetSecretString(key string) string {
	return api.secrets.GetString(key)
}

func (api *APIManager) GetSecretStringSlice(key string) []string {
	return api.secrets.GetStringSlice(key)
}

func (api *APIManager) GetSecretBool(key string) bool {
	return api.secrets.GetBool(key)
}

func (api *APIManager) GetSecretInt(key string) int {
	return api.secrets.GetInt(key)
}

func (api *APIManager) GetSecretIntSlice(key string) []int {
	return api.secrets.GetIntSlice(key)
}

func (api *APIManager) GetSecretUint(key string) uint {
	return api.secrets.GetUint(key)
}

func (api *APIManager) GetSecretUint16(key string) uint16 {
	return api.secrets.GetUint16(key)
}

func (api *APIManager) GetSecretUint32(key string) uint32 {
	return api.secrets.GetUint32(key)
}

func (api *APIManager) GetSecretUint64(key string) uint64 {
	return api.secrets.GetUint64(key)
}

func (api *APIManager) GetSecretTime(key string) time.Time {
	return api.secrets.GetTime(key)
}

func (api *APIManager) GetSecretFloat64(key string) float64 {
	return api.secrets.GetFloat64(key)
}