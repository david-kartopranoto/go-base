package util

import (
	"errors"
	"net/http"

	"golang.org/x/time/rate"
)

//LimiterService implements LimiterService interface
type LimiterService struct {
	limiter *rate.Limiter
}

//NewLimiterService create a new go rate limiter service
func NewLimiterService(config Config) (*LimiterService, error) {
	limiter := rate.NewLimiter(rate.Limit(config.Limiter.MaxEventsPerSec),
		config.Limiter.MaxBurstSize)

	s := &LimiterService{
		limiter: limiter,
	}

	return s, nil
}

func (s *LimiterService) Allow() bool {
	return s.limiter.Allow()
}

func (s *LimiterService) GetDefaultError() (error, int) {
	return errors.New("Limit exceeded"), http.StatusTooManyRequests
}
