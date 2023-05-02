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

package apiserver

import (
	"strconv"
	"testing"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

func BenchmarkRouter1(b *testing.B) {
	router := router.New()

	for i := 0; i < 250; i++ {
		router.GET("/"+strconv.Itoa(i), func(ctx *fasthttp.RequestCtx) {})
	}

	router.GET("/150/hsjdfhs/438yuhfsdf/24y8hfjsdhf/ghaw7rt2734sd/foo/bar/1234", func(ctx *fasthttp.RequestCtx) {})

	req := fasthttp.AcquireRequest()
	uri := fasthttp.AcquireURI()
	uri.SetPath("/150/hsjdfhs/438yuhfsdf/24y8hfjsdhf/ghaw7rt2734sd/foo/bar/1234")
	req.SetURI(uri)

	ctx := &fasthttp.RequestCtx{
		Request: *req,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.Handler(ctx)
	}
}
