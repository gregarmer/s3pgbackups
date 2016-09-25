package config

import (
	"../utils"
	"encoding/json"
	"errors"
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
	PostgresUsername     string   `json:"pg_username"`
	PostgresPassword     string   `json:"pg_password"`
	PostgresSSL          bool     `json:"pg_sslmode"`
	PostgresExcludeDb    []string `json:"pg_exclude_dbs"`
	PostgresExcludeTable []string `json:"pg_exclude_tables"`
}

func _shouldExclude(item string, excludes []string) bool {
	for _, b := range excludes {
		if b == item {
			return true
		}
	}
	return false
}

func (c *Config) Copy() Config {
	return *c
}

func (c *Config) ShouldExcludeDb(db string) bool {
	return _shouldExclude(db, c.PostgresExcludeDb)
}

func (c *Config) ShouldExcludeTable(table string) bool {
	return _shouldExclude(table, c.PostgresExcludeTable)
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

func LoadConfig() *Config {
	// init config if needed
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		InitConfig()
		utils.Fatalf("couldn't find config, created empty config at %s, please configure", configPath)
	}

	// load config
	file, _ := os.Open(configPath)
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
