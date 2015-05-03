package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Postgres struct {
	config *Config
}

func (postgres Postgres) GetDatabases() []string {
	var databases []string

	// XXX: Parameterize this, or figure out a way to use ~/.pgpass instead
	db, err := sql.Open("postgres", "user=vagrant password=vagrant dbname=postgres sslmode=require")
	checkErr(err)

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
