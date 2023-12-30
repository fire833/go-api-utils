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

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

type mockSubsystem struct {
	Subsystem

	name string

	isInit, isDown bool

	initDelay, reloadDelay, shutdownDelay float64
}

// Return the name of this subsystem for referencing.
func (m *mockSubsystem) Name() string {
	return m.name
}

// Optional hooks to be run before subsystem initialization.
func (m *mockSubsystem) PreInit() {
}

func (m *mockSubsystem) PostInit() {
}

// Configs and Secrets return the ConfigValues and SecretValues for the subsystem for management
// by the APIMAnager. This includes adding defaults to viper and populating configuration
// defaulting commands.
func (m *mockSubsystem) Configs() *[]*ConfigValue {
	return &[]*ConfigValue{}
}

func (m *mockSubsystem) Secrets() *[]*SecretValue {
	return &[]*SecretValue{}
}

// Starts up this subsystem, if it returns an error, will try to reinitalize
// the subsystem with backoff until an error is no longer returned.
//
// The waitgroup should be immediately deferred in the called function.
func (m *mockSubsystem) Initialize(reg *SystemRegistrar) error {
	time.Sleep(time.Second * time.Duration(m.initDelay))
	klog.Infof("initialized subsystem %s", m.name)
	m.isInit = true
	return nil
}

// If a subsystem needs to be synchronously called by the manager (IE you will call a method
// that should never return), wrap that method here and it will be held open by the manager
// as long as the process is alive.
func (m *mockSubsystem) SyncStart() {
	return
}

// If the configuration is found to be changed, the manager will call this callback
// to inform the subsystem to refresh itself given the new configuration changes.
//
// The waitgroup should be immediately deferred in the called function.
func (m *mockSubsystem) Reload() {
	time.Sleep(time.Second * time.Duration(m.reloadDelay))
	klog.Infof("reloaded subsystem %s", m.name)
	return
}

// If a kill signal is sent to the process, the manager will inform all subsystems of shutdown
// through this callback.
//
// The waitgroup should be immediately deferred in the called function.
func (m *mockSubsystem) Shutdown() {
	time.Sleep(time.Second * time.Duration(m.shutdownDelay))
	klog.Infof("shutdown subsystem %s", m.name)
	m.isDown = true
}

// This callback can be invoked at any point in execution by the manager to determine the
// current status of the subsystem.
func (m *mockSubsystem) Status() *SubsystemStatus {
	return &SubsystemStatus{
		Name:          m.name,
		IsInitialized: m.isInit,
		IsShutdown:    m.isDown,
		Meta:          nil,
	}
}

func newMockSubsystem(name string, initDelay, reloadDelay, shutdownDelay float64) *mockSubsystem {
	return &mockSubsystem{
		name:          name,
		initDelay:     initDelay,
		reloadDelay:   reloadDelay,
		shutdownDelay: shutdownDelay,
		isInit:        false,
		isDown:        false,
	}
}

func mockManager() *APIManager {
	return &APIManager{
		count:     0,
		systems:   make(map[string]Subsystem),
		shutdown:  make(chan uint8),
		config:    viper.New(),
		secrets:   viper.New(),
		sigHandle: make(chan os.Signal),
	}
}

func TestAPIManager_initializeSubsystems(t *testing.T) {
	type fields struct {
		count     uint
		systems   map[string]Subsystem
		config    *viper.Viper
		secrets   *viper.Viper
		shutdown  chan uint8
		sigHandle chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
		reg    *SystemRegistrar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &APIManager{
				count:     tt.fields.count,
				systems:   tt.fields.systems,
				config:    tt.fields.config,
				secrets:   tt.fields.secrets,
				shutdown:  tt.fields.shutdown,
				sigHandle: tt.fields.sigHandle,
			}
			m.initializeSubsystems(tt.reg)
		})
	}
}

func TestAPIManager_reloadSubsystems(t *testing.T) {
	type fields struct {
		count     uint
		systems   map[string]Subsystem
		config    *viper.Viper
		secrets   *viper.Viper
		shutdown  chan uint8
		sigHandle chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &APIManager{
				count:     tt.fields.count,
				systems:   tt.fields.systems,
				config:    tt.fields.config,
				secrets:   tt.fields.secrets,
				shutdown:  tt.fields.shutdown,
				sigHandle: tt.fields.sigHandle,
			}
			m.reloadSubsystems()
		})
	}
}

func TestAPIManager_shutdownSubsystems(t *testing.T) {
	type fields struct {
		count     uint
		systems   map[string]Subsystem
		config    *viper.Viper
		secrets   *viper.Viper
		shutdown  chan uint8
		sigHandle chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &APIManager{
				count:     tt.fields.count,
				systems:   tt.fields.systems,
				config:    tt.fields.config,
				secrets:   tt.fields.secrets,
				shutdown:  tt.fields.shutdown,
				sigHandle: tt.fields.sigHandle,
			}
			m.shutdownSubsystems()
		})
	}
}
