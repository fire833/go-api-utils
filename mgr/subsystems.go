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
	"os"
	"sync"
	"time"

	"k8s.io/klog/v2"
)

func (m *APIManager) initializeSubsystems(reg *SystemRegistrar) {
	wg := new(sync.WaitGroup)
	wg.Add(int(len(m.systems)))

	// Make a channel to asynchrnously collect if any subsystems fail,
	// in which case we want to then close out all subsystems and exit the process.
	errChan := make(chan bool, len(m.systems))

	for name, sys := range m.systems {
		klog.V(3).Infof("initializing subsystem %s", name)
		go func(s Subsystem, wg *sync.WaitGroup, errChan chan<- bool) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					klog.Errorf("subsystem %s panicked whilst initializing: %v", s.Name(), r)
				}
			}()

			for i := 1; i <= 3; i++ {
				if e := s.Initialize(reg); e != nil {
					klog.Errorf("unable to initialize subsystem %s (error: %s) %d times. Waiting 10 seconds to retry", s.Name(), e.Error(), i)
					time.Sleep(time.Second * 10) // Wait for 10 seconds to try and reinitialize
					continue
				} else {
					return
				}
			}

			klog.Errorf("3 retries attempted to initialize subsystem %s, all failed", s.Name())
			errChan <- true
		}(sys, wg, errChan)
	}

	// Wait for all subsystems to complete, then evauate if any
	// failed by checking the channel for any booleans.
	wg.Wait()

	// If we get more than 1 error, then we wish to shutdown the process.
	if len(errChan) > 0 {
		m.shutdownSubsystems()
		os.Exit(len(errChan)) // Kind of clever maybe, return code is number of subsystems that failed.
	}
}

func (m *APIManager) reloadSubsystems() {
	klog.V(4).Infof("reload signal received, forwarding to %d subsystems", len(m.systems))

	wg := new(sync.WaitGroup)
	wg.Add(len(m.systems))

	for name, sys := range m.systems {
		klog.V(5).Infof("sending reload update for subsystem %s", name)
		go func(sys Subsystem, wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					klog.Errorf("subsystem %s panicked whilst reloading: %v", sys.Name(), r)
				}
			}()

			sys.Reload()
		}(sys, wg)
	}

	wg.Wait()
	klog.V(5).Info("reloading of subsystems complete")
}

func (m *APIManager) shutdownSubsystems() {
	klog.V(4).Infof("shutdown signal received, forwarding to %d subsystems", len(m.systems))

	wg := new(sync.WaitGroup)
	wg.Add(len(m.systems))

	for name, sys := range m.systems {
		klog.V(5).Infof("sending shutdown update for subsystem %s", name)
		go func(sys Subsystem, wg *sync.WaitGroup) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					klog.Errorf("subsystem %s panicked whilst shutting down: %v", sys.Name(), r)
				}
			}()

			sys.Shutdown()
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
