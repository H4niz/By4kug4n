package detector

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/haniz/byakugan/scanner/engine/logger"
	pb "github.com/haniz/byakugan/scanner/proto"
)

type Detector interface {
	DetectVulnerability(ctx context.Context, dc *DetectionContext) (*pb.Evidence, error)
}

type detectorImpl struct {
	config   *Config
	logger   logger.Logger
	analyzer *Analyzer
	client   *http.Client
}

// Function to handle Analyzer
func NewDetector(config *Config, logger logger.Logger) Detector {
	if config == nil {
		config = DefaultConfig()
	}

	return &detectorImpl{
		config:   config,
		logger:   logger,
		analyzer: NewAnalyzer(config, logger),
		client: &http.Client{
			Timeout: config.RequestTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if !config.FollowRedirects {
					return http.ErrUseLastResponse
				}
				if len(via) >= config.MaxRedirects {
					return fmt.Errorf("stopped after %d redirects", config.MaxRedirects)
				}
				return nil
			},
		},
	}
}

func (d *detectorImpl) DetectVulnerability(ctx context.Context, dc *DetectionContext) (*pb.Evidence, error) {
	if d.logger != nil {
		d.logger.Info("Starting vulnerability detection", map[string]interface{}{
			"target":   dc.Target.Url,
			"method":   dc.Target.Method,
			"location": dc.InsertionPoint.Location,
		})
	}

	// Validate inputs
	if err := d.validateContext(dc); err != nil {
		return nil, fmt.Errorf("invalid detection context: %w", err)
	}

	// Create and execute request
	req, err := d.createHTTPRequest(ctx, dc)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, body, duration, err := d.executeRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request execution failed: %w", err)
	}

	// Create evidence
	evidence := &pb.Evidence{
		Request: &pb.HttpRequest{
			Url:     dc.Target.Url,
			Method:  dc.Target.Method,
			Headers: d.convertHeaders(req.Header),
		},
		Response: &pb.HttpResponse{
			StatusCode:   int32(resp.StatusCode),
			Headers:      d.convertHeaders(resp.Header),
			Body:         string(body),
			ResponseTime: duration.Milliseconds(),
		},
		Validated: d.validateEvidence(resp, body, dc.Validation),
		Timestamp: time.Now().Unix(),
	}

	if d.logger != nil {
		d.logger.Info("Vulnerability detection completed", map[string]interface{}{
			"status_code": evidence.Response.StatusCode,
			"validated":   evidence.Validated,
			"duration":    duration,
		})
	}

	return evidence, nil
}

func (d *detectorImpl) createHTTPRequest(ctx context.Context, dc *DetectionContext) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, dc.Target.Method, dc.Target.Url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers
	for k, v := range dc.AuthContext.Headers {
		req.Header.Set(k, v)
	}

	// Insert payload at insertion point
	if err := d.insertPayload(req, dc.InsertionPoint, dc.Payload); err != nil {
		return nil, err
	}

	return req, nil
}

func (d *detectorImpl) executeRequest(req *http.Request) (*http.Response, []byte, time.Duration, error) {
	start := time.Now()

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, 0, err
	}

	duration := time.Since(start)
	return resp, body, duration, nil
}

func (d *detectorImpl) validateEvidence(resp *http.Response, body []byte, validation *pb.Validation) bool {
	if validation == nil || validation.SuccessConditions == nil {
		return false
	}

	// Status code validation
	validStatus := false
	for _, code := range validation.SuccessConditions.StatusCodes {
		if int32(resp.StatusCode) == code {
			validStatus = true
			break
		}
	}

	// Response pattern validation
	validPattern := false
	bodyStr := string(body)
	for _, pattern := range validation.SuccessConditions.ResponsePatterns {
		if strings.Contains(bodyStr, pattern) {
			validPattern = true
			break
		}
	}

	return validStatus && validPattern
}

func (d *detectorImpl) insertPayload(req *http.Request, point *pb.InsertionPoint, payload string) error {
	parts := strings.Split(point.Location, ".")
	if len(parts) != 2 {
		return fmt.Errorf("invalid insertion point location: %s", point.Location)
	}

	location := parts[0]
	param := parts[1]

	switch location {
	case "header":
		req.Header.Set(param, payload)
	case "query":
		q := req.URL.Query()
		q.Set(param, payload)
		req.URL.RawQuery = q.Encode()
	case "path":
		// Handle path parameter replacement
		req.URL.Path = strings.ReplaceAll(req.URL.Path, fmt.Sprintf("{%s}", param), payload)
	default:
		return fmt.Errorf("unsupported insertion point location: %s", location)
	}

	return nil
}

func (d *detectorImpl) validateContext(dc *DetectionContext) error {
	if dc == nil {
		return fmt.Errorf("invalid detection context: context is nil")
	}

	// Validate Target
	if dc.Target == nil {
		return fmt.Errorf("invalid URL: target is nil")
	}

	if dc.Target.Url == "" {
		return fmt.Errorf("invalid URL: URL is empty")
	}

	// Validate URL format
	parsedURL, err := url.Parse(dc.Target.Url)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}

	// Additional URL validation
	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid URL: missing scheme or host")
	}

	// Validate required fields
	if dc.AuthContext == nil {
		dc.AuthContext = &pb.AuthContext{}
	}

	if dc.InsertionPoint == nil {
		dc.InsertionPoint = &pb.InsertionPoint{}
	}

	if dc.Validation == nil {
		dc.Validation = &pb.Validation{}
	}

	return nil
}

func (d *detectorImpl) convertHeaders(headers http.Header) map[string]string {
	converted := make(map[string]string)
	for k, v := range headers {
		if len(v) > 0 {
			converted[k] = v[0]
		}
	}
	return converted
}

// sanitizeRequest cleans and validates request
func sanitizeRequest(req *http.Request) (*http.Request, error) {
	if req == nil {
		return nil, ErrInvalidRequest
	}
	// Clone request to avoid modifying original
	clone := req.Clone(req.Context())

	// Sanitize headers
	for key := range clone.Header {
		if len(clone.Header[key]) > 0 {
			clone.Header[key][0] = sanitizeHeaderValue(clone.Header[key][0])
		}
	}

	return clone, nil
}

// sanitizeResponse cleans and validates response
func sanitizeResponse(resp *http.Response) (*http.Response, error) {
	if resp == nil {
		return nil, ErrInvalidResponse
	}
	// Clone response to avoid modifying original
	clone := &http.Response{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Header:     make(http.Header),
	}

	// Copy and sanitize headers
	for key, values := range resp.Header {
		clone.Header[key] = make([]string, len(values))
		for i, v := range values {
			clone.Header[key][i] = sanitizeHeaderValue(v)
		}
	}

	return clone, nil
}
