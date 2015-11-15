package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/gregarmer/s3pgbackups/config"
	"github.com/gregarmer/s3pgbackups/utils"
	_ "github.com/lib/pq"
	"log"
	"os/exec"
	"strings"
	"time"
)

type Postgres struct {
	Config *config.Config
}

func (postgres *Postgres) GetDSN() string {
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

	return dsn
}

func (postgres *Postgres) GetDatabases() []string {
	var databases []string

	db, err := sql.Open("postgres", postgres.GetDSN())
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

func (postgres *Postgres) GetTables(database string) []string {
	var tables []string

	db, err := sql.Open(database, postgres.GetDSN())
	utils.CheckErr(err)

	tablesQuery := "SELECT table_name FROM information_schema.tables "
	tablesQuery += "WHERE table_schema = 'public'"
	rows, err := db.Query(tablesQuery)
	utils.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		utils.CheckErr(err)
		tables = append(tables, name)
	}

	err = rows.Err() // get any error encountered during iteration
	utils.CheckErr(err)

	return tables
}

func (postgres *Postgres) DumpDatabase(db, workingDir string) string {
	backupFileName := fmt.Sprintf("%s-%s.sql", db,
		time.Now().Format("2006-01-02"))

	pgDumpCmd := fmt.Sprintf("-E UTF-8 -f %s ",
		fmt.Sprintf("%s/%s", workingDir, backupFileName))

	if postgres.Config.Excludes != nil {
		for _, table := range postgres.GetTables(db) {
			if postgres.Config.ShouldExcludeTable(db, table) {
				pgDumpCmd += fmt.Sprintf("-T %s ", table)
			}
		}
	}

	pgDumpCmd += db

	log.Printf("executing pg_dump %s", pgDumpCmd)
	cmd := exec.Command("pg_dump", strings.Split(pgDumpCmd, " ")...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	utils.CheckErr(err)
	// fmt.Printf("out: %q\n", out.String())
	return backupFileName
}
