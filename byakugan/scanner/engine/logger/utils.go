package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Helper function to merge fields
func mergeFields(base, override map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(base)+len(override))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		result[k] = v
	}
	return result
}

// Get caller info for better logging context
func getCallerInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

// Format error fields consistently
func formatError(err error) map[string]interface{} {
	if err == nil {
		return nil
	}
	return map[string]interface{}{
		"error":      err.Error(),
		"error_type": fmt.Sprintf("%T", err),
	}
}
