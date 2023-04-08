package manager

import (
	"reflect"
	"sync"
	"testing"
)

func TestDefaultSubsystem_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "1",
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultSubsystem{}
			if got := d.Name(); got != tt.want {
				t.Errorf("DefaultSubsystem.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultSubsystem_Initialize(t *testing.T) {
	wg1 := new(sync.WaitGroup)
	wg1.Add(1)

	tests := []struct {
		name    string
		wg      *sync.WaitGroup
		reg     AppRegistration
		wantErr bool
	}{
		{
			name:    "1",
			wg:      wg1,
			wantErr: false,
			reg:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultSubsystem{}
			if err := d.Initialize(tt.wg, tt.reg); (err != nil) != tt.wantErr {
				t.Errorf("DefaultSubsystem.Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultSubsystem_Status(t *testing.T) {
	tests := []struct {
		name string
		want *SubsystemStatus
	}{
		{
			name: "1",
			want: &SubsystemStatus{
				// IsInitialized: true,
				Name: "default",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DefaultSubsystem{}
			if got := d.Status(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultSubsystem.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
