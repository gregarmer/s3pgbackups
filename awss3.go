package main

import (
	//"bytes"
	//"flag"
	//"fmt"
	//"github.com/dustin/go-humanize"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	//_ "github.com/lib/pq"
	//"io/ioutil"
	"log"
	//"os"
	//"os/exec"
	//"path/filepath"
	//"strings"
	//"time"
)

type AwsS3 struct {
	config *Config
}

func (awsS3 AwsS3) GetAuth() aws.Auth {
	// setup aws auth
	auth := aws.Auth{}
	auth.AccessKey = awsS3.config.AwsAccessKey
	auth.SecretKey = awsS3.config.AwsSecretKey
	return auth
}

func (awsS3 AwsS3) GetBucket() *s3.Bucket {
	auth := awsS3.GetAuth()

	s := s3.New(auth, aws.USEast)
	bucket := s.Bucket(awsS3.config.S3Bucket)
	return bucket
}

func (awsS3 AwsS3) RotateBackups() {
	// We keep 1 backup per day for the last week, 1 backup per week for the
	//   last month, and 1 backup per month indefinitely.
	res, err := awsS3.GetBucket().List("", "", "", 5)
	checkErr(err)

	for _, v := range res.Contents {
		log.Printf("deleting %s", v.Key)
	}
}
