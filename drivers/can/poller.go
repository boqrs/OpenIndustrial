package can

import (
	"context"
	"log"
	"sync"
	"time"
)

// TxConfig defines a message to be transmitted periodically.
type TxConfig struct {
	SignalID string        `json:"signalId"`
	Interval time.Duration `json:"interval"`
	Value    any           `json:"value"`
}

// Poller handles periodic transmission of CAN frames.
type Poller struct {
	adapter Adapter
	encoder Encoder
	signals map[string]Signal
	txTasks []*TxConfig
	wg      sync.WaitGroup
	cancel  context.CancelFunc
}

// NewPoller creates a new Poller.
func NewPoller(adapter Adapter, encoder Encoder, signals []Signal, txTasks []TxConfig) *Poller {
	signalMap := make(map[string]Signal)
	for _, s := range signals {
		signalMap[s.ID] = s
	}

	tasks := make([]*TxConfig, len(txTasks))
	for i := range txTasks {
		tasks[i] = &txTasks[i]
	}

	return &Poller{
		adapter: adapter,
		encoder: encoder,
		signals: signalMap,
		txTasks: tasks,
	}
}

// Start begins the periodic transmission of all configured tasks.
func (p *Poller) Start(ctx context.Context) {
	ctx, p.cancel = context.WithCancel(ctx)
	for _, task := range p.txTasks {
		p.wg.Add(1)
		go p.runTask(ctx, task)
	}
}

// Stop gracefully stops all transmission tasks.
func (p *Poller) Stop() {
	if p.cancel != nil {
		p.cancel()
	}
	p.wg.Wait()
}

// runTask is the goroutine for a single periodic transmission task.
func (p *Poller) runTask(ctx context.Context, task *TxConfig) {
	defer p.wg.Done()

	ticker := time.NewTicker(task.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.send(ctx, task)
		case <-ctx.Done():
			return
		}
	}
}

func (p *Poller) send(ctx context.Context, task *TxConfig) {
	signal, ok := p.signals[task.SignalID]
	if !ok {
		log.Printf("Poller: signal not found for ID %s", task.SignalID)
		return
	}

	frame, err := p.encoder.Encode(signal, task.Value)
	if err != nil {
		log.Printf("Poller: failed to encode signal %s: %v", task.SignalID, err)
		return
	}

	if err := p.adapter.Send(ctx, frame); err != nil {
		log.Printf("Poller: failed to send frame for signal %s: %v", task.SignalID, err)
	}
}