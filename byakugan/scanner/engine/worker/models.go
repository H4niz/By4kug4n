package worker

// Target represents an API target
type Target struct {
	URL      string
	Method   string
	Protocol string
}

// AuthContext represents authentication context
type AuthContext struct {
	Type      string
	Token     string
	ExpiresAt int64
	Headers   map[string]string
}

// Validation represents validation rules
type Validation struct {
	SuccessConditions  SuccessConditions
	EvidenceCollection EvidenceCollection
}

// SuccessConditions defines conditions for successful detection
type SuccessConditions struct {
	StatusCodes      []int32
	ResponsePatterns []string
}

// EvidenceCollection defines what evidence to collect
type EvidenceCollection struct {
	SaveRequest  bool
	SaveResponse bool
	Screenshot   bool
}

// Payload represents a scan payload
type Payload struct {
	Headers         map[string]string
	QueryParams     map[string]string
	Body            []byte
	InsertionPoints []InsertionPoint
}

// InsertionPoint represents where to inject payloads
type InsertionPoint struct {
	Location string
	Type     string
	Payloads []string
	Encoding string
}

// ScanResult represents a scan result
type ScanResult struct {
	TaskID   string
	Success  bool
	Findings []Finding
	Evidence []Evidence
	EndTime  int64
	Metadata map[string]string
}

// Finding represents a vulnerability finding
type Finding struct {
	ID          string
	RuleID      string
	Severity    string
	Evidence    Evidence
	Description string
}

// Evidence represents detection evidence
type Evidence struct {
	Data        []byte
	Description string
	Request     string
	Response    string
}
