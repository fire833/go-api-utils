package manager

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
			m := New()
			if err := m.RegisterSysAPIHandler(tt.args.method, tt.args.path, tt.args.handler, tt.args.swaggerdoc, tt.args.schemas...); (err != nil) != tt.wantErr {
				t.Errorf("APIManager.RegisterSysAPIHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
