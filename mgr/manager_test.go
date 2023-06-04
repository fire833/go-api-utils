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

package mgr

import (
	"testing"

	"github.com/go-openapi/spec"
	"github.com/valyala/fasthttp"
)

func TestAPIManager_RegisterSysAPIHandler(t *testing.T) {
	type args struct {
		method     string
		path       string
		handler    fasthttp.RequestHandler
		swaggerdoc spec.PathItem
		schemas    []*spec.Schema
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = New()
			if err := RegisterSysAPIHandler(tt.args.method, tt.args.path, tt.args.handler, tt.args.swaggerdoc, tt.args.schemas...); (err != nil) != tt.wantErr {
				t.Errorf("APIManager.RegisterSysAPIHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
