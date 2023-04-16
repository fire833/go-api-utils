package manager

import (
	"os"
	"testing"

	"github.com/spf13/viper"
)

func MockManager() *APIManager {
	return &APIManager{
		count:     0,
		systems:   make(map[string]Subsystem),
		shutdown:  make(chan uint8),
		config:    viper.New(),
		secrets:   viper.New(),
		sigHandle: make(chan os.Signal),
	}
}

func TestAPIManager_initializeSubsystems(t *testing.T) {
	type fields struct {
		count     uint
		systems   map[string]Subsystem
		config    *viper.Viper
		secrets   *viper.Viper
		shutdown  chan uint8
		sigHandle chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
		reg    *SystemRegistrar
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &APIManager{
				count:     tt.fields.count,
				systems:   tt.fields.systems,
				config:    tt.fields.config,
				secrets:   tt.fields.secrets,
				shutdown:  tt.fields.shutdown,
				sigHandle: tt.fields.sigHandle,
			}
			m.initializeSubsystems(tt.reg)
		})
	}
}

func TestAPIManager_reloadSubsystems(t *testing.T) {
	type fields struct {
		count     uint
		systems   map[string]Subsystem
		config    *viper.Viper
		secrets   *viper.Viper
		shutdown  chan uint8
		sigHandle chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &APIManager{
				count:     tt.fields.count,
				systems:   tt.fields.systems,
				config:    tt.fields.config,
				secrets:   tt.fields.secrets,
				shutdown:  tt.fields.shutdown,
				sigHandle: tt.fields.sigHandle,
			}
			m.reloadSubsystems()
		})
	}
}

func TestAPIManager_shutdownSubsystems(t *testing.T) {
	type fields struct {
		count     uint
		systems   map[string]Subsystem
		config    *viper.Viper
		secrets   *viper.Viper
		shutdown  chan uint8
		sigHandle chan os.Signal
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &APIManager{
				count:     tt.fields.count,
				systems:   tt.fields.systems,
				config:    tt.fields.config,
				secrets:   tt.fields.secrets,
				shutdown:  tt.fields.shutdown,
				sigHandle: tt.fields.sigHandle,
			}
			m.shutdownSubsystems()
		})
	}
}
