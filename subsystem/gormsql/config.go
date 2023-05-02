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

import manager "github.com/fire833/go-api-utils/mgr"

var (
	gormSQLBackend *manager.ConfigKey = manager.NewConfigKey(
		"gormSQLbackend",
		"Specify the backend that you want to collect data from. Current valid values are sqlite, postgres, or mysql.",
		manager.String,
		"sqlite",
		false,
	)

	gormSqliteFile *manager.ConfigKey = manager.NewConfigKey(
		"gormSqliteFile",
		"Specify the relative or absolute path to a sqlite database file to be read or created by your application. This value will only be read if gormSQLbackend is set to 'sqlite'.",
		manager.String,
		"data.db",
		false,
	)

	gormSqlHost *manager.ConfigKey = manager.NewConfigKey(
		"gormSqlHost",
		"Specify the hostname of the remote SQL instance.",
		manager.String,
		"localhost",
		false,
	)

	gormSqlPort *manager.ConfigKey = manager.NewConfigKey(
		"gormSqlPort",
		"Specify the port of the remote SQL instance.",
		manager.Uint16,
		3306,
		false,
	)

	gormSqlUsername *manager.ConfigKey = manager.NewConfigKey(
		"gormSqlUsername",
		"Specify the username to connect to the remote SQL instance.",
		manager.String,
		nil,
		true,
	)

	gormSqlPassword *manager.ConfigKey = manager.NewConfigKey(
		"gormSqlPassword",
		"Specify the password to connect to the remote SQL instance.",
		manager.String,
		nil,
		true,
	)
)

func (g *GormSQLManager) Configs() *[]*manager.ConfigKey {
	return &[]*manager.ConfigKey{
		gormSQLBackend,
		gormSqliteFile,
		gormSqlHost,
		gormSqlPort,
		gormSqlUsername,
		gormSqlPassword,
	}
}
