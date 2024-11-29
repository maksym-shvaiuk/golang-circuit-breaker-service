package model

import "time"

// State represents the state of a circuit breaker.
type State int

const (
	// StateClosed means the circuit breaker is allowing requests.
	StateClosed State = iota
	// StateOpen means the circuit breaker is blocking requests.
	StateOpen
	// StateHalfOpen means the circuit breaker is allowing limited requests to test the system's health.
	StateHalfOpen
)

// NOTE: should match the primary key type in selected database package
type Key uint

// CircuitBreakerEntry represents a circuit breaker for a specific device.
type CircuitBreakerEntry struct {
	DeviceID                Key       `json:"deviceID"`
	State                   State     `json:"state"`                   // State of the circuit breaker
	LastChanged             time.Time `json:"lastChanged"`             // Timestamp of the last state change
	ErrorsThreshold         int       `json:"errorsThreshold"`         // Percentage of errors to trip the breaker
	ErrorsCntResetTimeoutMs int       `json:"errorsCntResetTimeoutMs"` // Time in milliseconds to reset the errors count
	ResetTimeoutMs          int       `json:"resetTimeoutMs"`          // Time in milliseconds to reset the breaker
}

// ConfigUpdateRequest represents the payload for updating circuit breaker configuration.
type ConfigUpdateRequest struct {
	ErrorsThreshold         int `json:"errorsThreshold"`
	ErrorsCntResetTimeoutMs int `json:"errorsCntResetTimeoutMs"`
	ResetTimeoutMs          int `json:"resetTimeoutMs"`
}

// ConfigUpdateResponse represents the response for updating circuit breaker configuration.
type ConfigUpdateResponse struct {
	DeviceID string              `json:"deviceID"`
	Config   ConfigUpdateRequest `json:"config"`
}

// ResetResponse represents the response for resetting a circuit breaker.
type ResetResponse struct {
	DeviceID string `json:"deviceID"`
	NewState State  `json:"newState"` // StateClosed
}

// StatusResponse represents the response for retrieving the status of a circuit breaker.
type StatusResponse struct {
	DeviceID    string    `json:"deviceID"`
	State       State     `json:"state"`       // State of the circuit breaker
	LastChanged time.Time `json:"lastChanged"` // Timestamp of the last state change
}

// PaginatedResponse represents the paginated response for retrieving multiple circuit breakers.
type PaginatedResponse struct {
	Page            int                   `json:"page"`
	PageSize        int                   `json:"pageSize"`
	TotalItems      int                   `json:"totalItems"`
	TotalPages      int                   `json:"totalPages"`
	CircuitBreakers []CircuitBreakerEntry `json:"circuitBreakers"`
}
