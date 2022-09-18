package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func InitializeSqlDB() *sql.DB {
	host := os.Getenv("APP_DB_HOST")
	port := os.Getenv("APP_DB_PORT")
	user := os.Getenv("APP_DB_USERNAME")
	password := os.Getenv("APP_DB_PASSWORD")
	dbname := os.Getenv("APP_DB_NAME")
	var psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=3",
		host, port, user, password, dbname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func InitializeDatabase(db *sql.DB) Database {
	StatementSelectShortURLMapping, err := db.Prepare(
		"SELECT short_url, long_url, auth0_sub FROM short_url_maps WHERE short_url = $1 LIMIT 1;",
	)
	if err != nil {
		panic(err)
	}
	StatementSelectShortURLMappingAuthenticated, err := db.Prepare(
		"SELECT short_url, long_url, auth0_sub FROM short_url_maps WHERE short_url = $1 AND auth0_sub = $2 LIMIT 1;",
	)
	if err != nil {
		panic(err)
	}
	StatementAddShortURLMapping, err := db.Prepare(
		"INSERT INTO short_url_maps (short_url, long_url, auth0_sub) VALUES ($1, $2, $3);",
	)
	if err != nil {
		panic(err)
	}
	StatementDeleteAll, err := db.Prepare(
		"DELETE FROM short_url_maps;",
	)
	if err != nil {
		panic(err)
	}
	return Database{
		DB:                             db,
		StatementSelectShortURLMapping: StatementSelectShortURLMapping,
		StatementSelectShortURLMappingAuthenticated: StatementSelectShortURLMappingAuthenticated,
		StatementAddShortURLMapping:                 StatementAddShortURLMapping,
		StatementDeleteAll:                          StatementDeleteAll,
	}
}
