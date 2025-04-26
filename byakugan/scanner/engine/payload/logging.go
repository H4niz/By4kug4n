package payload

func (g *Generator) LogPayloadGeneration(tmpl *PayloadTemplate, payload *Payload) {
	if g.logger == nil {
		return
	}

	g.logger.LogPayload(payload.ID, map[string]interface{}{
		"template_id":   tmpl.ID,
		"template_type": tmpl.Type,
		"payload_id":    payload.ID,
		"payload_type":  payload.Type,
		"payload_value": payload.Value,
		"encoded":       payload.Encoded,
		"metadata":      payload.Metadata,
	})
}

func (g *Generator) LogMetrics(metrics *GeneratorMetrics) {
	if g.logger == nil {
		return
	}

	g.logger.LogMetrics("payload", map[string]interface{}{
		"total_generated":   metrics.TotalGenerated,
		"failed_generation": metrics.FailedGeneration,
		"generation_time":   metrics.GenerationTime,
		"template_cache":    metrics.TemplateCache,
	})
}
