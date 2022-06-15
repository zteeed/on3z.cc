package main

import (
	"os"
)

func main() {
	app := App{}
	app.Initialize(os.Getenv("APP_BASE_URL"))
	defer app.Database.DB.Close()
	app.Run("0.0.0.0:8888")
}
