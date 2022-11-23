package batch

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"
)

/*
	mock service for this task
*/

// ErrBlocked reports if service is blocked.
var ErrBlocked = errors.New("blocked")

// Service defines external service that can process batches of items.
type Service interface {
	GetLimits() (n uint64, p time.Duration)
	Process(ctx context.Context, batch Batch) error
}

type someService struct {
	n           uint64
	p           time.Duration
	mu          sync.Mutex
	lastProcess time.Time
}

// Batch is a batch of items.
type Batch []Item

// Item is some abstract item.
type Item struct {
	Id int
}

func NewService(n uint64, p time.Duration) Service {
	return &someService{
		n:           n,
		p:           p,
		lastProcess: time.Now().Add(-2 * p),
	}
}

func (s *someService) Process(ctx context.Context, batch Batch) error {
	s.mu.Lock()
	if time.Since(s.lastProcess) <= s.p || uint64(len(batch)) > s.n {
		panic("Service locked!\n")
	}
	s.lastProcess = time.Now()
	s.mu.Unlock()
	log.Printf("Service: process batch with len %v\n", len(batch))
	for _, item := range batch {
		log.Printf("Service: process %v item", item.Id)
	}
	time.Sleep(s.p)
	return nil
}

func (s *someService) GetLimits() (n uint64, p time.Duration) {
	return s.n, s.p
}
