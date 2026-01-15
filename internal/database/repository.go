package database

import (
	"skynet-net-engine-api/internal/models"
	"go.uber.org/zap"
	"skynet-net-engine-api/pkg/logger"
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
func UpsertUser(username string, routerID int, profile string) error {
	query := `
		INSERT INTO pppoe_users (username, router_id, profile, is_enabled)
		VALUES (?, ?, ?, true)
		ON DUPLICATE KEY UPDATE
			profile = VALUES(profile),
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := DB.Exec(query, username, routerID, profile)
	if err != nil {
		logger.Error("Failed to upsert user", zap.String("user", username), zap.Error(err))
	}
	return err
}

// GetUsersByRouter fetches all users for a specific router
func GetUsersByRouter(routerID int) (map[string]string, error) {
	rows, err := DB.Query("SELECT username, profile FROM pppoe_users WHERE router_id = ? AND is_enabled = true", routerID)
	if err != nil {
		logger.Error("Failed to fetch users by router", zap.Int("router_id", routerID), zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	users := make(map[string]string) // username -> profile
	for rows.Next() {
		var username, profile string
		if err := rows.Scan(&username, &profile); err != nil {
			continue
		}
		users[username] = profile
	}

	return users, nil
}
