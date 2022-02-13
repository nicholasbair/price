package handler

import (
	"github.com/go-pg/pg/v10"
	"log"
	"price/client"
	"price/store"
	"strconv"
	"time"
)

const maxRestartWait = 60

func PriceStreamHandler(db *pg.DB, accountId string, instruments string, token string, restartCount int) {
	p := make(chan client.PriceEvent)
	go client.StartPriceStream(p, accountId, instruments, token)

	for priceEvent := range p {

		// Reset the counter if receiving events in the channel
		if restartCount != 0 {
			restartCount = 0
		}

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
						log.Println("INSERT ERROR", err)
					}
				}
			} else {
				log.Println("ERROR", timeErr)
			}

		}
	}
	newCount := getNewCount(restartCount, maxRestartWait)
	log.Println("PRICE: Restarting in", newCount)
	time.Sleep(time.Duration(newCount) * time.Second)
	PriceStreamHandler(db, accountId, instruments, token, newCount)
}

// -- Private --
func strToFloat(str string) float64 {
	flt, err := strconv.ParseFloat(str, 32)

	if err != nil {
		panic("Unable to convert float " + str)
	}
	return flt
}

func getNewCount(current int, max int) int {
	if current >= max {
		return max
	}
	return current + 1
}
