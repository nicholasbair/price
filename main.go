package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"price/handler"
	"price/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	db := store.ConnectToDatabase()

	setupErr := store.SetupTables(db)

	if setupErr != nil {
		fmt.Println(setupErr)
		panic("ERR setting up DB tables")
	}

	wait := make(chan bool)

	handler.PriceStreamHandler(db, "001-001-2420587-001", "EUR_USD")

	<-wait
}
