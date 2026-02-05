package worker

import (
	"context"
	"fmt"
	"hash/fnv"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

type Command struct {
	AccountID uuid.UUID
	Exec      func(ctx context.Context) error
	Err       chan error
}

type Pool struct {
	workers int
	queues  []chan Command
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewPool(workers, queueSize int) *Pool {
	if workers <= 0 {
		workers = 10
	}
	if queueSize <= 0 {
		queueSize = 100
	}

	ctx, cancel := context.WithCancel(context.Background())
	p := &Pool{
		workers: workers,
		queues:  make([]chan Command, workers),
		ctx:     ctx,
		cancel:  cancel,
	}

	for i := range workers {
		p.queues[i] = make(chan Command, queueSize)
		p.wg.Add(1)
		go p.processQueue(i, p.queues[i])
	}

	slog.Info("worker pool started", "workers", workers, "queue_size", queueSize)
	return p
}

func (p *Pool) Submit(cmd Command) error {
	select {
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool shutting down")
	default:
	}

	idx := p.shardIndex(cmd.AccountID)

	select {
	case p.queues[idx] <- cmd:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("worker pool shutting down")
	}
}

func (p *Pool) shardIndex(accountID uuid.UUID) int {
	h := fnv.New32a()
	h.Write(accountID[:])
	return int(h.Sum32()) % p.workers
}

func (p *Pool) processQueue(workerID int, ch chan Command) {
	defer p.wg.Done()

	for {
		select {
		case cmd, ok := <-ch:
			if !ok {
				return
			}
			err := cmd.Exec(p.ctx)
			if err != nil {
				slog.Error("command execution failed",
					"worker", workerID,
					"account_id", cmd.AccountID,
					"error", err,
				)
			}
			cmd.Err <- err
		case <-p.ctx.Done():
			p.drainQueue(ch)
			return
		}
	}
}

func (p *Pool) drainQueue(ch chan Command) {
	for {
		select {
		case cmd, ok := <-ch:
			if !ok {
				return
			}
			cmd.Err <- fmt.Errorf("worker pool shutting down")
		default:
			return
		}
	}
}

func (p *Pool) Shutdown() {
	p.cancel()

	for _, ch := range p.queues {
		close(ch)
	}

	p.wg.Wait()
	slog.Info("worker pool shut down")
}
