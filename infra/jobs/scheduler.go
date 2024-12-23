package jobs

import (
	"context"
	"github.com/wolfsblu/go-chef/domain"
	"time"
)

var quit chan struct{}

func StartScheduler(service *domain.RecipeService) {
	quit = make(chan struct{})
	initializeTickers()

	oneWeek := 7 * 24 * time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-getC(cleanupPasswordResets):
				go func() {
					_ = service.DeletePasswordResetsOlderThan(ctx, oneWeek)
				}()
			case <-getC(cleanupRegistrations):
				go func() {
					_ = service.DeleteRegistrationsOlderThan(ctx, oneWeek)
				}()
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
