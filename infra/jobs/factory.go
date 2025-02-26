package jobs

import "github.com/wolfsblu/go-chef/domain"

func NewScheduler(service *domain.RecipeService) *Scheduler {
	s := &Scheduler{
		service: service,
	}
	s.Start()
	return s
}
