package store

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"price/client"
)

func ConnectToDatabase() *pg.DB {
	db := pg.Connect(&pg.Options{
		User: "postgres",
	})
	//defer db.Close()
	return db
}

func Insert(db *pg.DB, record interface{}) error {
	_, err = db.Model(record).Insert()
	if err != nil {
		return err
	}
	return nil
}

func SetupTables(db *pg.DB) error {
	models := []interface{}{
		(*client.Price(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
