package server

import (
	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	r.PUT("/circuit-breaker/:deviceID/config", updateConfig)
	r.POST("/circuit-breaker/:deviceID/reset", resetCircuitBreaker)
	r.GET("/circuit-breaker/:deviceID/status", getCircuitBreakerStatus)
	r.GET("/circuit-breakers/", getAllCircuitBreakers)
}
