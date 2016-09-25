package dest

import (
	"../config"
	"testing"
)

func TestGetAuth(t *testing.T) {
	conf := config.Config{AwsAccessKey: "access", AwsSecretKey: "secret"}
	awsS3 := AwsS3{Config: &conf}
	auth := awsS3.GetAuth()
	if auth.AccessKey != "access" {
		t.Fatalf("auth.AccessKey should be 'access'")
	}
	if auth.SecretKey != "secret" {
		t.Fatalf("auth.SecretKey should be 'secret'")
	}
}
