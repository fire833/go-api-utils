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
	"math"
	"testing"
)

func TestConfigValues(t *testing.T) {
	var u16 uint16
	for u16 = 0; u16 < math.MaxUint16; u16++ {
		var testUint16 *ConfigValue = NewConfigValue("foo", "bar", uint16(u16))
		if testUint16.GetUint16() != u16 {
			t.Errorf("test config value does not return correct value: %d", u16)
		}
	}

	var u32 uint32
	for u32 = 0; u32 < math.MaxUint32/2; u32 += 10000 {
		var testUint16 *ConfigValue = NewConfigValue("foo", "bar", uint32(u32))
		if testUint16.GetUint32() != u32 {
			t.Errorf("test config value does not return correct value: %d", u32)
		}
	}

	var u64 uint64
	for u64 = 0; u64 < math.MaxUint64/2; u64 += 9999999999995 {
		var testUint16 *ConfigValue = NewConfigValue("foo", "bar", uint64(u64))
		if testUint16.GetUint64() != u64 {
			t.Errorf("test config value does not return correct value: %d", u64)
		}
	}
}
