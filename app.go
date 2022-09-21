package main

import (
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"llil.gq/go/auth0"
	"log"
	"net/http"
	"time"

	sw "llil.gq/go"
	"llil.gq/go/database"
)

type App struct {
	Router       *http.ServeMux
	Database     database.Database
	JwtValidator *validator.Validator
}

func (a *App) Initialize(baseUrl string) {
	db := database.InitializeSqlDB()
	err := db.Ping()
	if err != nil {
		time.Sleep(2 * time.Second)
		log.Printf(err.Error())
		a.Initialize(baseUrl)
	}
	dbObject := database.InitializeDatabase(db)
	a.Database = dbObject
	a.JwtValidator = auth0.GetJwtValidator()
	a.Router = sw.NewRouter(dbObject, a.JwtValidator, baseUrl)
}

func (a *App) Run(addr string) {
	log.Printf("Server started")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}
