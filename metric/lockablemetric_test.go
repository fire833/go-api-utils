package metric

import (
	"testing"
)

func BenchmarkLockableMetric(b *testing.B) {

	lock := NewLockableMetric()

	b.ResetTimer()

	// b.SetParallelism(3)
	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			lock.Increment()
		}
	})
}
