package serialization

import (
	"fmt"
	"testing"
)

func BenchmarkSprintFInt(b *testing.B) {
	testInt := 465643

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", testInt)
	}
}

func BenchmarkSprintFBool(b *testing.B) {
	testBool := true

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", testBool)
	}
}

func BenchmarkSprintFString(b *testing.B) {
	testStr := "this is a string of length 29"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", testStr)
	}
}

func BenchmarkSprintFFloat(b *testing.B) {
	testFloat := 1.234678236482467439873

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("%v", testFloat)
	}
}
