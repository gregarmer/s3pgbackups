package dest

import (
	"../config"
	"../utils"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type AwsS3 struct {
	Config       *config.Config
	bucketExists bool
}

func (awsS3 *AwsS3) GetAuth() aws.Auth {
	// setup aws auth
	auth := aws.Auth{}
	auth.AccessKey = awsS3.Config.AwsAccessKey
	auth.SecretKey = awsS3.Config.AwsSecretKey
	return auth
}

func (awsS3 *AwsS3) getOrCreateBucket(noop *bool) *s3.Bucket {
	auth := awsS3.GetAuth()

	s := s3.New(auth, aws.USEast)
	bucket := s.Bucket(awsS3.Config.S3Bucket)

	if !awsS3.bucketExists {
		exists, err := bucket.Exists("")
		utils.CheckErr(err)
		if !exists {
			if !*noop {
				log.Printf("creating bucket %s", awsS3.Config.S3Bucket)
				err := bucket.PutBucket(s3.BucketOwnerFull)
				utils.CheckErr(err)
			} else {
				log.Printf("would create bucket %s (noop)", awsS3.Config.S3Bucket)
			}
		}
		awsS3.bucketExists = true
	}

	return bucket
}

func (awsS3 *AwsS3) DeleteFile(bucket *s3.Bucket, fileName string, noop *bool) {
	if *noop {
		log.Printf("would delete %s (noop)", fileName)
	} else {
		log.Printf("deleting %s", fileName)
		err := bucket.Del(fileName)
		utils.CheckErr(err)
	}
}

func (awsS3 *AwsS3) UploadTree(path string, noop *bool) {
	bucket := awsS3.getOrCreateBucket(noop)
	filepath.Walk(path, func(localFile string, fi os.FileInfo, err error) (e error) {
		if !fi.IsDir() {
			file, err := os.Open(localFile)
			utils.CheckErr(err)
			defer file.Close()

			if *noop {
				log.Printf("would upload %s (%s) (noop)",
					fi.Name(), humanize.Bytes(uint64(fi.Size())))
			} else {
				log.Printf("uploading %s (%s)", localFile, humanize.Bytes(uint64(fi.Size())))
				err = bucket.PutReader("daily/"+fi.Name(), file, fi.Size(),
					"application/x-gzip", s3.BucketOwnerFull, s3.Options{})
				utils.CheckErr(err)
			}
		}
		return nil
	})
}

func (awsS3 *AwsS3) RotateBackups(noop *bool) {
	// We keep 1 backup per day for the last week, 1 backup per week for the
	//   last month, and 1 backup per month indefinitely.
	bucket := awsS3.getOrCreateBucket(noop)
	res, err := bucket.List("daily", "", "", 1000)
	utils.CheckErr(err)

	now := time.Now()
	toKeep := []string{}

	// keep the last 7 days worth of backups in daily/
	for i := 0; i <= 6; i++ {
		toKeep = append(toKeep, fmt.Sprintf("-%s.sql",
			now.Add(-time.Hour*24*time.Duration(i)).Format("2006-01-02")))
	}

	for _, v := range res.Contents {
		delete := true

		for _, keep := range toKeep {
			if strings.Contains(v.Key, keep) {
				delete = false
				break
			}
		}

		if delete {
			awsS3.DeleteFile(bucket, v.Key, noop)
		}
	}
}
