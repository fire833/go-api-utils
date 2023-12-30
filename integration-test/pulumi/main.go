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
	"fmt"

	"github.com/pulumi/pulumi-vault/sdk/v5/go/vault"
	"github.com/pulumi/pulumi-vault/sdk/v5/go/vault/kv"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		kvMnt, e := vault.NewMount(ctx, fmt.Sprintf("%sKv", ctx.Stack()), &vault.MountArgs{
			Description:        pulumi.String(fmt.Sprintf("%s K/V store.", ctx.Stack())),
			Type:               pulumi.String("kv"),
			Path:               pulumi.String("devkv"),
			MaxLeaseTtlSeconds: pulumi.Int(3600),
			Options: pulumi.Map{
				"version": pulumi.Any("2"),
			},
		})

		if e != nil {
			return e
		}

		_, e = kv.NewSecretBackendV2(ctx, fmt.Sprintf("%sKv", ctx.Stack()), &kv.SecretBackendV2Args{
			Mount:              kvMnt.Path,
			MaxVersions:        pulumi.Int(25),
			DeleteVersionAfter: pulumi.Int(86400),
			CasRequired:        pulumi.Bool(false),
		})

		if e != nil {
			return e
		}

		return nil
	})
}
