package fake

import (
	"fmt"
	"math/big"
	"math/rand"

	crand "crypto/rand"
)

func init() {
	// klog.Info("initializing math/random for creating mock plain data types")

	i, _ := crand.Int(crand.Reader, big.NewInt(10000000000))
	rand.Seed(i.Int64())
}

func FakeBool() bool {
	return rand.Intn(2) != 0
}

func FakeUint64() uint64 {
	return rand.Uint64()
}

func FakeUint32() uint32 {
	return rand.Uint32()
}

func FakeInt64() int64 {
	return int64(rand.Int())
}

func FakeInt32() int32 {
	return int32(rand.Int())
}

func FakeFloat32() float32 {
	return rand.Float32()
}

func FakeFloat64() float64 {
	return rand.Float64()
}

func FakeString() string {
	return FakeStringLen(72)
}

func FakeStringFrom(choices []string) string {
	return choices[rand.Intn(len(choices))]
}

func FakeStringCeil(ceil int) string {
	return FakeStringLen(rand.Intn(ceil))
}

func FakeStringLen(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	s := make([]rune, length)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func FakeStringSlice() []string {
	return FakeStringSliceLen(32)
}

func FakeStringSliceLen(length int) []string {
	strings := []string{}

	for i := 0; i < length; i++ {
		strings = append(strings, FakeString())
	}

	return strings
}

// Returns a mocked semantic version string.
func FakeSemVer() string {
	maj := rand.Intn(10)
	min := rand.Intn(10)
	patch := rand.Intn(10)

	return fmt.Sprintf("%d.%d.%d", maj, min, patch)
}
