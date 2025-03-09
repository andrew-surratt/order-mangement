package service

import (
	"reflect"
	"testing"
)

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			"GetConfig should return default config",
			&Config{Datapath: "data", Staticpath: "static", Basepath: "http://localhost:8080"},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if got := GetConfig(); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("GetConfig() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
