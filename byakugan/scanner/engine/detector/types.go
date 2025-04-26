package detector

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	pb "github.com/haniz/byakugan/scanner/proto"
)

// Severity represents vulnerability severity level
type Severity string

const (
	Critical Severity = "CRITICAL"
	High     Severity = "HIGH"
	Medium   Severity = "MEDIUM"
	Low      Severity = "LOW"
	Info     Severity = "INFO"
)

// DetectionStrategy defines scanning strategy interface
type DetectionStrategy interface {
	Name() string
	Analyze(req *http.Request, resp *http.Response, body []byte) (*Finding, error)
}

// Evidence represents detection evidence
type Evidence struct {
	Request      *http.Request
	Response     *http.Response
	ResponseBody []byte
	ResponseTime time.Duration
	Data         map[string]string
	Description  string // Add Description field
	Timestamp    int64
	MatchedRules []string
	Observations map[string]string
	Validated    bool
}

// Finding represents a vulnerability finding
type Finding struct {
	ID          string
	RuleID      string
	Type        string
	Pattern     string
	Severity    string
	Confidence  float64
	Evidence    *Evidence
	Description string
	Timestamp   time.Time
	Metadata    map[string]string
}

// DetectionPattern defines pattern matching configuration
type DetectionPattern struct {
	ID          string
	Pattern     string
	Type        string
	Description string
	Confidence  float64
}

// RequestContext contains request analysis context
type RequestContext struct {
	Request     *http.Request
	Response    *http.Response
	Body        []byte
	StartTime   time.Time
	InsertPoint string
	Payload     string
}

// Statistics tracks detection metrics
type Statistics struct {
	TotalScans     int64 `json:"total_scans"`
	TotalFindings  int64 `json:"total_findings"`
	RulesMatched   int64 `json:"rules_matched"`
	FalsePositives int64 `json:"false_positives"`
	scanDurations  []time.Duration
	mu             sync.RWMutex
}

// ValidationConfig defines validation configuration
type ValidationConfig struct {
	SuccessConditions *pb.SuccessConditions
	ResponseTimeout   time.Duration
	RetryCount        int
}

// DetectionContext contains scan context data
type DetectionContext struct {
	Target         *pb.Target
	AuthContext    *pb.AuthContext
	InsertionPoint *pb.InsertionPoint
	Payload        string
	Validation     *pb.Validation
}

// Validate validates the detection context
func (dc *DetectionContext) Validate() error {
	if dc == nil {
		return fmt.Errorf("detection context is nil")
	}

	if dc.Target == nil {
		return fmt.Errorf("target is missing")
	}

	if dc.Target.Url == "" {
		return fmt.Errorf("target URL is missing")
	}

	return nil
}
