package swagger

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/stretchr/testify/assert"
)

func newTestDB(t *testing.T) *pg.DB {
	t.Helper()
	host := os.Getenv("APP_DB_HOST")
	port := os.Getenv("APP_DB_PORT")
	user := os.Getenv("APP_DB_USERNAME")
	password := os.Getenv("APP_DB_PASSWORD")
	dbname := os.Getenv("APP_DB_NAME")
	db := pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		User:     user,
		Password: password,
		Database: dbname,
	})

	t.Cleanup(func() {
		var shortUrlMaps []ShortUrlMap
		res, err := db.Model(&shortUrlMaps).Where("true").Delete()
		if err != nil {
			panic(err)
		}
		fmt.Println("Rows deleted: ", res.RowsAffected())

	})

	return db
}

func newTestDataShortenHandler(t *testing.T) *DataShortenHandler {
	t.Helper()
	return &DataShortenHandler{
		db:      newTestDB(t),
		baseUrl: os.Getenv("APP_BASE_URL"),
	}
}

func TestCreateNewShortURL_generateShortUrl(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		h := newTestDataShortenHandler(t)
		res := h.generateShortUrl(longUrlPayload{
			LongURL: "longurl",
		})
		assert.Len(t, res, 7)
		assert.Equal(t, "IjZEPut", res)
	})
	t.Run("Collision", func(t *testing.T) {
		h := newTestDataShortenHandler(t)

		addShortUrl(h.db, "url", "IjZEPut")
		res := h.generateShortUrl(longUrlPayload{
			LongURL: "longurl",
		})
		assert.Len(t, res, 7)
		assert.NotEqual(t, "IjZEPut", res)
	})
}
