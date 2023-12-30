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
		EnableSysAPI: false,
		EnableVault:  true,
	})

	m.Initialize(&mgr.SystemRegistrar{
		AppName:      "testapp",
		Systems:      []mgr.Subsystem{},
		Registration: nil,
	})

	sec1 := mgr.NewSecretVaultValue("key1", "Get secret string from vault", "foobad", "devkv/test1")
	sec2 := mgr.NewSecretVaultValue("key2", "Get secret string from vault", "foobad", "devkv/test1")
	sec3 := mgr.NewSecretVaultValue("key1", "Get secret string from vault", "foobad", "devkv/test2")
	sec4 := mgr.NewSecretVaultValue("key2", "Get secret string from vault", "foobad", "devkv/test2")

	fmt.Printf("found: %s, %s, %v, %d\n", sec1.GetString(), sec2.GetString(), sec3.GetBool(), sec4.GetInt())
}
