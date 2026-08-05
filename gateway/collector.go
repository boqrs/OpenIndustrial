package gateway

import (
	"context"
	"sync"
	"time"

	"github.com/OpenGongChang/OpenIndustrial/runtime/driver"
)

// Collector receives events from drivers
// and distributes them inside gateway.
//
// Responsibilities:
// 1. Receive driver events
// 2. Update realtime cache
// 3. Broadcast events to subscribers
//
// Collector does NOT:
// - communicate with devices
// - persist data
// - upload cloud data
type Collector struct {

	// event input channel from drivers
	eventCh chan driver.Event

	// realtime state cache
	cache *Cache

	// event subscribers
	outputs []chan driver.Event

	// lifecycle
	stopCh chan struct{}

	wg sync.WaitGroup

	mu sync.RWMutex

	stopped bool
}


// NewCollector creates collector instance.
func NewCollector(
	cache *Cache,
) *Collector {

	return &Collector{

		eventCh: make(
			chan driver.Event,
			1024,
		),

		cache: cache,

		outputs: make(
			[]chan driver.Event,
			0,
		),

		stopCh: make(
			chan struct{},
		),
	}
}



// Submit submits event from driver.
//
// Driver should never access cache directly.
// Driver only produces events.
func (c *Collector) Submit(
	event driver.Event,
) {


	if event.Timestamp.IsZero() {

		event.Timestamp = time.Now()

	}


	select {

	case c.eventCh <- event:


	default:
		// overload protection
		// drop event
	}

}



// Subscribe creates event subscriber.
//
// Used by:
// - publisher
// - alarm engine
// - storage engine
// - rule engine
func (c *Collector) Subscribe() <-chan driver.Event {


	ch := make(
		chan driver.Event,
		256,
	)


	c.mu.Lock()
	defer c.mu.Unlock()


	if c.stopped {

		close(ch)
		return ch
	}


	c.outputs = append(
		c.outputs,
		ch,
	)


	return ch
}



// Start starts collector loop.
func (c *Collector) Start(
	ctx context.Context,
) {


	c.wg.Add(1)


	go c.run(ctx)

}



// run collector main loop.
func (c *Collector) run(
	ctx context.Context,
) {


	defer c.wg.Done()



	for {


		select {


		case event := <-c.eventCh:


			// 1. update cache

			if c.cache != nil {

				c.cache.Update(event)

			}


			// 2. broadcast

			c.broadcast(event)



		case <-ctx.Done():

			return



		case <-c.stopCh:

			return
		}

	}

}




// broadcast sends event to all subscribers.
func (c *Collector) broadcast(
	event driver.Event,
) {


	c.mu.RLock()
	defer c.mu.RUnlock()



	for _, ch := range c.outputs {


		select {


		case ch <- event:


		default:

			// subscriber slow
			// drop event

		}

	}

}




// Stop stops collector.
func (c *Collector) Stop() {


	c.mu.Lock()


	if c.stopped {

		c.mu.Unlock()
		return
	}


	c.stopped = true


	close(c.stopCh)


	for _, ch := range c.outputs {

		close(ch)

	}


	c.outputs = nil


	c.mu.Unlock()



	c.wg.Wait()

}




// Events returns input channel.
//
// Usually used internally by Gateway.
func (c *Collector) Events() chan<- driver.Event {

	return c.eventCh

}