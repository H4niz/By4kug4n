package logger

import "time"

// UnifiedLoggerAdapter provides a standardized logging interface for all scanner components
type UnifiedLoggerAdapter struct {
	engineLogger Logger
	commsLogger  Logger
}

func NewUnifiedLogger(engineLogger, commsLogger Logger) *UnifiedLoggerAdapter {
	return &UnifiedLoggerAdapter{
		engineLogger: engineLogger,
		commsLogger:  commsLogger,
	}
}

// Core logging methods
func (l *UnifiedLoggerAdapter) Info(msg string, fields map[string]interface{}) {
	if l.engineLogger != nil {
		l.engineLogger.Info(msg, fields)
	}
	if l.commsLogger != nil {
		l.commsLogger.Info(msg, fields)
	}
}

func (l *UnifiedLoggerAdapter) Error(msg string, err error, fields map[string]interface{}) {
	mergedFields := make(map[string]interface{})
	if fields != nil {
		for k, v := range fields {
			mergedFields[k] = v
		}
	}
	if err != nil {
		mergedFields["error"] = err.Error()
	}

	l.Info(msg, mergedFields)
}

// Component specific logging
func (l *UnifiedLoggerAdapter) LogScan(scanID string, fields map[string]interface{}) {
	mergedFields := map[string]interface{}{
		"scan_id":   scanID,
		"type":      "scan",
		"timestamp": time.Now().Format(time.RFC3339),
	}
	for k, v := range fields {
		mergedFields[k] = v
	}
	l.Info("Scan execution", mergedFields)
}

func (l *UnifiedLoggerAdapter) LogPayload(payloadID string, fields map[string]interface{}) {
	mergedFields := map[string]interface{}{
		"payload_id": payloadID,
		"type":       "payload",
	}
	for k, v := range fields {
		mergedFields[k] = v
	}
	l.Info("Payload generation", mergedFields)
}

func (l *UnifiedLoggerAdapter) LogVulnerability(vulnID string, fields map[string]interface{}) {
	mergedFields := map[string]interface{}{
		"vulnerability_id": vulnID,
		"type":             "vulnerability",
	}
	for k, v := range fields {
		mergedFields[k] = v
	}
	l.Info("Vulnerability detected", mergedFields)
}

type loggerAdapter struct {
	baseLogger Logger
}

func NewLoggerAdapter(base Logger) Logger {
	return &loggerAdapter{
		baseLogger: base,
	}
}

// Implement full Logger interface
func (a *loggerAdapter) Info(msg string, fields map[string]interface{}) {
	a.baseLogger.Info(msg, fields)
}

func (a *loggerAdapter) Error(msg string, err error, fields map[string]interface{}) {
	a.baseLogger.Error(msg, err, fields)
}

func (a *loggerAdapter) Debug(msg string, fields map[string]interface{}) {
	a.baseLogger.Debug(msg, fields)
}

func (a *loggerAdapter) WithField(key string, value interface{}) Logger {
	return a.baseLogger.WithField(key, value)
}

func (a *loggerAdapter) WithFields(fields map[string]interface{}) Logger {
	return a.baseLogger.WithFields(fields)
}

func (a *loggerAdapter) LogMetrics(component string, metrics map[string]interface{}) {
	a.baseLogger.LogMetrics(component, metrics)
}

func (a *loggerAdapter) LogPayload(id string, payload interface{}) {
	a.baseLogger.LogPayload(id, payload)
}

func (a *loggerAdapter) LogRequest(method string, url string, headers map[string]interface{}, params map[string]interface{}) {
	a.baseLogger.LogRequest(method, url, headers, params)
}

func (a *loggerAdapter) LogResponse(statusCode int, body map[string]interface{}, duration time.Duration) {
	a.baseLogger.LogResponse(statusCode, body, duration)
}
