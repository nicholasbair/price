package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"price/config"
)

// Transformed / flattened price
type Price struct {
	Type       string
	Time       float32
	Bid        float32
	Ask        float32
	Tradeable  bool
	Instrument string
}

type PriceEvent struct {
	Type       string      `json:"type"`
	Time       string      `json:"time"`
	Bids       []PriceItem `json:"bids"`
	Asks       []PriceItem `json:"asks"`
	Tradeable  bool        `json:"tradeable"`
	Instrument string      `json:"instrument"`
}

type PriceItem struct {
	Price string `json:"price"`
}

// StartPriceStream streams prices from Oanda for the provided instruments
func StartPriceStream(c chan PriceEvent, accountId string, instruments string) {
	req, reqErr := http.NewRequest("GET", getHost()+accountId+"/pricing/stream?instruments="+instruments, nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TOKEN"))
	req.Header.Set("Accept-Datetime-Format", "UNIX")

	if reqErr != nil {
		panic("Unable to configure request for price stream")
	}

	resp, respErr := http.DefaultClient.Do(req)

	if respErr != nil || resp.StatusCode != 200 {
		println("Restarting")
		StartPriceStream(c, accountId, instruments)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		p := new(PriceEvent)
		line, _ := reader.ReadBytes('\n')

		if err := json.Unmarshal([]byte(line), &p); err != nil {
			fmt.Println("Price: Can't unmarshal:", err)
			fmt.Println("Line:", line)
			if closed := req.Close; !closed {
				fmt.Println("Price: unable to close request")
			}
			StartPriceStream(c, accountId, instruments)
		}

		c <- *p
	}
}

// -- Private --
func getHost() string {
	if config.IsBacktest() {
		return "http://localhost:3000/accounts/"
	} else {
		return "https://stream-fxtrade.oanda.com/v3/accounts/"
	}
}
