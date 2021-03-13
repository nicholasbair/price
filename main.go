package main

import (
	"fmt"
	"price/config"
	"price/handler"
	"price/store"
	"strings"
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

	for _, inst := range strings.Split(config.GetEnv("INSTRUMENTS"), ",") {
		handler.PriceStreamHandler(db, config.GetEnv("ACCOUNT"), inst, config.GetEnv("TOKEN"))
	}

	<-wait
}
