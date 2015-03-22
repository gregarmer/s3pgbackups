package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
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
	PostgresUser         string   `json:"pg_user"`
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

func (c Config) ShouldExcludeDb(db string) bool {
	return _shouldExclude(db, c.PostgresExcludeDb)
}

func (c Config) ShouldExcludeTable(table string) bool {
	return _shouldExclude(table, c.PostgresExcludeTable)
}

func GetConfigPath() string {
	u, _ := user.Current()
	return u.HomeDir + "/" + configFile
}

func LoadConfig() *Config {
	// init config if needed
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		InitConfig()
		Fatalf("couldn't find config, created empty config at %s, please configure", configPath)
	}

	// load config
	file, _ := os.Open(configPath)
	decoder := json.NewDecoder(file)
	config := Config{}
	if err := decoder.Decode(&config); err != nil {
		Fatalf("error: %s", err)
	}

	return &config
}

func InitConfig() {
	// initialize a new configuration in ~/.s3pgbackups
	configPath := GetConfigPath()

	config := Config{}
	fileJson, _ := json.MarshalIndent(config, "", "  ")
	if err := ioutil.WriteFile(configPath, fileJson, 0600); err != nil {
		Fatalf("error: %+v", err)
	}
}
