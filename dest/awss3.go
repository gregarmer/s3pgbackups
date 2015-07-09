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
	res, err := bucket.List("", "", "", 1000)
	utils.CheckErr(err)

	// 1. Come up with a set of required dates
	// lastSevenDays := map[string]
	// lastFourWeeks := map[string]
	// datesToKeep := map[string]

	// 2. Run symmetric difference against all dates from list
	// for list of files
	//   if file not in datesToKeep
	//     delete

	// 3. Delete remainder
	for _, v := range res.Contents {
		awsS3.DeleteFile(bucket, v.Key, noop)
	}
}
