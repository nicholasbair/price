package main

import (
	"log"
	"os"
	"price/config"
	"price/handler"
	"price/store"
)

func main() {
	config.LoadEnv()

	{
		file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)
	}

	db := store.ConnectToDatabase()

	setupErr := store.SetupTables(db)

	if setupErr != nil {
		log.Println(setupErr)
		panic("ERR setting up DB tables")
	}

	wait := make(chan bool)

	go handler.PriceStreamHandler(db, config.GetEnv("ACCOUNT"), config.GetEnv("INSTRUMENTS"), config.GetEnv("TOKEN"), 0)

	<-wait
}
