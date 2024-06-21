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
	"os"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

// Sync start process should never return until process shutdown has been confirmed and all
// subsystems have exited as gracefully as possible.
func (m *APIManager) SyncStartProcess() {
	go m.handleSignals() // start signal handler

	if m.opts.EnableSysAPI {
		go m.startSysAPI() // start sysAPI.
	}

	for name, sys := range m.systems {
		klog.V(4).Infof("synchronously starting subsystem %s", name)
		go sys.SyncStart()
	}

	<-m.shutdown
}

func (m *APIManager) initializeSubsystems(reg *SystemRegistrar) {
	wg := new(sync.WaitGroup)
	wg.Add(int(len(m.systems)))

	for name, sys := range m.systems {
		klog.V(3).Infof("initializing subsystem %s", name)
		go func(s Subsystem, wg *sync.WaitGroup) {
			defer wg.Done()

			for i := 1; i <= 3; i++ {
				defer func() {
					if r := recover(); r != nil {
						klog.Errorf("subsystem %s panicked whilst initializing: %v", s.Name(), r)
					}
				}()

				if e := s.Initialize(reg); e != nil {
					klog.Errorf("unable to initialize subsystem %s (error: %s) %d times. Waiting 10 seconds to retry", s.Name(), e.Error(), i)
					time.Sleep(time.Second * 10) // Wait for 10 seconds to try and reinitialize
					continue
				} else {
					return
				}
			}

			klog.Errorf("3 retries attempted to initialize subsystem %s, aborting process startup", s.Name())

			m.shutdownSubsystems()

			os.Exit(1) // TODO invoke a more graceful shutdown process for subsystems that are already initialized.
		}(sys, wg)
	}

	wg.Wait()
}

func (m *APIManager) reloadSubsystems() {
	klog.V(4).Infof("reload signal received, forwarding to %d subsystems", m.count)

	wg := new(sync.WaitGroup)
	wg.Add(int(m.count))

	for name, sys := range m.systems {
		klog.V(5).Infof("sending reload update for subsystem %s", name)
		go func(sys Subsystem, wg *sync.WaitGroup) {
			defer func() {
				if r := recover(); r != nil {
					klog.Errorf("subsystem %s panicked whilst reloading: %v", sys.Name(), r)
				}
			}()

			sys.Reload()
			wg.Done()
		}(sys, wg)
	}

	wg.Wait()
	klog.V(5).Info("reloading of subsystems complete")
}

func (m *APIManager) shutdownSubsystems() {
	klog.V(4).Infof("shutdown signal received, forwarding to %d subsystems", m.count)

	wg := new(sync.WaitGroup)
	wg.Add(int(m.count))

	for name, sys := range m.systems {
		klog.V(5).Infof("sending shutdown update for subsystem %s", name)
		go func(sys Subsystem, wg *sync.WaitGroup) {
			defer func() {
				if r := recover(); r != nil {
					klog.Errorf("subsystem %s panicked whilst shutting down: %v", sys.Name(), r)
				}
			}()

			sys.Shutdown()
			wg.Done()
		}(sys, wg)
	}

	wg.Wait()
	klog.V(5).Info("shutdown of subsystems complete")

	if m.opts.EnableSysAPI {
		if e := m.server.Shutdown(); e != nil {
			klog.Errorf("unable to gracefully shutdown sysAPI: %v", e)
		}
	}
}

func (m *APIManager) preInit() {
	for _, sub := range m.systems {
		sub.PreInit()
	}
}

func (m *APIManager) postInit() {
	for _, sub := range m.systems {
		sub.PostInit()
	}
}
