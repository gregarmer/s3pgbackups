package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const workingDir = "temp"

var verbose = flag.Bool("v", false, "be verbose")
var noop = flag.Bool("n", false, "don't actually do anything, just print what would be done")

func Fatalf(error string, args ...interface{}) {
	criticalLog := log.New(os.Stderr, "", log.LstdFlags)
	criticalLog.Printf(error, args...)
	os.Exit(1)
}

func PreFlight(config *Config) {
	// are we Config.PostgresUser ?

	// do we have s3 config (keys, buckets etc) ?
	if config.AwsAccessKey == "" {
		Fatalf("missing AwsAccessKey, cannot continue")
	}

	if config.AwsSecretKey == "" {
		Fatalf("missing AwsSecretKey, cannot continue")
	}
}

func GetDatabases() []string {
	var databases []string

	db, err := sql.Open("postgres", "user=vagrant password=vagrant dbname=postgres sslmode=require")
	checkErr(err)

	log.Print(db)

	rows, err := db.Query("SELECT datname FROM pg_database")
	checkErr(err)

	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		checkErr(err)
		databases = append(databases, name)
	}

	err = rows.Err() // get any error encountered during iteration
	checkErr(err)

	return databases
}

// func GetTables(database string) []string {
// 	var tables []string

// 	db, err := sql.Open("postgres", fmt.Sprintf("user=vagrant password=vagrant dbname=%s sslmode=require", database))
// 	checkErr(err)

// 	rows, err := db.Query("SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
// 	checkErr(err)

// 	defer rows.Close()
// 	for rows.Next() {
// 		var name string
// 		err = rows.Scan(&name)
// 		checkErr(err)
// 		tables = append(tables, name)
// 	}

// 	err = rows.Err() // get any error encountered during iteration
// 	checkErr(err)

// 	return tables
// }

func checkErr(err error) {
	if err != nil {
		Fatalf("Error: %s", err)
	}
}

func main() {
	start_time := time.Now()

	flag.Parse()

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}

	log.Printf("starting postgres cluster backup")

	if *noop {
		log.Printf("running in no-op mode, no commands will actually be executed")
	}

	config := LoadConfig()
	log.Printf("config: %+v", config)

	// pre-flight check (s3 keys, access to postgres etc)
	PreFlight(config)

	// setup aws auth
	auth := aws.Auth{}
	auth.AccessKey = config.AwsAccessKey
	auth.SecretKey = config.AwsSecretKey

	// Open Bucket
	s := s3.New(auth, aws.USEast)
	bucket := s.Bucket(config.S3Bucket)

	// create a working directory to store the backups
	currentDir, _ := os.Getwd()
	fullWorkingDir := currentDir + "/" + workingDir
	if _, err := os.Stat(fullWorkingDir); !os.IsNotExist(err) {
		log.Printf("working directory already exists at %s, removing it", fullWorkingDir)
		os.RemoveAll(fullWorkingDir)
	}
	os.Mkdir(fullWorkingDir, 0700)

	// back up the databases
	for _, db := range GetDatabases() {
		if config.ShouldExcludeDb(db) {
			log.Printf("[database] skipping '%s' because it's in excludes", db)
		} else {
			log.Printf("[%s] backing up database", db)

			// create backup
			backupFileName := fmt.Sprintf("%s-%s.sql", db, time.Now().Format("2006-01-02"))
			pgDumpCmd := fmt.Sprintf("-E UTF-8 -T %s -f %s %s",
				strings.Join(config.PostgresExcludeTable, " -T "),
				fmt.Sprintf("%s/%s", fullWorkingDir, backupFileName),
				db)
			log.Printf("executing pg_dump %s", pgDumpCmd)
			cmd := exec.Command("/usr/bin/pg_dump", strings.Split(pgDumpCmd, " ")...)
			var out bytes.Buffer
			cmd.Stdout = &out
			err := cmd.Run()
			checkErr(err)
			// fmt.Printf("out: %q\n", out.String())

			// compress backup
			log.Printf("compressing %s", backupFileName)
		}
	}

	// create bucket incase it doesn't already exist
	log.Printf("creating bucket %s", config.S3Bucket)
	err := bucket.PutBucket(s3.BucketOwnerFull)
	checkErr(err)

	// walk temp and upload everything to S3
	filepath.Walk(fullWorkingDir, func(localFile string, fi os.FileInfo, err error) (e error) {
		if !fi.IsDir() {
			log.Printf("uploading %s (%s)", localFile, humanize.Bytes(uint64(fi.Size())))

			file, err := os.Open(localFile)
			checkErr(err)
			defer file.Close()

			if *noop {
				log.Printf("would upload %s (%s)", fi.Name(), humanize.Bytes(uint64(fi.Size())))
			} else {
				err = bucket.PutReader("daily/"+fi.Name(), file, fi.Size(), "application/x-gzip", s3.BucketOwnerFull, s3.Options{})
				checkErr(err)
			}
		}
		return nil
	})

	// cleanup working directory
	os.RemoveAll(fullWorkingDir)

	// rotate old s3 backups
	// We keey 1 backup per day for the last week, 1 backup per week for the
	//   last month, and 1 backup per month indefinitely.
	log.Printf("rotating old backups")
	res, err := bucket.List("", "", "", 5)
	checkErr(err)

	for _, v := range res.Contents {
		log.Printf("deleting %s", v.Key)
	}

	log.Printf("done - took %s", time.Since(start_time))
}
