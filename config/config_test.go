package config

import "testing"

func TestPreFlight(t *testing.T) {
	conf := Config{}
	err := conf.PreFlight()
	if err == nil {
		t.Fatalf("error should be set")
	}
}

func TestCopy(t *testing.T) {
	conf := Config{S3Bucket: "foo"}
	c := conf.Copy()
	c.S3Bucket = "bar"
	if conf.S3Bucket == c.S3Bucket {
		t.Fatalf("config.Copy() should return a copy that doesn't affect the original")
	}
}
