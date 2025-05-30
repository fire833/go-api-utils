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
	"testing"

	"github.com/fire833/go-api-utils/fake"
	"github.com/valyala/fasthttp"
)

// Mostly used to make sure we don't crash and the genericErrorHandlers are guaranteed to compile.

func TestOKResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			OKResponseHandler(ctx, 200, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() < 200 || ctx.Response.StatusCode() > 299 {
				t.Errorf("OkResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 200)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			OKResponseHandler(ctx, 200, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() < 200 || ctx.Response.StatusCode() > 299 {
				t.Errorf("OkResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 200)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			OKResponseHandler(ctx, 200, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() < 200 || ctx.Response.StatusCode() > 299 {
				t.Errorf("OkResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 200)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			OKResponseHandler(ctx, 200, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() < 200 || ctx.Response.StatusCode() > 299 {
				t.Errorf("OkResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 200)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			OKResponseHandler(ctx, 200, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() < 200 || ctx.Response.StatusCode() > 299 {
				t.Errorf("OkResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 200)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestBadRequestResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			BadRequestResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 400 {
				t.Errorf("BadRequestResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 400)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			BadRequestResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 400 {
				t.Errorf("BadRequestResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 400)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			BadRequestResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 400 {
				t.Errorf("BadRequestResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 400)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			BadRequestResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 400 {
				t.Errorf("BadRequestResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 400)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			BadRequestResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 400 {
				t.Errorf("BadRequestResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 400)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestNotFoundResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			NotFoundResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 404 {
				t.Errorf("NotFoundResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 404)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			NotFoundResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 404 {
				t.Errorf("NotFoundResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 404)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			NotFoundResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 404 {
				t.Errorf("NotFoundResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 404)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			NotFoundResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 404 {
				t.Errorf("NotFoundResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 404)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			NotFoundResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 404 {
				t.Errorf("NotFoundResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 404)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestUnauthorizedResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			UnauthorizedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 401 {
				t.Errorf("UnauthorizedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 401)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			UnauthorizedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 401 {
				t.Errorf("UnauthorizedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 401)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			UnauthorizedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 401 {
				t.Errorf("UnauthorizedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 401)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			UnauthorizedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 401 {
				t.Errorf("UnauthorizedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 401)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			UnauthorizedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 401 {
				t.Errorf("UnauthorizedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 401)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestForbiddenResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			ForbiddenResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 403 {
				t.Errorf("ForbiddenResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 403)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			ForbiddenResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 403 {
				t.Errorf("ForbiddenResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 403)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			ForbiddenResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 403 {
				t.Errorf("ForbiddenResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 403)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			ForbiddenResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 403 {
				t.Errorf("ForbiddenResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 403)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			ForbiddenResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 403 {
				t.Errorf("ForbiddenResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 403)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestInternalErrorResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			InternalErrorResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 500 {
				t.Errorf("InternalErrorResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 500)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			InternalErrorResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 500 {
				t.Errorf("InternalErrorResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 500)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			InternalErrorResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 500 {
				t.Errorf("InternalErrorResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 500)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			InternalErrorResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 500 {
				t.Errorf("InternalErrorResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 500)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			InternalErrorResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 500 {
				t.Errorf("InternalErrorResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 500)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestNotAcceptableResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			NotAcceptableResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 406 {
				t.Errorf("NotAcceptableResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 406)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			NotAcceptableResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 406 {
				t.Errorf("NotAcceptableResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 406)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			NotAcceptableResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 406 {
				t.Errorf("NotAcceptableResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 406)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			NotAcceptableResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 406 {
				t.Errorf("NotAcceptableResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 406)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			NotAcceptableResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 406 {
				t.Errorf("NotAcceptableResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 406)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestMethodNotAllowedResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			MethodNotAllowedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 405 {
				t.Errorf("MethodNotAllowedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 405)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			MethodNotAllowedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 405 {
				t.Errorf("MethodNotAllowedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 405)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			MethodNotAllowedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 405 {
				t.Errorf("MethodNotAllowedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 405)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			MethodNotAllowedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 405 {
				t.Errorf("MethodNotAllowedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 405)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			MethodNotAllowedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 405 {
				t.Errorf("MethodNotAllowedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 405)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestNotImplementedResponseHandler(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		// Test default marshalling first.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}

			NotImplementedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 501 {
				t.Errorf("NotImplementedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 501)
			}
			if i == 250 {
				t.Logf("default marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("protobuf", func(t *testing.T) {
		// Explicitly test protobuf marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/API+Protobuf")

			NotImplementedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 501 {
				t.Errorf("NotImplementedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 501)
			}
			if i == 250 {
				t.Logf("protobuf marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("xml", func(t *testing.T) {
		// Explicitly test xml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/xml")

			NotImplementedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 501 {
				t.Errorf("NotImplementedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 501)
			}
			if i == 250 {
				t.Logf("xml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("yaml", func(t *testing.T) {
		// Explicitly test yaml marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/yaml")

			NotImplementedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 501 {
				t.Errorf("NotImplementedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 501)
			}
			if i == 250 {
				t.Logf("yaml marshalling: %s", ctx.Response.Body())
			}
		}
	})

	t.Run("json", func(t *testing.T) {
		// Explicitly test json marshalling.
		for i := 0; i < 500; i++ {
			ctx := &fasthttp.RequestCtx{}
			ctx.Request.Header.Set("Accept", "application/json")

			NotImplementedResponseHandler(ctx, fake.FakeStringCeil(256))
			if ctx.Response.StatusCode() != 501 {
				t.Errorf("NotImplementedResponseHandler() = %d, want %d", ctx.Response.StatusCode(), 501)
			}
			if i == 250 {
				t.Logf("json marshalling: %s", ctx.Response.Body())
			}
		}
	})
}

func TestGenericOKResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericOKResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericOKResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericOKResponseHandler(&fasthttp.RequestCtx{})
	}
}

func TestGenericBadRequestResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericBadRequestResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericBadRequestResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericBadRequestResponseHandler(&fasthttp.RequestCtx{})
	}
}

func TestGenericNotFoundResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericNotFoundResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericNotFoundResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericNotFoundResponseHandler(&fasthttp.RequestCtx{})
	}
}

func TestGenericUnauthorizedResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericUnauthorizedResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericUnauthorizedResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericUnauthorizedResponseHandler(&fasthttp.RequestCtx{})
	}
}

func TestGenericForbiddenResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericForbiddenResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericForbiddenResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericForbiddenResponseHandler(&fasthttp.RequestCtx{})
	}
}

func TestGenericNotAcceptableResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericNotAcceptableResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericNotAcceptableResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericNotAcceptableResponseHandler(&fasthttp.RequestCtx{})
	}
}

func TestGenericMethodNotAllowedResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericMethodNotAllowedResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericMethodNotAllowedResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericMethodNotAllowedResponseHandler(&fasthttp.RequestCtx{})
	}
}

func TestGenericNotImplementedResponseHandler(t *testing.T) {
	for i := 0; i < 500; i++ {
		GenericNotImplementedResponseHandler(&fasthttp.RequestCtx{})
	}
}

func BenchmarkGenericNotImplementedResponseHandler(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenericNotImplementedResponseHandler(&fasthttp.RequestCtx{})
	}
}
