package core

import (
	"sync"
	"skynet-net-engine-api/internal/database"
	"skynet-net-engine-api/internal/models"
	"skynet-net-engine-api/pkg/logger"
	"go.uber.org/zap"
)

type Pool struct {
	Workers map[int]*Worker
	Lock    sync.RWMutex
}

var GlobalPool *Pool

func InitPool() {
	GlobalPool = &Pool{
		Workers: make(map[int]*Worker),
	}

	// 1. Fetch Routers
	routers, err := database.GetAllRouters()
	if err != nil {
		logger.Error("Failed to load routers for pool - Continuing with empty pool", zap.Error(err))
	}

	// 2. Spawn Workers
	for _, r := range routers {
		worker := NewWorker(r)
		
		GlobalPool.Lock.Lock()
		GlobalPool.Workers[r.ID] = worker
		GlobalPool.Lock.Unlock()

		// Start the engine in a persistent Goroutine
		go worker.Start()
	}

	logger.Info("Worker Pool Initialized", zap.Int("workers", len(routers)))
}

func (p *Pool) GetWorker(id int) *Worker {
	p.Lock.RLock()
	defer p.Lock.RUnlock()
	return p.Workers[id]
}

func (p *Pool) GetAllTargets() []models.ActiveUser {
	p.Lock.RLock()
	defer p.Lock.RUnlock()

	var total []models.ActiveUser
	for _, w := range p.Workers {
		w.Lock.RLock()
		// Copy users to avoid race conditions if underlying array changes
		if len(w.ActiveUsers) > 0 {
			batch := make([]models.ActiveUser, len(w.ActiveUsers))
			copy(batch, w.ActiveUsers)
			total = append(total, batch...)
		}
		w.Lock.RUnlock()
	}
	return total
}
