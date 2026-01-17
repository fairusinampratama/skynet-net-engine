package database

import (
	"skynet-net-engine-api/internal/models"
	"go.uber.org/zap"
	"skynet-net-engine-api/pkg/logger"
	"database/sql"
)

func GetAllRouters() ([]models.Router, error) {
	rows, err := DB.Query("SELECT id, name, host, port, username, password FROM routers")
	if err != nil {
		logger.Error("Failed to fetch routers", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var routers []models.Router
	for rows.Next() {
		var r models.Router
		if err := rows.Scan(&r.ID, &r.Name, &r.Host, &r.Port, &r.Username, &r.Password); err != nil {
			logger.Error("Failed to scan router row", zap.Error(err))
			continue
		}
		routers = append(routers, r)
	}

	return routers, nil
}

// UpsertUser inserts or updates a PPPoE user
func UpsertUser(username string, routerID int, profile string, remoteAddress string, isEnabled bool) error {
	query := `
		INSERT INTO pppoe_users (username, router_id, profile, remote_address, is_enabled)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			profile = VALUES(profile),
			remote_address = VALUES(remote_address),
			is_enabled = VALUES(is_enabled),
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := DB.Exec(query, username, routerID, profile, remoteAddress, isEnabled)
	if err != nil {
		logger.Error("Failed to upsert user", zap.String("user", username), zap.Error(err))
	}
	return err
}

// DBUser represents a user record from the database
type DBUser struct {
	Profile       string
	RemoteAddress string
	IsEnabled     bool
}

// GetUsersByRouter fetches all users for a specific router (including disabled ones)
func GetUsersByRouter(routerID int) (map[string]DBUser, error) {
	rows, err := DB.Query("SELECT username, profile, remote_address, is_enabled FROM pppoe_users WHERE router_id = ?", routerID)
	if err != nil {
		logger.Error("Failed to fetch users by router", zap.Int("router_id", routerID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	users := make(map[string]DBUser) // username -> DBUser
	for rows.Next() {
		var username, profile string
		var remoteAddress sql.NullString // Handle potential NULLs
		var isEnabled bool
		
		if err := rows.Scan(&username, &profile, &remoteAddress, &isEnabled); err != nil {
			logger.Error("Scan error", zap.Error(err))
			continue
		}
		users[username] = DBUser{
			Profile:       profile,
			RemoteAddress: remoteAddress.String,
			IsEnabled:     isEnabled,
		}
	}

	return users, nil
}
