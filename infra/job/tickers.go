package job

import "time"

type tickerType string

var tickerMap map[tickerType]*time.Ticker

var (
	cleanupPasswordResets = tickerType("cleanupPasswordResets")
	cleanupRegistrations  = tickerType("cleanupRegistrations")
)

func initializeTickers() {
	tickerMap = map[tickerType]*time.Ticker{
		cleanupPasswordResets: time.NewTicker(24 * time.Hour),
		cleanupRegistrations:  time.NewTicker(24 * time.Hour),
	}
}

func getC(t tickerType) <-chan time.Time {
	if ticker, ok := tickerMap[t]; ok {
		return ticker.C
	}
	return make(chan time.Time)
}

func stopTickers() {
	for _, t := range tickerMap {
		t.Stop()
	}
}
