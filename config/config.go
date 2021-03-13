package config

import (
	"github.com/joho/godotenv"
	"os"
	"strings"
)

/* Contents of .env file
OANDA_TOKEN=1234
OANDA_ACCOUNTS=001-001-2420587-001,001-001-2420587-002
OANDA_INSTRUMENTS=EUR_USD,USD_SEK
BACKTEST=true
*/

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func GetEnv(key string) []string {
	val := os.Getenv(key)
	if val == "" {
		panic("ENV var missing " + key)
	}

	return strings.Split(val, ",")
}

func IsBacktest() bool {
	switch os.Getenv("BACKTEST") {
	case "true", "TRUE":
		return true
	case "false", "FALSE":
		return false
	default:
		panic("Missing BACKTEST env flag")
	}
}
