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

	setupErr := store.SetupTables(db)

	if setupErr != nil {
		fmt.Println(setupErr)
		panic("ERR setting up DB tables")
	}

	wait := make(chan bool)

	go handler.PriceStreamHandler(db, config.GetEnv("ACCOUNT"), config.GetEnv("INSTRUMENTS"), config.GetEnv("TOKEN"), 0)

	<-wait
}
