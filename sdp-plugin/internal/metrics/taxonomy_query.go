package metrics

// SetClassification manually sets classification (AC5).
func (t *Taxonomy) SetClassification(eventID, failureType, notes string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if fc, exists := t.classifications[eventID]; exists {
		fc.FailureType = failureType
		fc.Notes = notes
	} else {
		t.classifications[eventID] = &FailureClassification{
			EventID:     eventID,
			FailureType: failureType,
			Severity:    t.severityForType(failureType),
			Notes:       notes,
		}
	}
}

// GetClassification retrieves classification by event ID.
func (t *Taxonomy) GetClassification(eventID string) (FailureClassification, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	fc, exists := t.classifications[eventID]
	if !exists {
		return FailureClassification{}, false
	}
	return *fc, true
}

// GetByModel returns all classifications for a model.
func (t *Taxonomy) GetByModel(modelID string) []FailureClassification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []FailureClassification
	for _, fc := range t.classifications {
		if fc.ModelID == modelID {
			result = append(result, *fc)
		}
	}
	return result
}

// GetByType returns all classifications of a failure type.
func (t *Taxonomy) GetByType(failureType string) []FailureClassification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var result []FailureClassification
	for _, fc := range t.classifications {
		if fc.FailureType == failureType {
			result = append(result, *fc)
		}
	}
	return result
}

// GetStats returns summary statistics.
func (t *Taxonomy) GetStats() TaxonomyStats {
	t.mu.RLock()
	defer t.mu.RUnlock()

	stats := TaxonomyStats{
		TotalClassifications: len(t.classifications),
		TotalByModel:         make(map[string]int),
		TotalByType:          make(map[string]int),
		TotalByLanguage:      make(map[string]int),
		TotalBySeverity:      make(map[string]int),
	}

	for _, fc := range t.classifications {
		stats.TotalByModel[fc.ModelID]++
		stats.TotalByType[fc.FailureType]++
		stats.TotalByLanguage[fc.Language]++
		stats.TotalBySeverity[fc.Severity]++
	}

	return stats
}

// GetAll returns all classifications.
func (t *Taxonomy) GetAll() []FailureClassification {
	t.mu.RLock()
	defer t.mu.RUnlock()

	result := make([]FailureClassification, 0, len(t.classifications))
	for _, fc := range t.classifications {
		result = append(result, *fc)
	}
	return result
}
