package detector

import (
	"net/http"
	"time"
)

// EvidenceCollector collects vulnerability evidence
type EvidenceCollector struct {
	config *Config
}

// NewEvidenceCollector creates a new evidence collector
func NewEvidenceCollector(config *Config) *EvidenceCollector {
	if config == nil {
		config = DefaultConfig()
	}
	return &EvidenceCollector{config: config}
}

// Add CollectEvidence method
func (ec *EvidenceCollector) CollectEvidence(req *http.Request, resp *http.Response) *Evidence {
	startTime := time.Now()

	evidence := &Evidence{
		Request:      req,
		Response:     resp,
		Data:         make(map[string]string),
		Description:  "Detection evidence",
		Timestamp:    time.Now().Unix(),
		MatchedRules: make([]string, 0),
		Observations: make(map[string]string),
	}

	// Add request data
	if req != nil {
		requestData := collectRequestData(req)
		for k, v := range requestData {
			evidence.Data[k] = v
		}
	}

	// Add response data
	if resp != nil {
		responseData := collectResponseData(resp, nil)
		for k, v := range responseData {
			evidence.Data[k] = v
		}
	}

	// Set response time
	evidence.ResponseTime = time.Since(startTime)
	return evidence
}
