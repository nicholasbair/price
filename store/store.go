package store

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"price/client"
	"price/config"
)

func ConnectToDatabase() *pg.DB {
	if config.IsBacktest() {
		db := pg.Connect(&pg.Options{
			User:     "postgres",
			Database: "postgres",
		})
		return db
	} else {
		opt, err := pg.ParseURL(config.GetEnv("DB_CONNECTION_STRING"))

		if err != nil {
			fmt.Println(err)
		}

		return pg.Connect(opt)
	}
}

func Insert(db *pg.DB, record *client.Price) error {
	_, err := db.Model(record).Insert()
	if err != nil {
		return err
	}
	return nil
}

func SetupTables(db *pg.DB) error {
	var p client.Price

	err := db.Model(&p).CreateTable(&orm.CreateTableOptions{
		Temp:        false,
		IfNotExists: true,
	})
	if err != nil {
		return err
	}

	return nil
}
