package store

import (
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
		db := pg.Connect(&pg.Options{
			User:     config.GetEnv("DB_USER"),
			Password: config.GetEnv("DB_PASSWORD"),
			Database: config.GetEnv("DB_NAME"),
			Addr:     config.GetEnv("DB_ADDR"),
		})
		return db
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
