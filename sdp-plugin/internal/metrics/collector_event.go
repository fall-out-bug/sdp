package metrics

func (c *Collector) processEvents(events []evidenceEvent, processedEventIDs map[string]bool, metrics *Metrics) (
	wsState map[string]*workstreamState, modelStats map[string]*modelVerificationStats, newProcessedIDs map[string]bool,
) {
	wsState = make(map[string]*workstreamState)
	modelStats = make(map[string]*modelVerificationStats)
	newProcessedIDs = make(map[string]bool)

	for _, ev := range events {
		if processedEventIDs != nil && processedEventIDs[ev.ID] {
			continue
		}
		newProcessedIDs[ev.ID] = true

		c.processEvent(ev, metrics, wsState, modelStats)
	}

	return wsState, modelStats, newProcessedIDs
}

func (c *Collector) processEvent(ev evidenceEvent, metrics *Metrics, wsState map[string]*workstreamState, modelStats map[string]*modelVerificationStats) {
	switch ev.Type {
	case "verification":
		c.processVerification(ev, metrics, wsState, modelStats)
	case "approval":
		c.processApproval(ev, metrics)
	case "generation":
		c.processGeneration(ev, wsState)
	}
}

func (c *Collector) computeModelPassRates(modelStats map[string]*modelVerificationStats, metrics *Metrics) {
	for modelID, stats := range modelStats {
		if stats.Total > 0 {
			modelStats[modelID].PassRate = float64(stats.Passed) / float64(stats.Total)
		}
		metrics.ModelPassRate[modelID] = stats.PassRate
	}
}

func (c *Collector) computeCatchRates(metrics *Metrics) {
	if metrics.TotalVerifications > 0 {
		metrics.CatchRate = float64(metrics.FailedVerifications) / float64(metrics.TotalVerifications)
	}
	if metrics.TotalApprovals > 0 {
		metrics.AcceptanceCatchRate = float64(metrics.FailedApprovals) / float64(metrics.TotalApprovals)
	}
}

func (c *Collector) processVerification(ev evidenceEvent, metrics *Metrics, wsState map[string]*workstreamState, modelStats map[string]*modelVerificationStats) {
	metrics.TotalVerifications++

	passed, ok := ev.Data["passed"].(bool)
	if !ok {
		return
	}

	if !passed {
		metrics.FailedVerifications++
	}

	// Track model pass rate
	modelID := c.wsModel[ev.WSID]
	if modelID != "" {
		stats, exists := modelStats[modelID]
		if !exists {
			stats = &modelVerificationStats{}
			modelStats[modelID] = stats
		}
		stats.Total++
		if passed {
			stats.Passed++
		}
	}

	// Track iteration count
	if _, exists := metrics.IterationCount[ev.WSID]; !exists {
		metrics.IterationCount[ev.WSID] = 0
	}
	metrics.IterationCount[ev.WSID]++
}

func (c *Collector) processApproval(ev evidenceEvent, metrics *Metrics) {
	metrics.TotalApprovals++

	approved, ok := ev.Data["approved"].(bool)
	if !ok {
		return
	}

	if !approved {
		metrics.FailedApprovals++
	}
}

func (c *Collector) processGeneration(ev evidenceEvent, wsState map[string]*workstreamState) {
	state := wsState[ev.WSID]
	if state == nil {
		state = &workstreamState{}
		wsState[ev.WSID] = state
	}
	state.generationCount++
	state.lastWasGen = true

	// Extract model ID from generation data
	modelID, ok := ev.Data["model_id"].(string)
	if ok && modelID != "" {
		c.wsModel[ev.WSID] = modelID
	}
}
