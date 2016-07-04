package database

import (
	"github.com/gregarmer/s3pgbackups/config"
	"testing"
)

var testConfig = config.Config{
	AwsAccessKey:     "access",
	AwsSecretKey:     "secret",
	PostgresUsername: "foo",
	PostgresPassword: "bar",
}

func TestGetDSNWithSSL(t *testing.T) {
	testConfig.PostgresSSL = true
	postgres := Postgres{Config: &testConfig}
	expectedDsn := "user=foo password=bar dbname=postgres sslmode=require"
	actualDsn := postgres.GetDSN()
	if expectedDsn != actualDsn {
		t.Fatalf("GetDSN with SSL failed. Expected '%s' got '%s'",
			expectedDsn, actualDsn)
	}
}

func TestGetDSNWithoutSSL(t *testing.T) {
	testConfig.PostgresSSL = false
	postgres := Postgres{Config: &testConfig}
	expectedDsn := "user=foo password=bar dbname=postgres sslmode=disable"
	actualDsn := postgres.GetDSN()
	if expectedDsn != actualDsn {
		t.Fatalf("GetDSN without SSL failed. Expected '%s' got '%s'",
			expectedDsn, actualDsn)
	}
}
