package jobs

import (
	"context"
	"github.com/wolfsblu/go-chef/domain"
)

var quit chan struct{}

func StartScheduler(service *domain.RecipeService) {
	quit = make(chan struct{})
	initializeTickers()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-getC(cleanupPasswordResets):
				_ = service.RemoveObsoletePasswordResets(ctx)
			case <-getC(cleanupRegistrations):
				_ = service.RemoveObsoleteRegistrations(ctx)
			case <-quit:
				cancel()
				stopTickers()
				return
			}
		}
	}()
}

func QuitScheduler() {
	close(quit)
}
