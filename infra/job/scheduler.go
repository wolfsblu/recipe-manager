package job

import (
	"context"
	"github.com/wolfsblu/go-chef/domain"
	"time"
)

type Scheduler struct {
	quit    chan struct{}
	service *domain.RecipeService
}

func (s *Scheduler) Start() {
	s.quit = make(chan struct{})
	initializeTickers()

	oneWeek := 7 * 24 * time.Hour
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-getC(cleanupPasswordResets):
				go func() {
					_ = s.service.DeletePasswordResetsOlderThan(ctx, oneWeek)
				}()
			case <-getC(cleanupRegistrations):
				go func() {
					_ = s.service.DeleteRegistrationsOlderThan(ctx, oneWeek)
				}()
			case <-s.quit:
				cancel()
				stopTickers()
				return
			}
		}
	}()
}

func (s *Scheduler) Quit() {
	close(s.quit)
}
