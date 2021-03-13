package config

import (
	"github.com/joho/godotenv"
	"os"
)

/* Contents of .env file
TOKEN=1234
ACCOUNT=001-001-2420587-001
INSTRUMENTS=EUR_USD,USD_SEK
BACKTEST=true
DB_CONNECTION_STRING=...
*/

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}

func GetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic("ENV var missing " + key)
	}

	return val
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
