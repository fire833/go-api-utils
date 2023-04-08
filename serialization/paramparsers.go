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
