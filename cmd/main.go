package main

import (
	"github.com/ivanbaug/go-elecshops/internal/dbdriver"
	"github.com/ivanbaug/go-elecshops/internal/routes"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := dbdriver.ConnectSQL(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Cannot connect to database! Dying...")
	}
	defer db.SQL.Close()

	r := routes.SetupRouter(db)
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
