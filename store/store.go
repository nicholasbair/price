package store

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"price/client"
)

func ConnectToDatabase() *pg.DB {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Database: "postgres",
	})
	return db
}

func Insert(db *pg.DB, record *client.PriceEvent) error {
	_, err := db.Model(record).Insert()
	if err != nil {
		return err
	}
	return nil
}

func SetupTables(db *pg.DB) error {
	var p client.PriceEvent

	err := db.Model(&p).CreateTable(&orm.CreateTableOptions{
		Temp:        false,
		IfNotExists: true,
	})
	if err != nil {
		return err
	}

	return nil
}
