package allocator

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Allocator interface {
	Allocate(ctx context.Context, matchID string) (string, error)
}

type MockAllocator struct {
	Region string
}

func NewMockAllocator(region string) *MockAllocator {
	return &MockAllocator{Region: region}
}

func (m *MockAllocator) Allocate(ctx context.Context, matchID string) (string, error) {
	select {
	case <-time.After(200 * time.Millisecond):
	case <-ctx.Done():
		return "", ctx.Err()
	}

	fakeIP := fmt.Sprintf("10.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
	port := rand.Intn(1000) + 9000

	address := fmt.Sprintf("%s:%d", fakeIP, port)
	
	fmt.Printf("☁️  [Infrastructure] Allocated Server %s for Match %s in %s\n", 
		address, matchID, m.Region)

	return address, nil
}
