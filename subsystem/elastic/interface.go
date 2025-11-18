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

package elastic

import (
	"errors"
	"os"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v9"
	manager "github.com/fire833/go-api-utils/mgr"
	"github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
)

const ElasticSubsystemName string = "ELASTIC"

var ELASTIC *ElasticManager

type ElasticManager struct {
	manager.DefaultSubsystem

	client *elasticsearch.Client

	credsRenewer *api.LifetimeWatcher
	creds        *api.Secret

	isInitialized bool
	isShutdown    bool
}

func New() *ElasticManager {
	return &ElasticManager{
		client:        nil,
		isInitialized: false,
		isShutdown:    false,
	}
}

func NewWithCreds(creds *api.Secret, watcher *api.LifetimeWatcher) *ElasticManager {
	return &ElasticManager{
		client:        nil,
		creds:         creds,
		credsRenewer:  watcher,
		isInitialized: false,
		isShutdown:    false,
	}
}

func (s *ElasticManager) Name() string { return ElasticSubsystemName }

func (s *ElasticManager) SetGlobal() { ELASTIC = s }

func (s *ElasticManager) Initialize(reg *manager.SystemRegistrar) error {
	var user, pass string
	if s.creds != nil { // Try and read in credentials from the vault API secret.
		u, okUser := s.creds.Data["username"]
		p, okPass := s.creds.Data["password"]

		if !okUser || !okPass {
			return errors.New("elastic creds from vault couldn't be retrieved, keys do not exist in secret")
		}

		user = u.(string)
		pass = p.(string)
	} else { // Otherwise, try from envvars
		user = os.Getenv("ELASTIC_USER")
		pass = os.Getenv("ELASTIC_PASS")
	}

	client, e := elasticsearch.NewClient(elasticsearch.Config{
		Username: user,
		Password: pass,
		Logger:   &elastictransport.JSONLogger{},
	})
	if e != nil {
		return e
	}

	s.client = client

	s.isInitialized = true
	return nil
}

// NOP PreInit
func (s *ElasticManager) PreInit() {}

// NOP SyncStart
func (s *ElasticManager) SyncStart() {}

// NOP PostInit
func (s *ElasticManager) PostInit() {}

func (s *ElasticManager) Configs() *[]*manager.ConfigValue {
	return &[]*manager.ConfigValue{}
}

func (s *ElasticManager) Secrets() *[]*manager.SecretValue {
	return &[]*manager.SecretValue{}
}

// NOP to reload the subsystem
func (s *ElasticManager) Reload() {}

// NOP to shutdown the subsystem
func (s *ElasticManager) Shutdown() {
	s.isShutdown = true
}

// Return nothing since this subsystem does nothing, but you should be able to fill this out at runtime
// so the APIManager can effectively make decision on process lifecycle.
func (s *ElasticManager) Status() *manager.SubsystemStatus {
	return &manager.SubsystemStatus{
		Name:          ElasticSubsystemName,
		IsInitialized: s.isInitialized,
		IsShutdown:    s.isShutdown,
		Meta:          nil,
	}
}

// NOP to implement prometheus Collector interface.
func (s *ElasticManager) Describe(chan<- *prometheus.Desc) {}

// NOP to implement prometheus Collector interface.
func (s *ElasticManager) Collect(chan<- prometheus.Metric) {}
