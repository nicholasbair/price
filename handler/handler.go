package handler

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"price/client"
	"price/store"
	"strconv"
	"time"
)

func PriceStreamHandler(db *pg.DB, accountId string, instruments string, token string) {
	p := make(chan client.PriceEvent)
	go client.StartPriceStream(p, accountId, instruments, token)

	for priceEvent := range p {
		if priceEvent.Tradeable && priceEvent.Type != "HEARTBEAT" {

			t, timeErr := time.Parse(time.RFC3339, priceEvent.Time)
			var price client.Price

			if timeErr == nil {
				price = client.Price{
					Type:       "PRICE",
					Time:       t,
					Bid:        strToFloat(priceEvent.Bids[0].Price),
					Ask:        strToFloat(priceEvent.Asks[0].Price),
					Tradeable:  priceEvent.Tradeable,
					Instrument: priceEvent.Instrument,
				}

				if price.Bid != 0.0 && price.Ask != 0.0 {
					err := store.Insert(db, &price)
					if err != nil {
						fmt.Println("INSERT ERROR", err)
					} else {
						fmt.Println("INSERT: ", price.Instrument)
					}
				}
			} else {
				fmt.Println("ERROR", timeErr)
			}

		}
	}
}

// -- Private --
func strToFloat(str string) float64 {
	flt, err := strconv.ParseFloat(str, 32)

	if err != nil {
		panic("Unable to convert float " + str)
	}
	return flt
}
