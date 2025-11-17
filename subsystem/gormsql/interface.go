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

package gormsql

import (
	"errors"
	"fmt"
	"os"

	manager "github.com/fire833/go-api-utils/mgr"
	"github.com/hashicorp/vault/api"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"k8s.io/klog/v2"
)

const GormSQLSubsystemName = "gormsql"

var SQL *GormSQLManager

type GormSQLManager struct {
	manager.DefaultSubsystem

	totalTransactions prometheus.Counter

	credsRenewer *api.LifetimeWatcher
	creds        *api.Secret
	db           *gorm.DB
}

func New() *GormSQLManager {
	return &GormSQLManager{
		db:           nil,
		credsRenewer: nil,
	}
}

func NewWithCreds(creds *api.Secret, watcher *api.LifetimeWatcher) *GormSQLManager {
	return &GormSQLManager{
		db:           nil,
		creds:        creds,
		credsRenewer: watcher,
	}
}

func (g *GormSQLManager) Initialize(reg *manager.SystemRegistrar) error {
	g.totalTransactions = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: reg.AppName,
		Subsystem: GormSQLSubsystemName,
		Name:      "total_transactions",
		Help:      "Metrics on the total number of transactions made with this subsystem.",
	})

	config := &gorm.Config{
		Logger: &gormLogger{},
	}

	var user, pass string
	if g.creds != nil { // Try and read in credentials from the vault API secret.
		u, okUser := g.creds.Data["username"]
		p, okPass := g.creds.Data["password"]

		if !okUser || !okPass {
			return errors.New("gorm creds from vault couldn't be retrieved, keys do not exist in secret")
		}

		user = u.(string)
		pass = p.(string)
	} else { // Otherwise, try from envvars
		user = os.Getenv("GORM_SQL_USER")
		pass = os.Getenv("GORM_SQL_PASS")
	}

	switch gormSQLBackend.GetString() {
	case "postgres", "POSTGRES", "Postgres", "pg", "cockroach":
		{
			if db, e := gorm.Open(postgres.Open(g.createPostgresConnstring(user, pass)), config); e == nil {
				g.db = db
			} else {
				return e
			}
		}
	case "mysql", "MYSQL", "MySQL":
		{
			if db, e := gorm.Open(mysql.Open(g.createMysqlConnstring(user, pass)), config); e == nil {
				g.db = db
			} else {
				return e
			}
		}
	default:
		{
			if db, e := gorm.Open(sqlite.Open(gormSqliteFile.GetString()), config); e == nil {
				g.db = db
			} else {
				return e
			}
		}
	}

	return nil
}

func (g *GormSQLManager) createPostgresConnstring(user, pass string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		gormSqlHost.GetString(), user, pass, gormSqlDb.GetString(), gormSqlPort.GetUint16(),
		gormTlsverifyLevel.GetString())
}

func (g *GormSQLManager) createMysqlConnstring(user, pass string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, gormSqlHost.GetString(), gormSqlPort.GetUint16(), gormSqlDb.GetString())
}

func (g *GormSQLManager) Name() string { return GormSQLSubsystemName }

func (g *GormSQLManager) SetGlobal() { SQL = g }

func (g *GormSQLManager) Collect(ch chan<- prometheus.Metric) {
	ch <- g.totalTransactions
}

func (g *GormSQLManager) SyncStart() {
	if g.credsRenewer != nil {
		go g.credsRenewer.Start()
		defer g.credsRenewer.Stop()

		for {
			select {
			case done := <-g.credsRenewer.DoneCh():
				klog.Errorf("received error for db credential renewals: %s", done)
				return
			case renew := <-g.credsRenewer.RenewCh():
				klog.Infof("successfully renewed db credentials at %s for %d seconds, restarting connections", renew.RenewedAt, renew.Secret.LeaseDuration)
			}
		}
	}
}
