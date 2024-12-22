package jobs

import (
	"context"
	"github.com/wolfsblu/go-chef/domain"
	"github.com/wolfsblu/go-chef/infra/jobs/tickers"
)

var quit chan struct{}

func StartScheduler(service *domain.RecipeService) {
	quit = make(chan struct{})
	tickers.Initialize()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-tickers.GetC(tickers.CleanupPasswordResets):
				_ = service.RemoveObsoletePasswordResets(ctx)
			case <-tickers.GetC(tickers.CleanupRegistrations):
				_ = service.RemoveObsoleteRegistrations(ctx)
			case <-quit:
				cancel()
				tickers.Stop()
				return
			}
		}
	}()
}

func QuitScheduler() {
	close(quit)
}
