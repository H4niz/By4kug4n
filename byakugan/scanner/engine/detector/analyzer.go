package detector

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	"github.com/haniz/byakugan/scanner/engine/logger"
)

func (s *Statistics) averageDuration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(s.scanDurations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range s.scanDurations {
		total += d
	}
	return total / time.Duration(len(s.scanDurations))
}

// Analyzer handles response analysis for vulnerability detection
type Analyzer struct {
	config     *Config
	logger     logger.Logger
	patterns   []*DetectionPattern
	strategies []DetectionStrategy
	evidence   *EvidenceCollector
	rules      map[string]*Rule
	analyzers  map[string]*RuleAnalyzer
	stats      *Statistics // Add stats field
	mu         sync.RWMutex
}

type AnalyzerMetrics struct {
	TotalScans     int64         `json:"total_scans"`
	TotalFindings  int64         `json:"total_findings"`
	ScanDuration   time.Duration `json:"scan_duration"`
	RulesMatched   int64         `json:"rules_matched"`
	FalsePositives int64         `json:"false_positives"`
}

type RuleAnalyzer struct {
	enabled bool
	name    string
	fn      AnalyzerFunc
}

type AnalyzerFunc func(req *http.Request, resp *http.Response) ([]Finding, error)

// Initialize analyzer with evidence collector
func NewAnalyzer(config *Config, logger logger.Logger) *Analyzer {
	if config == nil {
		config = DefaultConfig()
	}

	return &Analyzer{
		config:     config,
		logger:     logger,
		patterns:   make([]*DetectionPattern, 0),
		strategies: make([]DetectionStrategy, 0),
		evidence:   NewEvidenceCollector(config),
		rules:      make(map[string]*Rule),
		analyzers:  make(map[string]*RuleAnalyzer),
		stats:      &Statistics{},
	}
}

// AddPattern adds a detection pattern
func (a *Analyzer) AddPattern(pattern *DetectionPattern) error {
	if pattern == nil {
		return ErrInvalidPattern
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.patterns = append(a.patterns, pattern)
	return nil
}

// AddStrategy adds a detection strategy
func (a *Analyzer) AddStrategy(strategy DetectionStrategy) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.strategies = append(a.strategies, strategy)
}

// Update Analyze method
func (a *Analyzer) Analyze(req *http.Request, resp *http.Response) ([]Finding, error) {
	startTime := time.Now()
	findings := make([]Finding, 0)

	if resp == nil || resp.Body == nil {
		return nil, ErrInvalidResponse
	}

	// Read response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Collect initial evidence
	evidence, err := a.CollectEvidence(req, resp, body)
	if err != nil {
		if a.logger != nil {
			a.logger.Error("Evidence collection failed", err, nil)
		}
		// Continue with basic evidence
		evidence = &Evidence{
			Request:      req,
			Response:     resp,
			ResponseBody: body,
			ResponseTime: time.Since(startTime),
			Data:         make(map[string]string),
			Timestamp:    time.Now().Unix(),
		}
	}

	// Pattern matching
	for _, p := range a.patterns {
		re, err := regexp.Compile(p.Pattern)
		if err != nil {
			if a.logger != nil {
				a.logger.Error("Invalid pattern", err, map[string]interface{}{
					"pattern": p.Pattern,
				})
			}
			continue
		}

		matches := re.FindAll(body, -1)
		if len(matches) > 0 {
			finding := Finding{
				ID:          generateID(),
				Type:        p.Type,
				Pattern:     p.Pattern,
				Confidence:  calculateConfidence(p, evidence),
				Evidence:    evidence,
				Timestamp:   time.Now(),
				Description: p.Description,
			}
			findings = append(findings, finding)
		}
	}

	a.updateStats(time.Since(startTime), len(findings))
	return findings, nil
}

func (a *Analyzer) RulesCount() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.rules)
}

func (a *Analyzer) EnabledAnalyzers() []string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	var enabled []string
	for name, analyzer := range a.analyzers {
		if analyzer.enabled {
			enabled = append(enabled, name)
		}
	}
	return enabled
}

// Add rule registration method
func (a *Analyzer) RegisterRule(rule *Rule) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	if rule.ID == "" {
		return fmt.Errorf("rule ID cannot be empty")
	}

	a.rules[rule.ID] = rule
	return nil
}

// Add analyzer registration method
func (a *Analyzer) RegisterAnalyzer(name string, fn AnalyzerFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.analyzers[name] = &RuleAnalyzer{
		enabled: true,
		name:    name,
		fn:      fn,
	}
}

