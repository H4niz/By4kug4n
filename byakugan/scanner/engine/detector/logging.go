package detector

func (a *Analyzer) LogDetection(finding Finding) {
	if a.logger == nil {
		return
	}

	a.logger.Info("Detection found", map[string]interface{}{
		"finding_id":  finding.ID,
		"type":        finding.Type,
		"confidence":  finding.Confidence,
		"description": finding.Description,
		"evidence":    finding.Evidence,
		"metadata":    finding.Metadata,
	})
}
