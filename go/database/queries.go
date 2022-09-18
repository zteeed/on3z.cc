package database

import (
	"fmt"
	"log"
)

func SelectShortURL(database Database, shortURL string) (ShortURLMapping, error) {
	result := ShortURLMapping{}
	err := database.StatementSelectShortURLMapping.QueryRow(shortURL).Scan(&result.ShortURL, &result.LongURL, &result.Auth0Sub)
	return result, err
}

func SelectShortURLByUser(database Database, shortURL string, auth0Sub string) (ShortURLMapping, error) {
	result := ShortURLMapping{}
	err := database.StatementSelectShortURLMappingAuthenticated.QueryRow(shortURL, auth0Sub).Scan(&result.ShortURL, &result.LongURL, &result.Auth0Sub)
	return result, err
}

func AddShortUrl(database Database, longURL string, shortURL string, auth0Sub string) {
	_, err := database.StatementAddShortURLMapping.Exec(shortURL, longURL, auth0Sub)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteAll(database Database) {
	result, err := database.StatementDeleteAll.Exec()
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, _ := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Rows deleted: ", rowsAffected)
}
