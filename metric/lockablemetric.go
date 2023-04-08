package metric

import "sync"

// Lockable metrics wraps a float64 metric value with a mutex for
// allowing key-level locking for a metric within a ObjectManager.
type LockableMetric struct {
	m sync.Mutex

	value float64
}

func NewLockableMetric() *LockableMetric {
	return &LockableMetric{
		value: 0,
	}
}

// Internally locks and increments the current metric value by 1.
func (m *LockableMetric) Increment() {
	m.m.Lock()
	m.value++
	m.m.Unlock()
}

// Retrieve the current metric value for exporting to outside sources.
func (m *LockableMetric) Get() float64 {
	return m.value
}

// Set the metric value to a value. This shouldn't really be used for most metrics,
// but could be used to set an initial metric value if that is required.
func (m *LockableMetric) Set(value float64) {
	m.m.Lock()
	m.value = value
	m.m.Unlock()
}
