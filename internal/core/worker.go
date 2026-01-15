package core

import (
	"time"
	"sync"
	"skynet-net-engine-api/internal/mikrotik"
	"skynet-net-engine-api/internal/models"
	"skynet-net-engine-api/pkg/logger"
	"go.uber.org/zap"
)

type Worker struct {
	Router   models.Router
	CmdChan  chan Command
	Client   *mikrotik.Client
	IsOnline bool
	
	// Cache
	ActiveUsers    []models.ActiveUser
	SystemResource *models.SystemResource
	Lock           sync.RWMutex
}

func NewWorker(r models.Router) *Worker {
	return &Worker{
		Router:  r,
		CmdChan: make(chan Command, 10), // Buffered channel
	}
}

// Start begins the persistent loop
func (w *Worker) Start() {
	go w.metricsLoop() // Start background metrics/keepalive

	for {
		// 1. Try to Connect
		logger.Info("Dialing router...", zap.String("host", w.Router.Host), zap.String("user", w.Router.Username))
		client, err := mikrotik.NewClient(w.Router)
		
		if err != nil {
			logger.Error("Connection failed, retrying in 5s...", zap.String("host", w.Router.Host), zap.Error(err))
			w.IsOnline = false
			time.Sleep(5 * time.Second)
			continue // Retry loop
		}

		// 2. Connected!
		w.Client = client
		w.IsOnline = true
		logger.Info("Router Connected!", zap.String("host", w.Router.Host))
		SendWebhook("router.up", w.Router.ID, w.Router.Host, nil)

		// 3. Command Loop (Blocks until connection dies)
		w.handleCommands()

		// 4. Cleanup after disconnect
		logger.Warn("Router Disconnected. Cleaning up...", zap.String("host", w.Router.Host))
		if w.IsOnline {
			SendWebhook("router.down", w.Router.ID, w.Router.Host, "Connection lost")
		}
		w.IsOnline = false
		if w.Client != nil {
			w.Client.Close()
		}
		
		// 5. Backoff before reconnecting
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
			// TODO: Implement Sync
			cmd.Result <- "Synced"
		
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
		if !w.IsOnline || w.Client == nil {
			continue
		}

		users, err := w.Client.GetActiveUsers()
		if err != nil {
			logger.Error("Failed to fetch active users", zap.String("host", w.Router.Host), zap.Error(err))
		} // Continue anyway to try fetching resources
		
		res, errRes := w.Client.GetSystemResource()
		if errRes != nil {
			// logging error
		}

		// Update Cache
		w.Lock.Lock()
		if err == nil {
			w.ActiveUsers = users
		}
		if errRes == nil {
			w.SystemResource = res
		}
		w.Lock.Unlock()
		
		// logger.Info("Metrics refreshed", zap.String("host", w.Router.Host), zap.Int("users", len(users)))
	}
}
