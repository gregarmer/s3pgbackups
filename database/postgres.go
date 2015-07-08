package database

import (
	"database/sql"
	"fmt"
	"github.com/gregarmer/s3pgbackups/config"
	"github.com/gregarmer/s3pgbackups/utils"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Config *config.Config
}

func (postgres Postgres) GetDatabases() []string {
	var databases []string
	var sslmode string

	if postgres.Config.PostgresSSL {
		sslmode = "require"
	} else {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=%s",
		postgres.Config.PostgresUsername,
		postgres.Config.PostgresPassword,
		sslmode)

	db, err := sql.Open("postgres", dsn)
	utils.CheckErr(err)

	rows, err := db.Query("SELECT datname FROM pg_database")
	utils.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		utils.CheckErr(err)
		databases = append(databases, name)
	}

	err = rows.Err() // get any error encountered during iteration
	utils.CheckErr(err)

	return databases
}
