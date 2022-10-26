package main

import (
	"final-project/database"
	"final-project/router"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	PORT := os.Getenv("PORT")
	database.StartDB()
	r := router.StartApp()
	r.Run(PORT)
}