// Rule represents a security detection rule
type Rule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	CVSS        float64  `json:"cvss_score"`
	CWE         string   `json:"cwe"`
	References  []string `json:"references"`

	// Detection configuration
	Patterns []struct {
		Pattern string `json:"pattern"`
		Type    string `json:"type"`
	} `json:"patterns"`

	Timeouts struct {
		Connect int `json:"connect"`
		Read    int `json:"read"`
		Write   int `json:"write"`
	} `json:"timeouts"`

	RateLimit int `json:"rate_limit"`

	// Validation criteria
	Validation struct {
		StatusCodes []int    `json:"status_codes"`
		Headers     []string `json:"headers"`
		Content     []string `json:"content"`
	} `json:"validation"`
}

// Validate checks if a rule is valid
func (r *Rule) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("rule ID cannot be empty")
	}

	if r.Type == "" {
		return fmt.Errorf("rule type cannot be empty")
	}

	if len(r.Patterns) == 0 {
		return fmt.Errorf("rule must have at least one pattern")
	}

	// Validate CVSS score
	if r.CVSS < 0 || r.CVSS > 10 {
		return fmt.Errorf("invalid CVSS score: must be between 0 and 10")
	}

	return nil
}

// MatchPattern checks if a rule pattern matches content
func (r *Rule) MatchPattern(content []byte) (bool, error) {
	for _, p := range r.Patterns {
		re, err := regexp.Compile(p.Pattern)
		if err != nil {
			return false, fmt.Errorf("invalid pattern: %w", err)
		}

		if re.Match(content) {
			return true, nil
		}
	}
	return false, nil
}

// ValidateResponse checks if a response matches the rule criteria
func (r *Rule) ValidateResponse(resp *http.Response) bool {
	// Check status code
	if len(r.Validation.StatusCodes) > 0 {
		found := false
		for _, code := range r.Validation.StatusCodes {
			if resp.StatusCode == code {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Check headers
	for _, header := range r.Validation.Headers {
		if resp.Header.Get(header) == "" {
			return false
		}
	}

	return true
}

// LogMetrics logs analyzer metrics
func (a *Analyzer) LogMetrics() {
	if a.logger == nil {
		return
	}

	metrics := &AnalyzerMetrics{
		TotalScans:     atomic.LoadInt64(&a.stats.TotalScans),
		TotalFindings:  atomic.LoadInt64(&a.stats.TotalFindings),
		ScanDuration:   a.stats.averageDuration(),
		RulesMatched:   atomic.LoadInt64(&a.stats.RulesMatched),
		FalsePositives: atomic.LoadInt64(&a.stats.FalsePositives),
	}

	a.logger.LogMetrics("detector", map[string]interface{}{
		"total_scans":     metrics.TotalScans,
		"total_findings":  metrics.TotalFindings,
		"scan_duration":   metrics.ScanDuration,
		"rules_matched":   metrics.RulesMatched,
		"false_positives": metrics.FalsePositives,
	})
}

// Update statistics
func (a *Analyzer) updateStats(duration time.Duration, findings int) {
	atomic.AddInt64(&a.stats.TotalScans, 1)
	atomic.AddInt64(&a.stats.TotalFindings, int64(findings))

	a.stats.mu.Lock()
	defer a.stats.mu.Unlock()

	a.stats.scanDurations = append(a.stats.scanDurations, duration)
	if len(a.stats.scanDurations) > 100 {
		a.stats.scanDurations = a.stats.scanDurations[1:]
	}
}

func (a *Analyzer) analyze(pattern *DetectionPattern, evidence *Evidence) (*Finding, error) {
	confidence := calculateConfidence(pattern, evidence)

	return &Finding{
		ID:         generateID(),
		Type:       pattern.Type,
		Pattern:    pattern.Pattern,
		Confidence: confidence,
		Evidence:   evidence,
		Timestamp:  time.Now(),
	}, nil
}

// Update CollectEvidence to properly handle response time
func (a *Analyzer) CollectEvidence(req *http.Request, resp *http.Response, body []byte) (*Evidence, error) {
	startTime := time.Now()

	evidence := &Evidence{
		Request:      req,
		Response:     resp,
		ResponseBody: body,
		Data:         make(map[string]string),
		Timestamp:    time.Now().Unix(),
	}

	// Ensure timing measurement
	if resp != nil {
		evidence.ResponseTime = time.Since(startTime)
		evidence.Data["response_time_ms"] = fmt.Sprintf("%d", evidence.ResponseTime.Milliseconds())
	}

	// Add response data
	if resp != nil && body != nil {
		evidence.Data["response_body_length"] = fmt.Sprintf("%d", len(body))
		evidence.Data["status_code"] = fmt.Sprintf("%d", resp.StatusCode)
		evidence.Data["content_type"] = resp.Header.Get("Content-Type")
	}

	// Set validation status
	evidence.Validated = true

	return evidence, nil
}
