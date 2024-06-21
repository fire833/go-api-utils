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

package gormsql

import manager "github.com/fire833/go-api-utils/mgr"

var (
	gormSQLBackend *manager.ConfigValue = manager.NewConfigValue(
		"gormSQLbackend",
		"Specify the backend that you want to collect data from. Current valid values are sqlite, postgres, or mysql.",
		"sqlite",
	)

	gormSqliteFile *manager.ConfigValue = manager.NewConfigValue(
		"gormSqliteFile",
		"Specify the relative or absolute path to a sqlite database file to be read or created by your application. This value will only be read if gormSQLbackend is set to 'sqlite'.",
		"data.db",
	)

	gormSqlHost *manager.ConfigValue = manager.NewConfigValue(
		"gormSqlHost",
		"Specify the hostname of the remote SQL instance.",
		"localhost",
	)

	gormSqlPort *manager.ConfigValue = manager.NewConfigValue(
		"gormSqlPort",
		"Specify the port of the remote SQL instance.",
		uint16(3306),
	)

	gormSqlUsername *manager.SecretValue = manager.NewSecretValue(
		"gormSqlUsername",
		"Specify the username to connect to the remote SQL instance.",
		"",
	)

	gormSqlPassword *manager.SecretValue = manager.NewSecretValue(
		"gormSqlPassword",
		"Specify the password to connect to the remote SQL instance.",
		"",
	)
)

func (g *GormSQLManager) Configs() *[]*manager.ConfigValue {
	return &[]*manager.ConfigValue{
		gormSQLBackend,
		gormSqliteFile,
		gormSqlHost,
		gormSqlPort,
	}
}

func (g *GormSQLManager) Secrets() *[]*manager.SecretValue {
	return &[]*manager.SecretValue{
		gormSqlUsername,
		gormSqlPassword,
	}
}
