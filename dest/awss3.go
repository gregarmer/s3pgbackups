package dest

import (
	"github.com/dustin/go-humanize"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"github.com/gregarmer/s3pgbackups/config"
	"github.com/gregarmer/s3pgbackups/utils"
	"log"
	"os"
	"path/filepath"
)

type AwsS3 struct {
	Config *config.Config
}

func (awsS3 AwsS3) GetAuth() aws.Auth {
	// setup aws auth
	auth := aws.Auth{}
	auth.AccessKey = awsS3.Config.AwsAccessKey
	auth.SecretKey = awsS3.Config.AwsSecretKey
	return auth
}

func (awsS3 AwsS3) GetOrCreateBucket(noop *bool) *s3.Bucket {
	auth := awsS3.GetAuth()

	s := s3.New(auth, aws.USEast)
	bucket := s.Bucket(awsS3.Config.S3Bucket)
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
	return bucket
}

func (awsS3 AwsS3) DeleteFile(fileName string, noop *bool) {
	if *noop {
		log.Printf("would delete %s (noop)", fileName)
	} else {
		log.Printf("deleting %s", fileName)
		err := awsS3.GetOrCreateBucket(noop).Del(fileName)
		utils.CheckErr(err)
	}
}

func (awsS3 AwsS3) UploadTree(path string, noop *bool) {
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
				err = awsS3.GetOrCreateBucket(noop).PutReader("daily/"+fi.Name(),
					file, fi.Size(), "application/x-gzip",
					s3.BucketOwnerFull, s3.Options{})
				utils.CheckErr(err)
			}
		}
		return nil
	})
}

func (awsS3 AwsS3) RotateBackups(noop *bool) {
	// We keep 1 backup per day for the last week, 1 backup per week for the
	//   last month, and 1 backup per month indefinitely.
	res, err := awsS3.GetOrCreateBucket(noop).List("", "", "", 1000)
	utils.CheckErr(err)

	// 1. Come up with a set of required dates

	// 2. Run symmetric difference against all dates from list

	// 3. Delete remainder
	for _, v := range res.Contents {
		awsS3.DeleteFile(v.Key, noop)
	}
}
