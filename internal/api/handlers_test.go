package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Expert Go developers use 'httptest' to record real responses 
// without needing to run the actual server.

func TestHealthCheck(t *testing.T) {
	// 1. Setup - Create a fresh Gin engine for testing
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", HealthCheck)

	// 2. Execute - Create a fake request
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	
	// Serve the request
	r.ServeHTTP(w, req)

	// 3. Verify (Assert)
	// Check Status Code
	assert.Equal(t, http.StatusOK, w.Code)
	
	// Check Body
	expected := `{"muscle":"alive","status":"ok"}`
	assert.JSONEq(t, expected, w.Body.String())
}
