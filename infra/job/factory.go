package job

import "github.com/wolfsblu/recipe-manager/domain"

func NewScheduler(service *domain.UserService) *Scheduler {
	s := &Scheduler{
		service: service,
	}
	s.Start()
	return s
}
