package config

import "testing"

func TestPreFlight(t *testing.T) {
	conf := Config{}
	err := conf.PreFlight()
	if err == nil {
		t.Fatalf("error should be set")
	}
}
