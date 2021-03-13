package main

import (
	"fmt"
	"price/config"
	"price/handler"
	"price/store"
)

func main() {
	config.LoadEnv()

	db := store.ConnectToDatabase()

	fmt.Println(db)

	setupErr := store.SetupTables(db)

	if setupErr != nil {
		fmt.Println(setupErr)
		panic("ERR setting up DB tables")
	}

	wait := make(chan bool)

	handler.PriceStreamHandler(db, "001-001-2420587-001", "EUR_USD")

	<-wait
}
