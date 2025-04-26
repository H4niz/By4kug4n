package detector

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// PatternMatcher implements pattern-based vulnerability detection
type PatternMatcher struct {
	patterns map[string]*regexp.Regexp
}

// NewPatternMatcher creates a new pattern matcher
func NewPatternMatcher() *PatternMatcher {
	return &PatternMatcher{
		patterns: make(map[string]*regexp.Regexp),
	}
}

// LoadPatternsFromFile loads patterns from a JSON file
func (pm *PatternMatcher) LoadPatternsFromFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Match the JSON structure in pattern files
	var patternsFile struct {
		Patterns []*DetectionPattern `json:"patterns"`
	}

	if err := json.Unmarshal(data, &patternsFile); err != nil {
		return err
	}

	for _, p := range patternsFile.Patterns {
		re, err := regexp.Compile(p.Pattern)
		if err != nil {
			continue
		}
		pm.patterns[p.ID] = re
	}

	return nil
}

// Analyze implements the DetectionStrategy interface
func (pm *PatternMatcher) Analyze(req *http.Request, resp *http.Response, body []byte) (*Finding, error) {
	for patternID, re := range pm.patterns {
		if matches := re.FindAll(body, -1); len(matches) > 0 {
			evidence, err := createEvidence(req, resp, body, time.Now())
			if err != nil {
				return nil, err
			}

			evidence.MatchedRules = []string{patternID}

			return &Finding{
				ID:       patternID,
				Type:     "pattern_match",
				Pattern:  re.String(),
				Evidence: evidence,
			}, nil
		}
	}
	return nil, nil
}

func (pm *PatternMatcher) Name() string {
	return "PatternMatcher"
}

// Update createEvidence to handle proper types
func createEvidence(req *http.Request, resp *http.Response, body []byte, startTime time.Time) (*Evidence, error) {
	evidence := &Evidence{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
		ResponseTime: time.Since(startTime),
		Data:         make(map[string]string),
		Description:  "Pattern match evidence",
		Timestamp:    time.Now().Unix(),
		MatchedRules: make([]string, 0),
		Observations: make(map[string]string),
	}

	// Add request data
	requestData := collectRequestData(req)
	for k, v := range requestData {
		evidence.Data[k] = v
	}

	// Add response data
	responseData := collectResponseData(resp, body)
	for k, v := range responseData {
		evidence.Data[k] = v
	}

	return evidence, nil
}

func (p *PatternMatcher) processMatch(evidence *Evidence, pattern *DetectionPattern) {
	if evidence.MatchedRules == nil {
		evidence.MatchedRules = make([]string, 0)
	}
	evidence.MatchedRules = append(evidence.MatchedRules, pattern.ID)
	evidence.Observations["pattern_match"] = pattern.Pattern
}
