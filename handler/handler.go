package handler

import (
	"fmt"
	"github.com/go-pg/pg"
	"price/client"
	"price/store"
	"strconv"
)

func PriceStreamHandler(db *pg.DB, accountId string, instrument string) {
	p := make(chan client.PriceEvent)
	go client.StartPriceStream(p, accountId, instrument)

	for price := range p {
		if price.Tradeable && price.Type != "HEARTBEAT" {
			err := store.Insert(db, &price)
			if err != nil {
				fmt.Println("INSERT ERROR", err)
			}
		}
	}
}

// -- Private --
func strToFloat(str string) float32 {
	flt64, err := strconv.ParseFloat(str, 32)
	flt := float32(flt64)

	if err != nil {
		panic("Unable to convert float " + str)
	}
	return flt
}
