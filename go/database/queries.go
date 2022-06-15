package database

import (
	"database/sql"
	"fmt"
	"log"
)

func SelectShortURL(database Database, shortURL string) (ShortURLMapping, error) {
	result := ShortURLMapping{}
	err := database.StatementSelectShortURLMapping.QueryRow(shortURL).Scan(&result.ShortURL, &result.LongURL)
	return result, err
}

func AddShortUrl(database Database, longURL string, shortURL string) {
	_, err := database.StatementAddShortURLMapping.Exec(shortURL, longURL)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteAll(database Database) sql.Result {
	result, err := database.StatementDeleteAll.Exec()
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	fmt.Println("Rows deleted: ", rowsAffected)
	return result
}
