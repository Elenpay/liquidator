package customerrors

// Error type for when a swap is in progress
type SwapInProgressError struct {
	Message string
}

func (e *SwapInProgressError) Error() string {
	return e.Message
}
