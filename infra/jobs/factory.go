package jobs

import "github.com/wolfsblu/go-chef/domain"

func NewScheduler(service *domain.RecipeService) *Scheduler {
	return &Scheduler{
		service: service,
	}
}
