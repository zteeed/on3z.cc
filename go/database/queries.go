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

func ListShortURLByUser(database Database, auth0Sub string, pageLength int, offset int) ([]ShortURLMappingRestrict, error) {
	result := make([]ShortURLMappingRestrict, 0)
	rows, err := database.StatementListShortURLMappingAuthenticated.Query(auth0Sub, pageLength, offset)
	var resultLine ShortURLMappingRestrict
	for rows.Next() {
		rows.Scan(&resultLine.ShortURL, &resultLine.LongURL)
		result = append(result, resultLine)
	}
	return result, err
}

func SelectShortURLByUserByShortURL(database Database, shortURL string, auth0Sub string) (ShortURLMappingRestrict, error) {
	result := ShortURLMappingRestrict{}
	err := database.StatementSelectShortURLMappingAuthenticated.QueryRow(shortURL, auth0Sub).Scan(&result.ShortURL, &result.LongURL)
	return result, err
}

func AddShortUrl(database Database, longURL string, shortURL string, auth0Sub string) {
	_, err := database.StatementAddShortURLMapping.Exec(shortURL, longURL, auth0Sub)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateShortUrl(database Database, longURL string, shortURL string, auth0Sub string) {
	_, err := database.StatementUpdateShortURLMapping.Exec(longURL, shortURL, auth0Sub)
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
