package tickers

import "time"

type Type string

var tickerMap map[Type]*time.Ticker

var (
	CleanupPasswordResets = Type("CleanupPasswordResets")
	CleanupRegistrations  = Type("CleanupRegistrations")
)

func Initialize() {
	tickerMap = map[Type]*time.Ticker{
		CleanupPasswordResets: time.NewTicker(24 * time.Hour),
		CleanupRegistrations:  time.NewTicker(24 * time.Hour),
	}
}

func GetC(t Type) <-chan time.Time {
	if ticker, ok := tickerMap[t]; ok {
		return ticker.C
	}
	return make(chan time.Time)
}

func Stop() {
	for _, t := range tickerMap {
		t.Stop()
	}
}
