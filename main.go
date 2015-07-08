package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/goamz/goamz/s3"
	"github.com/gregarmer/s3pgbackups/config"
	"github.com/gregarmer/s3pgbackups/database"
	"github.com/gregarmer/s3pgbackups/dest"
	"github.com/gregarmer/s3pgbackups/utils"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const workingDir = "temp"

var verbose = flag.Bool("v", false, "be verbose")
var noop = flag.Bool("n", false,
	"don't actually do anything, just print what would be done")

func main() {
	start_time := time.Now()

	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	log.Printf("starting postgres cluster backup")

	if *noop {
		log.Printf("running in no-op mode, no commands will really be executed")
	}

	conf := config.LoadConfig()

	// Don't print real passwords and secret keys in verbose mode
	verbose_config := conf.Copy()
	verbose_config.PostgresPassword = "****"
	verbose_config.AwsSecretKey = "****"
	log.Printf("config: %+v", verbose_config)

	// AwsS3
	awsS3 := dest.AwsS3{conf}
	bucket := awsS3.GetBucket()

	// Postgres
	postgres := database.Postgres{conf}

	// create a working directory to store the backups
	currentDir, _ := os.Getwd()
	fullWorkingDir := currentDir + "/" + workingDir
	if _, err := os.Stat(fullWorkingDir); !os.IsNotExist(err) {
		log.Printf("working directory already exists at %s, removing it",
			fullWorkingDir)
		os.RemoveAll(fullWorkingDir)
	}
	os.Mkdir(fullWorkingDir, 0700)

	// back up the databases
	for _, db := range postgres.GetDatabases() {
		if conf.ShouldExcludeDb(db) {
			log.Printf("[database] skipping '%s' because it's in excludes", db)
		} else {
			log.Printf("[%s] backing up database", db)

			// create backup
			backupFileName := fmt.Sprintf("%s-%s.sql", db,
				time.Now().Format("2006-01-02"))
			pgDumpCmd := fmt.Sprintf("-E UTF-8 -T %s -f %s %s",
				strings.Join(conf.PostgresExcludeTable, " -T "),
				fmt.Sprintf("%s/%s", fullWorkingDir, backupFileName),
				db)
			log.Printf("executing pg_dump %s", pgDumpCmd)
			cmd := exec.Command("pg_dump", strings.Split(pgDumpCmd, " ")...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			utils.CheckErr(err)
			// fmt.Printf("out: %q\n", out.String())

			// compress backup
			log.Printf("compressing %s", backupFileName)
			cmd = exec.Command("gzip", fmt.Sprintf("%s/%s", fullWorkingDir,
				backupFileName))
			cmd.Stdout = &out
			err = cmd.Run()
			utils.CheckErr(err)
		}
	}

	// create bucket incase it doesn't already exist
	log.Printf("creating bucket %s", conf.S3Bucket)
	err := bucket.PutBucket(s3.BucketOwnerFull)
	utils.CheckErr(err)

	// walk temp and upload everything to S3
	awsS3.UploadTree(fullWorkingDir, noop)
	fmt.Println(bucket, s3.BucketOwnerFull)

	// cleanup working directory
	os.RemoveAll(fullWorkingDir)

	// rotate old s3 backups
	log.Printf("rotating old backups")
	awsS3.RotateBackups(noop)

	log.Printf("done - took %s", time.Since(start_time))
}
