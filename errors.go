package main

// Error type for when a swap is in progress
type SwapInProgressError struct {
	message string
}

func (e *SwapInProgressError) Error() string {
	return e.message
}
