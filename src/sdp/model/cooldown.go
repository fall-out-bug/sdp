package model

import (
	"context"
	"sync"
	"time"
)

// CooldownConfig configures cooldown behavior
type CooldownConfig struct {
	DefaultCooldown time.Duration // Default cooldown between calls
	GlobalRate      time.Duration // Minimum time between any calls
	MaxBackoff      time.Duration // Maximum backoff on rate limit
}

// DefaultCooldownConfig provides sensible defaults
var DefaultCooldownConfig = CooldownConfig{
	DefaultCooldown: 100 * time.Millisecond,
	GlobalRate:      50 * time.Millisecond,
	MaxBackoff:      30 * time.Second,
}

// CooldownManager manages rate limiting between API calls
type CooldownManager struct {
	config   CooldownConfig
	lastCall map[string]time.Time
	globalMu sync.Mutex
	last     time.Time
	cooldown map[string]time.Duration
	mu       sync.RWMutex
}

// NewCooldownManager creates a new cooldown manager
func NewCooldownManager(config CooldownConfig) *CooldownManager {
	return &CooldownManager{
		config:   config,
		lastCall: make(map[string]time.Time),
		cooldown: make(map[string]time.Duration),
	}
}

// WaitForCooldown waits until the cooldown period has passed
func (c *CooldownManager) WaitForCooldown(ctx context.Context, model string) error {
	// First check global rate limit
	if err := c.waitForGlobal(ctx); err != nil {
		return err
	}

	// Then check model-specific cooldown
	c.mu.RLock()
	cooldown := c.cooldown[model]
	if cooldown == 0 {
		cooldown = c.config.DefaultCooldown
	}
	lastCall := c.lastCall[model]
	c.mu.RUnlock()

	if waitTime := time.Since(lastCall); waitTime < cooldown {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(cooldown - waitTime):
		}
	}

	return nil
}

// waitForGlobal handles global rate limiting
func (c *CooldownManager) waitForGlobal(ctx context.Context) error {
	c.globalMu.Lock()
	defer c.globalMu.Unlock()

	if waitTime := time.Since(c.last); waitTime < c.config.GlobalRate {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.config.GlobalRate - waitTime):
		}
	}

	c.last = time.Now()
	return nil
}

// RecordCall records that a call was made to a model
func (c *CooldownManager) RecordCall(model string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastCall[model] = time.Now()
}

// SetCooldown sets a custom cooldown for a model
func (c *CooldownManager) SetCooldown(model string, cooldown time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cooldown[model] = cooldown
}

// Backoff increases cooldown after rate limit error
func (c *CooldownManager) Backoff(model string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	current := c.cooldown[model]
	if current == 0 {
		current = c.config.DefaultCooldown
	}

	// Double the cooldown, capped at max
	newCooldown := current * 2
	if newCooldown > c.config.MaxBackoff {
		newCooldown = c.config.MaxBackoff
	}

	c.cooldown[model] = newCooldown
}

// Reset resets cooldown for a model
func (c *CooldownManager) Reset(model string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cooldown, model)
}

// ResetAll resets all cooldowns
func (c *CooldownManager) ResetAll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cooldown = make(map[string]time.Duration)
	c.lastCall = make(map[string]time.Time)
}

// GetCooldown returns the current cooldown for a model
func (c *CooldownManager) GetCooldown(model string) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if cd, ok := c.cooldown[model]; ok {
		return cd
	}
	return c.config.DefaultCooldown
}

// TimeUntilNext returns time until next call is allowed
func (c *CooldownManager) TimeUntilNext(model string) time.Duration {
	c.mu.RLock()
	cooldown := c.cooldown[model]
	if cooldown == 0 {
		cooldown = c.config.DefaultCooldown
	}
	lastCall := c.lastCall[model]
	c.mu.RUnlock()

	remaining := cooldown - time.Since(lastCall)
	if remaining < 0 {
		return 0
	}
	return remaining
}
