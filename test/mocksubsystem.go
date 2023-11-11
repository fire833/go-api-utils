package test

import (
	"time"

	"github.com/fire833/go-api-utils/mgr"
	"k8s.io/klog/v2"
)

type MockSubsystem struct {
	mgr.Subsystem

	name string

	initDelay, reloadDelay, shutdownDelay int
}

// Return the name of this subsystem for referencing.
func (m *MockSubsystem) Name() string

// Optional hooks to be run before subsystem initialization.
func (m *MockSubsystem) PreInit() {
}

func (m *MockSubsystem) PostInit() {
}

// Configs and Secrets return the ConfigValues and SecretValues for the subsystem for management
// by the APIMAnager. This includes adding defaults to viper and populating configuration
// defaulting commands.
func (m *MockSubsystem) Configs() *[]*mgr.ConfigValue {
	return &[]*mgr.ConfigValue{}
}

func (m *MockSubsystem) Secrets() *[]*mgr.SecretValue {
	return &[]*mgr.SecretValue{}
}

// Starts up this subsystem, if it returns an error, will try to reinitalize
// the subsystem with backoff until an error is no longer returned.
//
// The waitgroup should be immediately deferred in the called function.
func (m *MockSubsystem) Initialize(reg *mgr.SystemRegistrar) error {
	return nil
}

// If a subsystem needs to be synchronously called by the manager (IE you will call a method
// that should never return), wrap that method here and it will be held open by the manager
// as long as the process is alive.
func (m *MockSubsystem) SyncStart() {
	return
}

// If the configuration is found to be changed, the manager will call this callback
// to inform the subsystem to refresh itself given the new configuration changes.
//
// The waitgroup should be immediately deferred in the called function.
func (m *MockSubsystem) Reload() {
	time.Sleep(time.Second * time.Duration(m.reloadDelay))
	klog.Infof("reloaded subsystem")
	return
}

// If a kill signal is sent to the process, the manager will inform all subsystems of shutdown
// through this callback.
//
// The waitgroup should be immediately deferred in the called function.
func (m *MockSubsystem) Shutdown() {
}

// This callback can be invoked at any point in execution by the manager to determine the
// current status of the subsystem.
func (m *MockSubsystem) Status() *mgr.SubsystemStatus {
	return nil
}

func NewMockSubsystem(name string, initDelay, reloadDelay, shutdownDelay int) *MockSubsystem {
	return &MockSubsystem{
		name:          name,
		initDelay:     initDelay,
		reloadDelay:   reloadDelay,
		shutdownDelay: shutdownDelay,
	}
}
