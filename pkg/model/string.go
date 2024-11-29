package model

import "fmt"

// String provides a string representation for the State type.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "StateClosed"
	case StateOpen:
		return "StateOpen"
	case StateHalfOpen:
		return "StateHalfOpen"
	default:
		return fmt.Sprintf("UnknownState(%d)", s)
	}
}
