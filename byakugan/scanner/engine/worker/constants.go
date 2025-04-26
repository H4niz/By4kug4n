package worker

const (
	// Task statuses
	TaskStatusPending   = "PENDING"
	TaskStatusRunning   = "RUNNING"
	TaskStatusCompleted = "COMPLETED"
	TaskStatusFailed    = "FAILED"

	// Default values
	DefaultTimeout = 30 // seconds
	DefaultRetries = 3

	// Error messages
	ErrNilTask     = "nil task provided"
	ErrNilResult   = "nil result provided"
	ErrNilEvidence = "nil evidence provided"
)
