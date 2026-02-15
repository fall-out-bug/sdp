package runtime

import (
	"context"
	"sync"
	"time"
)

// DegradedMode tracks degraded mode state
type DegradedMode struct {
	mu          sync.RWMutex
	active      bool
	reason      string
	since       time.Time
	cachedItems map[string]interface{}
}

// GlobalDegradedMode is the global degraded mode instance
var GlobalDegradedMode = NewDegradedMode()

// NewDegradedMode creates a new degraded mode tracker
func NewDegradedMode() *DegradedMode {
	return &DegradedMode{
		cachedItems: make(map[string]interface{}),
	}
}

// IsActive returns whether degraded mode is active
func (d *DegradedMode) IsActive() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.active
}

// Reason returns the reason for degraded mode
func (d *DegradedMode) Reason() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.reason
}

// Since returns when degraded mode started
func (d *DegradedMode) Since() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.since
}

// Duration returns how long degraded mode has been active
func (d *DegradedMode) Duration() time.Duration {
	d.mu.RLock()
	defer d.mu.RUnlock()
	if !d.active {
		return 0
	}
	return time.Since(d.since)
}

// Enter activates degraded mode
func (d *DegradedMode) Enter(reason string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.active {
		d.active = true
		d.reason = reason
		d.since = time.Now()
	}
}

// Exit deactivates degraded mode
func (d *DegradedMode) Exit() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.active = false
	d.reason = ""
	d.since = time.Time{}
}

// SetCache stores a value in cache
func (d *DegradedMode) SetCache(key string, value interface{}) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cachedItems[key] = value
}

// GetCache retrieves a value from cache
func (d *DegradedMode) GetCache(key string) (interface{}, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	val, ok := d.cachedItems[key]
	return val, ok
}

// ClearCache clears all cached items
func (d *DegradedMode) ClearCache() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cachedItems = make(map[string]interface{})
}

// CachedKeys returns all cached keys
func (d *DegradedMode) CachedKeys() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	keys := make([]string, 0, len(d.cachedItems))
	for k := range d.cachedItems {
		keys = append(keys, k)
	}
	return keys
}

// ServiceChecker checks if a service is available
type ServiceChecker struct {
	endpoint    string
	timeout     time.Duration
	lastCheck   time.Time
	lastStatus  bool
	checkPeriod time.Duration
}

// NewServiceChecker creates a new service checker
func NewServiceChecker(endpoint string, timeout, checkPeriod time.Duration) *ServiceChecker {
	return &ServiceChecker{
		endpoint:    endpoint,
		timeout:     timeout,
		checkPeriod: checkPeriod,
	}
}

// IsAvailable checks if the service is available
func (sc *ServiceChecker) IsAvailable(ctx context.Context) bool {
	// Use cached status if within check period
	if time.Since(sc.lastCheck) < sc.checkPeriod {
		return sc.lastStatus
	}

	// Perform actual check
	ctx, cancel := context.WithTimeout(ctx, sc.timeout)
	defer cancel()

	// Simple check - in real implementation would make HTTP request
	// For now, just return true
	available := true

	sc.lastCheck = time.Now()
	sc.lastStatus = available
	return available
}

// CheckServiceHealth checks service health and updates degraded mode
func CheckServiceHealth(ctx context.Context, checker *ServiceChecker, dm *DegradedMode, serviceName string) error {
	if !checker.IsAvailable(ctx) {
		dm.Enter(serviceName + " unavailable")
		return context.DeadlineExceeded
	}

	// Exit degraded mode if service recovered
	if dm.IsActive() && dm.Reason() == serviceName+" unavailable" {
		dm.Exit()
	}

	return nil
}

// WithDegradedMode executes fn with degraded mode handling
func WithDegradedMode(ctx context.Context, dm *DegradedMode, fn func() error, fallback func() error) error {
	if dm.IsActive() {
		if fallback != nil {
			return fallback()
		}
	}

	err := fn()
	if err != nil {
		dm.Enter("operation failed: " + err.Error())
		if fallback != nil {
			return fallback()
		}
		return err
	}

	return nil
}
