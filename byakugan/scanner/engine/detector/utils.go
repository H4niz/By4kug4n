package detector

import (
	"fmt"
	"net/http"
	"time"
)

// collectRequestData gathers data from request
func collectRequestData(req *http.Request) map[string]string {
	data := make(map[string]string)
	if req != nil {
		data["request_method"] = req.Method
		data["request_url"] = req.URL.String()
		for k, v := range req.Header {
			if len(v) > 0 {
				data[fmt.Sprintf("request_header_%s", k)] = sanitizeHeaderValue(v[0])
			}
		}
	}
	return data
}

// collectResponseData gathers data from response
func collectResponseData(resp *http.Response, body []byte) map[string]string {
	data := make(map[string]string)
	if resp != nil {
		data["response_status"] = fmt.Sprintf("%d", resp.StatusCode)
		for k, v := range resp.Header {
			if len(v) > 0 {
				data[fmt.Sprintf("response_header_%s", k)] = sanitizeHeaderValue(v[0])
			}
		}
		if body != nil {
			data["response_body_length"] = fmt.Sprintf("%d", len(body))
		}
	}
	return data
}

// sanitizeHeaderValue cleans header values
func sanitizeHeaderValue(value string) string {
	clean := []rune{}
	for _, r := range value {
		if r >= 32 && r < 127 {
			clean = append(clean, r)
		}
	}
	return string(clean)
}

func generateID() string {
	return fmt.Sprintf("FINDING-%d", time.Now().UnixNano())
}

func calculateConfidence(pattern *DetectionPattern, evidence *Evidence) float64 {
	confidence := pattern.Confidence

	if evidence.ResponseTime > 0 {
		confidence *= 1.2
	}
	if len(evidence.Observations) > 0 {
		confidence *= 1.1
	}

	if confidence > 1.0 {
		confidence = 1.0
	}

	return confidence
}
