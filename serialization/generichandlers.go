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

package serialization

import (
	"net/http"

	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func newErrorResponse(code uint32, err string, description string) *GenericErrorResponse {
	return &GenericErrorResponse{
		Error:       err,
		Description: description,
		Code:        code,
		Timestamp:   timestamppb.Now(),
	}
}

func newOKResponse(code uint32, message string) *OKResponse {
	return &OKResponse{
		Code:      code,
		Message:   message,
		Timestamp: timestamppb.Now(),
	}
}

func newBadRequestResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), description)
}

func newNotFoundResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusNotFound, http.StatusText(http.StatusNotFound), description)
}

func newUnauthorizedResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized), description)
}

func newForbiddenResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusForbidden, http.StatusText(http.StatusForbidden), description)
}

func newInternalErrorResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), description)
}

func newNotAcceptableResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable), description)
}

func newMethodNotAllowedResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed), description)
}

func newNotImplementedResponse(description string) *GenericErrorResponse {
	return newErrorResponse(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented), description)
}

func OKResponseHandler(ctx *fasthttp.RequestCtx, code uint32, message string) {
	MarshalBodyByAcceptHeader(ctx, newOKResponse(code, message))
	ctx.SetStatusCode(int(code))
}

func BadRequestResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newBadRequestResponse(description))
	ctx.SetStatusCode(http.StatusBadRequest)
}

func NotFoundResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newNotFoundResponse(description))
	ctx.SetStatusCode(http.StatusNotFound)
}

func UnauthorizedResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newUnauthorizedResponse(description))
	ctx.SetStatusCode(http.StatusUnauthorized)
}

func ForbiddenResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newForbiddenResponse(description))
	ctx.SetStatusCode(http.StatusForbidden)
}

func InternalErrorResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newInternalErrorResponse(description))
	ctx.SetStatusCode(http.StatusInternalServerError)
}

func NotAcceptableResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newNotAcceptableResponse(description))
	ctx.SetStatusCode(http.StatusNotAcceptable)
}

func MethodNotAllowedResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newMethodNotAllowedResponse(description))
	ctx.SetStatusCode(http.StatusMethodNotAllowed)
}

func NotImplementedResponseHandler(ctx *fasthttp.RequestCtx, description string) {
	MarshalBodyByAcceptHeader(ctx, newNotImplementedResponse(description))
	ctx.SetStatusCode(http.StatusNotImplemented)
}

func GenericOKResponseHandler(ctx *fasthttp.RequestCtx) {
	OKResponseHandler(ctx, http.StatusOK, "operation successful")
}

func GenericBadRequestResponseHandler(ctx *fasthttp.RequestCtx) {
	BadRequestResponseHandler(ctx, "request was malformed and could not be processed")
}

func GenericNotFoundResponseHandler(ctx *fasthttp.RequestCtx) {
	NotFoundResponseHandler(ctx, "requested object was not found on the server")
}

func GenericUnauthorizedResponseHandler(ctx *fasthttp.RequestCtx) {
	UnauthorizedResponseHandler(ctx, "request did not have proper authentication credentials")
}

func GenericForbiddenResponseHandler(ctx *fasthttp.RequestCtx) {
	ForbiddenResponseHandler(ctx, "request is forbidden with your current credentials")
}

func GenericInternalErrorResponseHandler(ctx *fasthttp.RequestCtx) {
	InternalErrorResponseHandler(ctx, "request cannot be processed due to internal server error")
}

func GenericNotAcceptableResponseHandler(ctx *fasthttp.RequestCtx) {
	NotAcceptableResponseHandler(ctx, "request was not formmatted in an acceptable manner")
}

func GenericMethodNotAllowedResponseHandler(ctx *fasthttp.RequestCtx) {
	MethodNotAllowedResponseHandler(ctx, "request method is not allowed on this endpoint")
}

func GenericNotImplementedResponseHandler(ctx *fasthttp.RequestCtx) {
	NotImplementedResponseHandler(ctx, "request endpoint is not yet implemented")
}
