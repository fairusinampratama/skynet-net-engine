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
