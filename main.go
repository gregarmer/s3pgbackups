package main

import (
	"flag"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var verbose = flag.Bool("verbose", false, "be verbose")
var noop = flag.Bool("noop", false, "don't actually do anything, just print what would be done")

func init() {
	flag.BoolVar(verbose, "v", false, "be verbose")
	flag.BoolVar(noop, "n", false, "don't actually do anything, just print what would be done")
}

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

	// do they work ?
}

func GetDatabases() []string {
	return []string{"foo", "bar"}
}

func GetTables(database string) []string {
	return []string{"baz", "quux", "django_session"}
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
	res, err := bucket.List("", "", "", 5)
	if err != nil {
		Fatalf("%s", err)
	}
	for _, v := range res.Contents {
		log.Printf(v.Key)
	}

	// back up the databases
	for _, db := range GetDatabases() {
		if config.ShouldExcludeDb(db) {
			log.Printf("[database] skipping '%s' because it's in excludes", db)
		} else {
			log.Printf("[database] backing up %s", db)
			for _, table := range GetTables(db) {
				if config.ShouldExcludeTable(table) {
					log.Printf("[table] skipping '%s' because it's in excludes", table)
				} else {
					log.Printf("[table] backing up %s", table)
				}
			}
		}
	}

	// compress backups

	// send to s3

	// rotate old s3 backups

	log.Printf("done - took %s", time.Since(start_time))
}
