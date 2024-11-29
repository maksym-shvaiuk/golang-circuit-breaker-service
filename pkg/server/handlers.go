package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maksym-shvaiuk/circuit-breaker-golang-test-excercise/pkg/model"
)

// updateConfig updates the configuration of a circuit breaker.
func updateConfig(c *gin.Context) {
	service, err := getServiceSafely(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service instance"})
		return
	}

	deviceIDStr := c.Param("deviceID")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deviceID"})
		return
	}

	var req model.CircuitBreakerEntry
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	req.DeviceID = model.Key(deviceID)
	err = service.storage.UpsertEntry(c.Request.Context(), req.DeviceID, req)
	if err != nil {
		service.logger.Error("Failed to upsert entry", "deviceID", deviceID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update config"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deviceID": req.DeviceID, "config": req})
}

// resetCircuitBreaker resets a circuit breaker to the CLOSED state.
func resetCircuitBreaker(c *gin.Context) {
	service, err := getServiceSafely(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service instance"})
		return
	}

	deviceIDStr := c.Param("deviceID")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deviceID"})
		return
	}

	entry, err := service.storage.GetEntry(c.Request.Context(), model.Key(deviceID))
	if err != nil {
		service.logger.Error("Failed to get entry", "deviceID", deviceID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	entry.State = model.StateClosed
	entry.LastChanged = time.Now()

	err = service.storage.UpsertEntry(c.Request.Context(), entry.DeviceID, entry)
	if err != nil {
		service.logger.Error("Failed to reset entry", "deviceID", deviceID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reset circuit breaker"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deviceID": entry.DeviceID, "newState": entry.State})
}

// getCircuitBreakerStatus retrieves the status of a specific circuit breaker.
func getCircuitBreakerStatus(c *gin.Context) {
	service, err := getServiceSafely(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service instance"})
		return
	}

	deviceIDStr := c.Param("deviceID")
	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deviceID"})
		return
	}

	entry, err := service.storage.GetEntry(c.Request.Context(), model.Key(deviceID))
	if err != nil {
		service.logger.Error("Failed to get entry", "deviceID", deviceID, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

// getAllCircuitBreakers retrieves all circuit breakers with optional pagination.
func getAllCircuitBreakers(c *gin.Context) {
	service, err := getServiceSafely(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve service instance"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	allEntries, err := service.storage.GetAllEntries(c.Request.Context())
	if err != nil {
		service.logger.Error("Failed to get all entries", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve circuit breakers"})
		return
	}

	// Apply pagination
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(allEntries) {
		start, end = len(allEntries), len(allEntries)
	}
	if end > len(allEntries) {
		end = len(allEntries)
	}

	paginatedEntries := allEntries[start:end]

	c.JSON(http.StatusOK, gin.H{
		"page":            page,
		"pageSize":        pageSize,
		"totalItems":      len(allEntries),
		"totalPages":      (len(allEntries) + pageSize - 1) / pageSize,
		"circuitBreakers": paginatedEntries,
	})
}
