/*
*	Copyright (C) 2024 Kendall Tauser
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
	"testing"

	"k8s.io/klog/v2"
)

func TestMain(m *testing.M) {
	m.Run()
}

func loadSimpleApp() *APIManager {
	// Reset the manager to nil on every test
	mgr = nil

	reg := &SystemRegistrar{
		AppName: "foo",
		Systems: []Subsystem{
			newMockSubsystem("thing1", 0.3, 0.2, 1),
			newMockSubsystem("thing2", 0.5, 0.1, 1.12),
			newMockSubsystem("thing3", 0.2, 0.3, 0.86),
			newMockSubsystem("thing4", 0.11, 0.2, 0.6),
			newMockSubsystem("thing5", 0.13, 0.5, 0.3),
		},
	}

	m := New(&APIManagerOpts{})
	m.Initialize(reg)
	return m
}

func TestStartup(t *testing.T) {
	t.Run("simpleStartStop", func(t *testing.T) {
		m := loadSimpleApp()

		go m.SyncStartProcess()

		klog.Info("sending shutdown for process")
		m.shutdownSubsystems()
	})

	t.Run("reload", func(t *testing.T) {
		m := loadSimpleApp()

		go m.SyncStartProcess()

		m.reloadSubsystems()
		m.reloadSubsystems()

		m.shutdownSubsystems()
	})
}
