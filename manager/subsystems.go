package manager

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

	go m.startSysAPI() // start sysAPI.

	for name, sys := range m.systems {
		klog.V(4).Infof("synchronously starting subsystem %s", name)
		go sys.SyncStart()
	}

	<-m.shutdown
}

func (m *APIManager) initializeSubsystems(reg AppRegistration) {
	wg := new(sync.WaitGroup)
	wg.Add(int(m.count))

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

				wgInt := new(sync.WaitGroup)
				wgInt.Add(1)

				if e := s.Initialize(wgInt, reg); e != nil {
					klog.Errorf("unable to initialize subsystem %s %d times. Waiting 10 seconds to retry", s.Name(), i)
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

			sys.Reload(wg)
		}(sys, wg)

		sys.Reload(wg)
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

			sys.Shutdown(wg)
		}(sys, wg)
	}

	wg.Wait()
	klog.V(5).Info("shutdown of subsystems complete")

	if e := m.server.Shutdown(); e != nil {
		klog.Errorf("unable to gracefully shutdown sysAPI: %v", e)
	}
}

func (m *APIManager) setGlobals() {
	for _, sub := range m.systems {
		sub.SetGlobal()
	}
}
