package main

import (
	"database/sql"
	"fmt"
	sw "llil.gq/go"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, port, user, password, dbname string) {
	var psqlInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=3",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	a.DB = db
}

func (a *App) Run(addr string) {
	log.Printf("Server started")

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(addr, router))
}

func (a *App) ensureTableExists() {
	tableCreationQuery := "CREATE TABLE urls(longURL VARCHAR NOT NULL UNIQUE, shortURL VARCHAR NOT NULL UNIQUE);"
	a.DB.Exec(tableCreationQuery)
}
