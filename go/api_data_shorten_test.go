package swagger

import (
	"github.com/stretchr/testify/assert"
	"llil.gq/go/database"
	"os"
	"testing"
)

func newTestDB(t *testing.T) database.Database {
	t.Helper()
	db := database.InitializeSqlDB()
	dbObject := database.InitializeDatabase(db)

	t.Cleanup(func() {
		database.DeleteAll(dbObject)
	})

	return database.InitializeDatabase(db)
}

func newTestDataShortenHandler(t *testing.T) *DataShortenHandler {
	t.Helper()
	return &DataShortenHandler{
		database: newTestDB(t),
		baseUrl:  os.Getenv("APP_BASE_URL"),
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

		database.AddShortUrl(h.database, "url", "IjZEPut")
		res := h.generateShortUrl(longUrlPayload{
			LongURL: "longurl",
		})
		assert.Len(t, res, 7)
		assert.NotEqual(t, "IjZEPut", res)
	})
}
