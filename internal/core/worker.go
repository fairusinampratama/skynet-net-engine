package core

import (
	"time"
	"sync"
	"skynet-net-engine-api/internal/mikrotik"
	"skynet-net-engine-api/internal/models"
	"skynet-net-engine-api/internal/database"
	"skynet-net-engine-api/pkg/logger"
	"go.uber.org/zap"
)

type Worker struct {
	Router   models.Router
	CmdChan  chan Command
	Client   *mikrotik.Client
	IsOnline bool
	
	// Synchronization
	once sync.Once
	wg   *sync.WaitGroup

	// Cache
	ActiveUsers    []models.ActiveUser
	SystemResource *models.SystemResource
	Lock           sync.RWMutex
}

func NewWorker(r models.Router, wg *sync.WaitGroup) *Worker {
	return &Worker{
		Router:  r,
		CmdChan: make(chan Command, 10), // Buffered channel
		wg:      wg,
	}
}

// Start begins the persistent loop
func (w *Worker) Start() {
	// Ensure we always mark as done eventually, even if we crash or never connect
	// But ideally we mark done inside the loop upon success/failure decision
	// For now, simpler: Signal ready on first connect OR first timeout failure
	
	// Helper to signal readiness once
	signalReady := func() {
		w.once.Do(func() {
			if w.wg != nil {
				w.wg.Done()
			}
		})
	}

	go w.metricsLoop() // Start background metrics/keepalive

	for {
		// 1. Try to Connect
		logger.Info("Dialing router...", zap.String("host", w.Router.Host), zap.Int("port", w.Router.Port), zap.String("user", w.Router.Username))
		client, err := mikrotik.NewClient(w.Router)
		
		if err != nil {
			logger.Error("Connection failed, retrying in 5s...", zap.String("host", w.Router.Host), zap.Int("port", w.Router.Port), zap.Error(err))
			w.IsOnline = false
			
			// If we fail the first connect, we consider this worker "warmed up" (but failed)
			// so we don't block the entire server forever.
			signalReady()
			
			time.Sleep(5 * time.Second)
			continue // Retry loop
		}

		// 2. Connected!
		w.Client = client
		w.IsOnline = true
		logger.Info("Router Connected!", zap.String("host", w.Router.Host))
		SendWebhook("router.up", w.Router.ID, w.Router.Host, nil)

		// 3. WARMUP: Fetch initial data IMMEDIATELY
		logger.Info("Warming up cache...", zap.String("host", w.Router.Host))
		w.refreshMetrics() // Force immediate fetch
		
		// Trigger initial Sync of Secrets (Async)
		w.CmdChan <- Command{Type: CmdSync}
		
		signalReady()      // Signal we are ready to serve

		// 4. Command Loop (Blocks until connection dies)
		w.handleCommands()

		// 5. Cleanup after disconnect
		logger.Warn("Router Disconnected. Cleaning up...", zap.String("host", w.Router.Host))
		if w.IsOnline {
			SendWebhook("router.down", w.Router.ID, w.Router.Host, "Connection lost")
		}
		w.IsOnline = false
		if w.Client != nil {
			w.Client.Close()
		}
		
		// 6. Backoff before reconnecting
		time.Sleep(3 * time.Second)
	}
}

func (w *Worker) handleCommands() {
	for cmd := range w.CmdChan {
		// Process command here
		// If TCP fails, we break the loop and let Start() reconnect
		logger.Info("Received command", zap.String("type", string(cmd.Type)))
		
		var err error
		switch cmd.Type {
		case CmdSync:
			secrets, errSync := w.Client.GetAllSecrets()
			if errSync != nil {
				err = errSync
				logger.Error("Failed to fetch secrets for sync", zap.String("router", w.Router.Name), zap.Error(err))
			} else {
				go func(secrets []models.PPPoESecret, routerID int) {
					count := 0
					for _, s := range secrets {
						if dbErr := database.UpsertUser(s.Name, routerID, s.Profile, s.RemoteAddress, !s.Disabled); dbErr == nil {
							count++
						}
					}
					logger.Info("Synced Secrets to DB", zap.String("router", w.Router.Name), zap.Int("synced", count), zap.Int("total", len(secrets)))
				}(secrets, w.Router.ID)
				
				cmd.Result <- "Synced"
			}
		
		case CmdCreateSecret:
			payload := cmd.Payload.(map[string]string)
			err = w.Client.AddSecret(
				payload["user"], 
				payload["password"], 
				payload["profile"], 
				payload["local_ip"], 
				payload["remote_ip"], 
				payload["comment"],
			)
		
		case CmdUpdateSecret:
			payload := cmd.Payload.(map[string]string)
			err = w.Client.SetSecretProfile(
				payload["user"],
				payload["profile"],
			)

		case CmdIsolate:
			payload := cmd.Payload.(map[string]string)
			if payload["action"] == "add" {
				err = w.Client.AddAddressList(payload["ip"], payload["list"], payload["comment"])
			} else {
				err = w.Client.RemoveAddressList(payload["ip"], payload["list"])
			}

		case CmdGetTraffic:
			target := cmd.Payload.(string)
			stats, errT := w.Client.GetQueueTraffic(target)
			if errT != nil {
				err = errT
			} else {
				cmd.Result <- stats
				return // Early return since we sent result
			}

		case CmdBackup:
			name := cmd.Payload.(string)
			err = w.Client.RunBackup(name)

		case CmdRefreshMetrics:
			w.refreshMetrics()
		}

		if cmd.Error != nil {
			cmd.Error <- err
		}
		if cmd.Result != nil && err == nil {
			cmd.Result <- "Success"
		}
	}
}

func (w *Worker) metricsLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !w.IsOnline {
			continue
		}
		// Thread Safety: Send command instead of direct call
		w.CmdChan <- Command{Type: CmdRefreshMetrics}
	}
}

func (w *Worker) refreshMetrics() {
	if w.Client == nil {
		return
	}

	users, err := w.Client.GetActiveUsers()
	if err != nil {
		logger.Error("Failed to fetch active users", zap.String("host", w.Router.Host), zap.Error(err))
	} 
	
	res, errRes := w.Client.GetSystemResource()
	if errRes != nil {
		// logging error optional
	}

	// Update Cache
	w.Lock.Lock()
	if err == nil {
		w.ActiveUsers = users
		logger.Info("Worker Cache Updated", zap.String("router", w.Router.Name), zap.Int("active_users", len(w.ActiveUsers)))
	}
	if errRes == nil {
		w.SystemResource = res
	}
	w.Lock.Unlock()
	
	// logger.Info("Metrics refreshed", zap.String("host", w.Router.Host), zap.Int("users", len(users)))
}
