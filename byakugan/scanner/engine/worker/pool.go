package worker

import (
	"context"
	"fmt"

	"github.com/haniz/byakugan/scanner/engine/detector"
	pb "github.com/haniz/byakugan/scanner/proto"
)

// NewPool creates new worker pool
func NewPool(config *pb.TaskConfig, det detector.Detector, logger Logger) (*Pool, error) {
	if config == nil {
		return nil, fmt.Errorf("nil config")
	}

	ctx, cancel := context.WithCancel(context.Background())

	pool := &Pool{
		taskChan:   make(chan *pb.ScanTask, int(config.MaxRetries)),
		resultChan: make(chan *pb.ScanResult, int(config.MaxRetries)),
		detector:   det,
		logger:     logger,
		metrics:    &Metrics{},
		ctx:        ctx,
		cancel:     cancel,
	}

	// Initialize workers
	for i := 0; i < int(config.MaxRetries); i++ {
		worker := NewWorker(i, pool.taskChan, pool.resultChan, det, logger)
		pool.workers = append(pool.workers, worker)
	}

	return pool, nil
}

// Start starts all workers in the pool
func (p *Pool) Start() error {
	for _, w := range p.workers {
		p.wg.Add(1)
		go func(worker *Worker) {
			defer p.wg.Done()
			worker.Start(p.ctx)
		}(w)
	}
	return nil
}

// Submit submits a task to the worker pool
func (p *Pool) Submit(task *pb.ScanTask) error {
	select {
	case p.taskChan <- task:
		return nil
	case <-p.ctx.Done():
		return fmt.Errorf("pool is shutting down")
	}
}

// Stop gracefully stops the worker pool
func (p *Pool) Stop() {
	p.cancel()
	close(p.taskChan)
	p.wg.Wait()
}
