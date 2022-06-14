package main

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	sw "llil.gq/go"
	"log"
	"net/http"
)

type App struct {
	Router *http.ServeMux
	DB     *pg.DB
}

// createSchema creates database schema for ShortUrlMap model.
func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*sw.ShortUrlMap)(nil),
	}

	for _, model := range models {
		_, err := db.Model(model).Exists()
		if err != nil {
			err := db.Model(model).CreateTable(&orm.CreateTableOptions{})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *App) Initialize(host string, port string, user string, password string, dbname string, baseUrl string) {
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		User:     user,
		Password: password,
		Database: dbname,
	})
	defer db.Close()

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	a.Router = sw.NewRouter(db, baseUrl)
}

func (a *App) Run(addr string) {
	log.Printf("Server started")
	log.Fatal(http.ListenAndServe(addr, a.Router))
}