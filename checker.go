package health

import (
	"sync"
	"time"
)

// Checker is a interface used to provide an indication of application health.
type Checker interface {
	Check() Health
}

// CheckerFunc is an adapter to allow the use of
// ordinary go functions as Checkers.
type CheckerFunc func() Health

func (f CheckerFunc) Check() Health {
	return f()
}

type checkerItem struct {
	name    string
	checker Checker
}

// CompositeChecker aggregate a list of Checkers
type CompositeChecker struct {
	checkers           []checkerItem
	info               map[string]interface{}
	useCachedLastState bool
	lastState          *lastState
	done               chan struct{}
}

type lastState struct {
	health Health
	sync.RWMutex
}

// NewCompositeChecker creates a new CompositeChecker
func NewCompositeChecker() *CompositeChecker {
	return &CompositeChecker{done: make(chan struct{})}
}

// Start will start a background Ticker to call Check() in an given interval
// Once Start() was called Check() will return the cached lastState result
func (c *CompositeChecker) Start(interval time.Duration) {
	if c.useCachedLastState == true {
		return
	}
	ticker := time.NewTicker(interval)
	c.useCachedLastState = true
	lastState := lastState{health: NewHealth()}
	c.lastState = &lastState

	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-ticker.C:
				currentHealth := c.check()
				lastState.Lock()
				lastState.health = currentHealth
				lastState.Unlock()
			}
		}
	}()
}

// Stop the background Ticker
func (c *CompositeChecker) Stop() {
	if c.useCachedLastState == false {
		return
	}
	c.useCachedLastState = false
	c.done <- struct{}{}
}

// AddInfo adds a info value to the Info map
func (c *CompositeChecker) AddInfo(key string, value interface{}) *CompositeChecker {
	if c.info == nil {
		c.info = make(map[string]interface{})
	}

	c.info[key] = value

	return c
}

// AddChecker add a Checker to the aggregator
func (c *CompositeChecker) AddChecker(name string, checker Checker) {
	c.checkers = append(c.checkers, checkerItem{name: name, checker: checker})
}

// Check returns the combination of all checkers added
// if some check is not up, the combined is marked as down
func (c CompositeChecker) Check() Health {
	if c.useCachedLastState {
		c.lastState.RLock()
		defer c.lastState.RUnlock()
		return c.lastState.health
	}
	return c.check()
}

func (c CompositeChecker) check() Health {
	health := NewHealth()
	health.Up()

	healths := make(map[string]interface{})

	type state struct {
		h    Health
		name string
	}
	ch := make(chan state, len(c.checkers))
	var wg sync.WaitGroup
	for _, item := range c.checkers {
		wg.Add(1)
		item := item
		go func() {
			ch <- state{h: item.checker.Check(), name: item.name}
			wg.Done()
		}()
	}
	wg.Wait()
	close(ch)

	for s := range ch {
		if !s.h.IsUp() && !health.IsDown() {
			health.Down()
		}
		healths[s.name] = s.h
	}

	health.info = healths

	// Extra Info
	for key, value := range c.info {
		health.AddInfo(key, value)
	}
	return health
}
