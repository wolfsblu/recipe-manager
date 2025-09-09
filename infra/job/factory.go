package job

import "github.com/wolfsblu/go-chef/domain"

func NewScheduler(service *domain.UserService) *Scheduler {
	s := &Scheduler{
		service: service,
	}
	s.Start()
	return s
}
