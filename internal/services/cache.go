package services

import (
	"sync"
	"time"
)

type CacheResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
	Ttl        time.Duration
	createdAt  time.Time
}

type CacheService struct {
	lock     *sync.RWMutex
	requests map[string]CacheResponse
}

func NewCacheService() *CacheService {
	return &CacheService{
		requests: make(map[string]CacheResponse),
	}
}

func (service *CacheService) Get(key string) (CacheResponse, bool) {
	service.lock.RLock()
	defer service.lock.RUnlock()

	response, ok := service.requests[key]

	if !ok {
		return CacheResponse{}, false
	}

	if response.createdAt.Add(response.Ttl).Before(time.Now()) {
		return CacheResponse{}, false
	}

	return response, true
}

func (service *CacheService) Set(key string, response CacheResponse) {
	service.lock.Lock()
	defer service.lock.Unlock()

	response.createdAt = time.Now()

	service.requests[key] = response
}
