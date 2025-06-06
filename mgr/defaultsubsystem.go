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

package mgr

import (
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

func (d *DefaultSubsystem) SetGlobal() {}

func (d *DefaultSubsystem) Initialize(reg *SystemRegistrar) error {
	d.IsInitialized = true
	return nil
}

// NOP PreInit
func (d *DefaultSubsystem) PreInit() {}

// NOP SyncStart
func (d *DefaultSubsystem) SyncStart() {}

// NOP PostInit
func (d *DefaultSubsystem) PostInit() {}

func (d *DefaultSubsystem) Configs() *[]*ConfigValue {
	return &[]*ConfigValue{}
}

func (d *DefaultSubsystem) Secrets() *[]*SecretValue {
	return &[]*SecretValue{}
}

// NOP to reload the subsystem
func (d *DefaultSubsystem) Reload() {}

// NOP to shutdown the subsystem
func (d *DefaultSubsystem) Shutdown() {
	d.IsShutdown = true
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
