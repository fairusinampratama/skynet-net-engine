package api

import (
	"net/http"
	"strings"
	"io"
	"github.com/gin-gonic/gin"
	"skynet-net-engine-api/pkg/logger"
	"skynet-net-engine-api/internal/assets"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "skynet-net-engine-api/docs" // Import generated docs
)

func Start(port string) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	
	// Middleware
	r.Use(gin.Recovery())
	r.Use(func(c *gin.Context) {
		// CORS for Dev
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-App-Key")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		logger.Info("API Request", 
			logger.Field("method", c.Request.Method),
			logger.Field("path", c.Request.URL.Path),
			logger.Field("ip", c.ClientIP()),
		)
		c.Next()
	})

	// Serve Embedded Frontend
	staticFS, err := assets.GetFS()
	if err != nil {
		logger.Fatal("Failed to load embedded assets")
	}
	
	// API Routes (must come before NoRoute)
	// ... (V1 routes defined below) ...

	// SPA Fallback Handler
	// We want to serve index.html for unknown routes (React Router)
	// BUT we must not interfere with /api/v1
	// SPA Fallback Handler
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"error": "API Route Not Found"})
			return
		}
		
		path := c.Request.URL.Path
		path = strings.TrimPrefix(path, "/")

		// 1. Try to serve if it's a specific static file (JS/CSS)
		if path != "" {
			f, err := staticFS.Open(path)
			if err == nil {
				defer f.Close()
				stat, _ := f.Stat()
				if !stat.IsDir() {
					c.FileFromFS(path, http.FS(staticFS))
					return
				}
			}
		}

		// 2. Fallback to index.html (SPA)
		// We explicitly open and serve content to avoid 301 redirects from FileServer
		f, err := staticFS.Open("index.html")
		if err != nil {
			c.String(500, "Dashboard not found")
			return
		}
		defer f.Close()
		
		stat, _ := f.Stat()
		http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), f.(io.ReadSeeker))
	})

	// Public V1 Routes
	v1 := r.Group("/api/v1")
	{
		// Documentation (Public)
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		v1.GET("/health", HealthCheck) // Health check is often public too
	}

	// Secured V1 Routes
	secured := v1.Group("/")
	secured.Use(func(c *gin.Context) {
		key := c.GetHeader("X-App-Key")
		// Hardcoded for MVP, move to ENV
		if key != "netengine_secret_key_123" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Next()
	})
	{
		// Internal Control
		secured.POST("/sync/:id", SyncRouter)
		secured.POST("/kick", KickUser)

		// CRUD Bridge
		secured.POST("/secret", CreateSecret)
		secured.PUT("/secret/:user", UpdatePlan)
		
		// Advanced
		secured.POST("/isolate", IsolateUser)
		secured.GET("/routers", GetRouters)
		secured.GET("/monitoring/targets", GetTargets)
		secured.GET("/router/:id/health", GetRouterHealth)
		secured.GET("/router/:id/traffic", GetUserTraffic)
		secured.POST("/router/:id/backup", TriggerBackup)
	}

	logger.Info("Starting API Server on " + port)
	if err := r.Run(port); err != nil {
		logger.Fatal("Failed to start API server")
	}
}
