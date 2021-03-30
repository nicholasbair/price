package client

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"price/config"
	"runtime"
	"time"
)

// Transformed / flattened price
type Price struct {
	Id         int64
	Type       string
	Time       time.Time
	Bid        float64
	Ask        float64
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
func StartPriceStream(c chan PriceEvent, accountId string, instruments string, token string) {
	req, reqErr := http.NewRequest("GET", getHost()+accountId+"/pricing/stream", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	if reqErr != nil {
		panic("Unable to configure request for price stream")
	}

	q := req.URL.Query()
	q.Add("instruments", instruments)
	req.URL.RawQuery = q.Encode()

	resp, respErr := http.DefaultClient.Do(req)

	if respErr != nil || resp.StatusCode != 200 {
		fmt.Println("Restarting due to", respErr, resp.StatusCode)
		runtime.GC()
		StartPriceStream(c, accountId, instruments, token)
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
			runtime.GC()
			StartPriceStream(c, accountId, instruments, token)
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
