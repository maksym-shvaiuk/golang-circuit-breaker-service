package server

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// ErrServiceNotFound indicates that the Service instance was not found in the context.
var ErrServiceNotFound = errors.New("service instance not found in context")

// getServiceSafely retrieves the Service instance from the Gin context safely.
func getServiceSafely(c *gin.Context) (*Service, error) {
	service, exists := c.Get("service")
	if !exists {
		return nil, ErrServiceNotFound
	}
	srv, ok := service.(*Service)
	if !ok {
		return nil, ErrServiceNotFound
	}
	return srv, nil
}
