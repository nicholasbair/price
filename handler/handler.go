package handler

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"price/client"
	"price/store"
	"strconv"
)

func PriceStreamHandler(db *pg.DB, accountId string, instrument string) {
	p := make(chan client.PriceEvent)
	go client.StartPriceStream(p, accountId, instrument)

	for priceEvent := range p {
		if priceEvent.Tradeable && priceEvent.Type != "HEARTBEAT" {

			price := client.Price{
				Type:       "PRICE",
				Time:       strToFloat(priceEvent.Time),
				Bid:        strToFloat(priceEvent.Bids[0].Price),
				Ask:        strToFloat(priceEvent.Asks[0].Price),
				Tradeable:  priceEvent.Tradeable,
				Instrument: instrument,
			}

			err := store.Insert(db, &price)
			if err != nil {
				fmt.Println("INSERT ERROR", err)
			} else {
				fmt.Println("INSERT: ", price.Instrument, price.Time)
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
