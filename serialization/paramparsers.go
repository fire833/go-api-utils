/*
*	Copyright (C) 2025 Kendall Tauser
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

import "strconv"

func Parsestring(in []byte) (string, error) {
	return string(in), nil
}

func Parsebool(in []byte) (bool, error) {
	return strconv.ParseBool(string(in))
}

func Parseint(in []byte) (int, error) {
	return strconv.Atoi(string(in))
}

func Parseint32(in []byte) (int32, error) {
	i, e := strconv.Atoi(string(in))
	return int32(i), e
}

func Parseint64(in []byte) (int64, error) {
	i, e := strconv.Atoi(string(in))
	return int64(i), e
}

func Parseuint(in []byte) (uint, error) {
	i, e := Parseint(in)
	return uint(i), e
}

func Parseuint32(in []byte) (uint32, error) {
	i, e := strconv.Atoi(string(in))
	return uint32(i), e
}

func Parseuint64(in []byte) (uint64, error) {
	i, e := strconv.Atoi(string(in))
	return uint64(i), e
}

func Parsefloat32(in []byte) (float32, error) {
	f, e := strconv.ParseFloat(string(in), 32)
	return float32(f), e
}

func Parsefloat64(in []byte) (float64, error) {
	return strconv.ParseFloat(string(in), 64)
}
