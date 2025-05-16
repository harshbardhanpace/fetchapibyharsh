package throttling

import (
	"sync"
	"time"
)

type APIThrottler struct {
	RequestsPerSecond int
	LastRequestTime   time.Time
	Mu                sync.Mutex
	ClientName        string
	VendorName        string
	Count             int
}

func (t *APIThrottler) Throttle() {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	if t.Count == 0 {
		t.LastRequestTime = time.Now()
	}
	t.Count++

	if t.Count >= t.RequestsPerSecond {
		currTime := time.Now()
		sleepDurationNew := time.Second - currTime.Sub(t.LastRequestTime)
		time.Sleep(sleepDurationNew)
		t.Count = 0
	}

}

func (t *APIThrottler) ThrottleDeprecate() {
	t.Mu.Lock()
	defer t.Mu.Unlock()

	currentTime := time.Now()
	timeSinceLastRequest := currentTime.Sub(t.LastRequestTime)
	desiredInterval := time.Second / time.Duration(t.RequestsPerSecond)

	if timeSinceLastRequest < desiredInterval {
		sleepDuration := desiredInterval - timeSinceLastRequest
		time.Sleep(sleepDuration)
	}

	t.LastRequestTime = time.Now()
}

// APIClient represents an API client.
type APIClient struct {
	Throttler *APIThrottler
}

// NewAPIClient creates a new API client instance.
func NewAPIClient(throttler *APIThrottler) *APIClient {
	return &APIClient{
		Throttler: throttler,
	}
}
