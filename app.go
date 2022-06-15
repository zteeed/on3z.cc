package main

import (
	"log"
	"net/http"
	"time"

	sw "llil.gq/go"
	"llil.gq/go/database"
)

type App struct {
	Router   *http.ServeMux
	Database database.Database
}

func (a *App) Initialize(baseUrl string) {
	db := database.InitializeSqlDB()
	err := db.Ping()
	if err != nil {
		time.Sleep(2 * time.Second)
		a.Initialize(baseUrl)
		log.Fatal(err)
	}
	dbObject := database.InitializeDatabase(db)
	a.Database = dbObject
	a.Router = sw.NewRouter(dbObject, baseUrl)
}

func (a *App) Run(addr string) {
	log.Printf("Server started")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
