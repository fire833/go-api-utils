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

package gormsql

import (
	"sync"

	manager "github.com/fire833/go-api-utils/mgr"
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const GormSQLSubsystemName = "gormsql"

var SQL *GormSQLManager

type GormSQLManager struct {
	manager.DefaultSubsystem

	totalTransactions prometheus.Counter

	db *gorm.DB
}

func New() *GormSQLManager {
	return &GormSQLManager{
		db: nil,
	}
}

func (g *GormSQLManager) Initialize(wg *sync.WaitGroup, reg *manager.SystemRegistrar) error {
	defer wg.Done()

	g.totalTransactions = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: reg.AppName,
		Subsystem: GormSQLSubsystemName,
		Name:      "total_transactions",
		Help:      "Metrics on the total number of transactions made with this subsystem.",
	})

	config := &gorm.Config{
		Logger: &gormLogger{},
	}

	switch gormSQLBackend.GetString() {
	case "postgres", "POSTGRES", "Postgres":
		{
			if db, e := gorm.Open(postgres.Open(""), config); e == nil {
				g.db = db
			} else {
				return e
			}
		}
	case "mysql", "MYSQL", "MySQL":
		{
			if db, e := gorm.Open(mysql.Open(""), config); e == nil {
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

func (g *GormSQLManager) Name() string { return GormSQLSubsystemName }

func (g *GormSQLManager) SetGlobal() { SQL = g }

func (g *GormSQLManager) Collect(ch chan<- prometheus.Metric) {
}

func (g *GormSQLManager) Reload(wg *sync.WaitGroup) {
	defer wg.Done()
}

func (g *GormSQLManager) Shutdown(wg *sync.WaitGroup) {
	defer wg.Done()
}
