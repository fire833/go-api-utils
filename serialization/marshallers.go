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

package serialization

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/BurntSushi/toml"
	object "github.com/fire833/go-api-utils/object"
	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/proto"
	"gopkg.in/yaml.v3"
)

// Default unmarshaller for unmarshalling request bodies based on thier Content-Type header.
// This API will only support yaml and json as the content mediums.
func UnmarshalBodyByContentHeader(ctx *fasthttp.RequestCtx, data object.Object) error {
	switch string(ctx.Request.Header.Peek("Content-Type")) {
	case "application/yaml":
		{
			if e := yaml.Unmarshal(ctx.Request.Body(), data); e != nil {
				return e
			} else {
				return nil
			}
		}
	case "application/toml":
		{
			if e := toml.Unmarshal(ctx.Request.Body(), data); e != nil {
				return e
			} else {
				return nil
			}
		}
	case "application/xml":
		{
			if e := xml.Unmarshal(ctx.Request.Body(), data); e != nil {
				return e
			} else {
				return nil
			}
		}
	case "application/protobuf":
		{
			if e := proto.Unmarshal(ctx.Request.Body(), data); e != nil {
				return e
			} else {
				return nil
			}
		}
	default:
		{
			if e := json.Unmarshal(ctx.Request.Body(), data); e != nil {
				return e
			} else {
				return nil
			}
		}
	}
}

// Default marshaller to take interface and marshal it into the body of the response body.
// Will marshal to the correct format depending on the "Accept" header, defaults to json.
func MarshalBodyByAcceptHeader(ctx *fasthttp.RequestCtx, in object.Object) error {
	switch string(ctx.Request.Header.Peek("Accept")) {
	default:
		{
			if data, e := json.Marshal(&in); e != nil {
				InternalErrorResponseHandler(ctx, e.Error())
				return e
			} else {
				ctx.Response.SetBody(data)
				ctx.Response.SetStatusCode(http.StatusOK)
				return nil
			}
		}
	case "application/yaml":
		{
			if data, e := yaml.Marshal(&in); e != nil {
				InternalErrorResponseHandler(ctx, e.Error())
				return e
			} else {
				ctx.Response.SetBody(data)
				ctx.Response.SetStatusCode(http.StatusOK)
				return nil
			}
		}
	case "application/xml":
		{
			if data, e := xml.Marshal(&in); e != nil {
				InternalErrorResponseHandler(ctx, e.Error())
				return e
			} else {
				ctx.Response.SetBody(data)
				ctx.Response.SetStatusCode(http.StatusOK)
				return nil
			}
		}
	case "application/protobuf":
		{
			if data, e := proto.Marshal(in); e != nil {
				InternalErrorResponseHandler(ctx, e.Error())
				return e
			} else {
				ctx.Response.SetBody(data)
				ctx.Response.SetStatusCode(http.StatusOK)
				return nil
			}
		}
	}
}
