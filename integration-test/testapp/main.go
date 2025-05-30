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

package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/fire833/go-api-utils/mgr"
	"k8s.io/klog/v2"
)

func main() {
	var flag flag.FlagSet
	klog.InitFlags(&flag)
	flag.Set("v", strconv.Itoa(10))

	m := mgr.New(&mgr.APIManagerOpts{
		EnableSysAPI: true,
		EnableVault:  false,
	})

	m.Initialize(&mgr.SystemRegistrar{
		AppName:      "testapp",
		Systems:      []mgr.Subsystem{},
		Registration: nil,
	})

	sec1 := mgr.NewSecretVaultValue("data", "Get secret string from vault", "foobad", "kttools/kv", "internal/ktnotify/api_key_hash")
	sec2 := mgr.NewSecretVaultValue("data", "Get secret string from vault", "foobad", "kttools/kv", "oidc/jupyter/client_secret")
	sec3 := mgr.NewSecretVaultValue("data", "Get secret string from vault", "foobad", "kttools/kv", "internal/github/webhooks/secret")
	sec4 := mgr.NewSecretVaultValue("data", "Get secret string from vault", "foobad", "kttools/kv", "internal/gitea/webhooks/secret")

	fmt.Printf("found: %s, %s, %s, %s\n", sec1.GetString(), sec2.GetString(), sec3.GetString(), sec4.GetString())
	m.SyncStartProcess()
}
