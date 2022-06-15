package swagger

import (
	"github.com/go-pg/pg/v10"
)

func selectShortURL(db *pg.DB, shortURL string) (bool, *ShortUrlMap) {
	shortUrlMap := new(ShortUrlMap)
	err := db.Model(shortUrlMap).Limit(1).Where("short_url = ?", shortURL).Select()
	shortUrlExist := err == nil
	return shortUrlExist, shortUrlMap
}

func addShortUrl(db *pg.DB, longURL string, shortURL string) {
	shortUrlMapInsert := &ShortUrlMap{
		LongURL:  longURL,
		ShortURL: shortURL,
	}
	_, err := db.Model(shortUrlMapInsert).Insert()
	if err != nil {
		panic(err)
	}
}
