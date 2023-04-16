package manager

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// DefaultSubsystem is a NOP subsystem that simply fulfills the Subsystem
// interface and does nothing else. If you want forwards compatibility with subsystems
// that you develop, please utilize the DefaultSubsystem as an embedded
// member of your objects so that the Subsystem interface is always fulfilled with your types.
type DefaultSubsystem struct {
	prometheus.Collector

	IsInitialized bool
	IsShutdown    bool
}

// Default naming
func (d *DefaultSubsystem) Name() string { return "default" }

func (d *DefaultSubsystem) Configs() *[]*ConfigKey {
	return &[]*ConfigKey{}
}

func (d *DefaultSubsystem) SetGlobal() {}

func (d *DefaultSubsystem) Initialize(wg *sync.WaitGroup, reg *SystemRegistrar) error {
	defer wg.Done()
	d.IsInitialized = true
	return nil
}

// NOP SyncStart
func (d *DefaultSubsystem) SyncStart() {}

// NOP to reload the subsystem
func (d *DefaultSubsystem) Reload(wg *sync.WaitGroup) { wg.Done() }

// NOP to shutdown the subsystem
func (d *DefaultSubsystem) Shutdown(wg *sync.WaitGroup) {
	d.IsShutdown = true
	wg.Done()
}

// Return nothing since this subsystem does nothing, but you should be able to fill this out at runtime
// so the APIManager can effectively make decision on process lifecycle.
func (d *DefaultSubsystem) Status() *SubsystemStatus {
	return &SubsystemStatus{
		Name:          d.Name(),
		IsInitialized: d.IsInitialized,
		IsShutdown:    d.IsShutdown,
		Meta:          nil,
	}
}

// NOP to implement prometheus Collector interface.
func (d *DefaultSubsystem) Describe(chan<- *prometheus.Desc) {}

// NOP to implement prometheus Collector interface.
func (d *DefaultSubsystem) Collect(chan<- prometheus.Metric) {}
