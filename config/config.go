package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gregarmer/s3pgbackups/utils"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

const configFile = ".s3pgbackups"

type Config struct {
	// AWS
	AwsAccessKey string `json:"aws_access_key"`
	AwsSecretKey string `json:"aws_secret_key"`

	// S3
	S3Bucket    string `json:"s3_bucket"`
	S3RotateOld bool   `json:"s3_rotate_old"`

	// PostgreSQL
	PostgresUsername string `json:"pg_username"`
	PostgresPassword string `json:"pg_password"`
	PostgresSSL      bool   `json:"pg_sslmode"`

	// New style excludes - see issue #1 - example:
	// ["database.*", "*.table", "database"]
	Excludes []string `json:"excludes"`
}

func (c *Config) Copy() Config {
	return *c
}

func (c *Config) ShouldExcludeDb(db string) bool {
	for _, cmp := range c.Excludes {
		if cmp == fmt.Sprintf("%s.*", db) || cmp == db {
			return true
		}
	}
	return false
}

func (c *Config) ShouldExcludeTable(db string, t string) bool {
	for _, cmp := range c.Excludes {
		// all tables are excluded for this db
		if cmp == fmt.Sprintf("%s.*", db) {
			return true
		}

		// table match is excluded for this db
		if cmp == fmt.Sprintf("%s.%s", db, t) || cmp == fmt.Sprintf("*.%s", t) {
			return true
		}
	}
	return false
}

func (c *Config) PreFlight() error {
	// do we have s3 config (keys, buckets etc) ?
	if c.AwsAccessKey == "" {
		return errors.New("missing AwsAccessKey, cannot continue")
	}

	if c.AwsSecretKey == "" {
		return errors.New("missing AwsSecretKey, cannot continue")
	}

	return nil
}

func GetConfigPath() string {
	u, _ := user.Current()
	return filepath.Join(u.HomeDir, configFile)
}

func LoadConfig(configFile string) *Config {
	// make sure the config file actually exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		utils.Fatalf("couldn't load config from %s", configFile)
	}

	// init config if needed
	defaultConfigPath := GetConfigPath()
	if configFile == defaultConfigPath {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			InitConfig()
			utils.Fatalf("created empty config at %s, please configure", configFile)
		}
	}

	// load config
	file, _ := os.Open(configFile)
	decoder := json.NewDecoder(file)
	config := Config{}
	if err := decoder.Decode(&config); err != nil {
		utils.Fatalf("error: %s", err)
	}

	// pre-flight check (s3 keys, access to postgres etc)
	utils.CheckErr(config.PreFlight())

	return &config
}

func InitConfig() {
	// initialize a new configuration in ~/.s3pgbackups
	configPath := GetConfigPath()

	config := Config{}
	fileJson, _ := json.MarshalIndent(config, "", "  ")
	if err := ioutil.WriteFile(configPath, fileJson, 0600); err != nil {
		utils.Fatalf("error: %+v", err)
	}
}
