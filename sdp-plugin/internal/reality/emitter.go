package reality

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const (
	specVersion            = "v1.0"
	fallbackGeneratedAtUTC = "1970-01-01T00:00:00Z"
	hotspotLineThreshold   = 800
)

type Mode string

const (
	ModeQuick        Mode = "quick"
	ModeDeep         Mode = "deep"
	ModeBootstrapSDP Mode = "bootstrap_sdp"
)

var validFocuses = map[string]bool{
	"":             true,
	"architecture": true,
	"quality":      true,
	"testing":      true,
	"docs":         true,
	"security":     true,
}

type Options struct {
	Mode  Mode
	Focus string
}

type claim struct {
	ClaimID          string   `json:"claim_id"`
	Title            string   `json:"title"`
	Statement        string   `json:"statement"`
	Status           string   `json:"status"`
	Confidence       float64  `json:"confidence"`
	SourceIDs        []string `json:"source_ids"`
	ReviewState      string   `json:"review_state"`
	AffectedPaths    []string `json:"affected_paths,omitempty"`
	OpenQuestions    []string `json:"open_questions,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	AffectedRepos    []string `json:"affected_repos,omitempty"`
	CounterEvidence  []string `json:"counter_evidence,omitempty"`
	AffectedElements []string `json:"affected_components,omitempty"`
}

type source struct {
	SourceID string `json:"source_id"`
	Kind     string `json:"kind"`
	Locator  string `json:"locator"`
	Revision string `json:"revision"`
	Repo     string `json:"repo,omitempty"`
	Path     string `json:"path,omitempty"`
}

type scanResult struct {
	Root               string
	RepoName           string
	GeneratedAt        string
	RepoType           string
	TotalFiles         int
	SourceFiles        int
	TestsFiles         int
	DocFiles           int
	ConfigFiles        int
	ManifestFiles      int
	Modules            []string
	Entrypoints        []string
	HotspotFiles       []hotspot
	Integrations       []integrationObservation
	DocReferenceDrifts []docReferenceDrift
}

type hotspot struct {
	Path  string
	Lines int
}

type integrationObservation struct {
	Name          string
	Kind          string
	ContractType  string
	EvidencePaths []string
	Confidence    float64
}

type docReferenceDrift struct {
	DocPath         string
	ReferencedPath  string
	ObservationNote string
}

// EmitOSS writes the OSS reality artifact set to docs/reality and .sdp/reality.
func EmitOSS(projectRoot string) ([]string, error) {
	return EmitOSSWithOptions(projectRoot, Options{})
}

// EmitOSSWithOptions writes the OSS reality artifact set with an explicit analysis mode.
func EmitOSSWithOptions(projectRoot string, opts Options) ([]string, error) {
	opts = normalizeOptions(opts)
	if err := validateOptions(opts); err != nil {
		return nil, err
	}

	scan, err := scanProject(projectRoot)
	if err != nil {
		return nil, err
	}

	verdict := readinessVerdict(scan)
	dim := readinessDimensions(scan, opts)
	constraints := readinessConstraints(scan)
	claims := buildClaims(scan, verdict, opts)
	sources := buildSources(scan)
	claimIDs := claimIDList(claims)
	topFindingIDs := topFindingClaimIDs(scan, claims)
	integrations := integrationEntries(scan)
	driftFindings := driftFindings(scan)
	bootstrapRecommendations := bootstrapRecommendations(scan, verdict, opts)

	artifactPaths := []string{
		".sdp/reality/reality-summary.json",
		".sdp/reality/feature-inventory.json",
		".sdp/reality/architecture-map.json",
		".sdp/reality/integration-map.json",
		".sdp/reality/quality-report.json",
		".sdp/reality/drift-report.json",
		".sdp/reality/readiness-report.json",
		"docs/reality/summary.md",
		"docs/reality/architecture.md",
		"docs/reality/quality.md",
		"docs/reality/bootstrap.md",
	}

	realitySummary := map[string]any{
		"spec_version":          specVersion,
		"run_id":                runID(scan.RepoName),
		"generated_at":          scan.GeneratedAt,
		"scope":                 summaryScope(scan, opts),
		"readiness_verdict":     verdict,
		"top_finding_claim_ids": topFindingIDs,
		"review_strategy": map[string]any{
			"primary":   "local source-first pass",
			"secondary": "heuristic cross-check across tests, configs, manifests, and docs",
			"scope":     "single_repository",
		},
		"artifacts": artifactPaths,
		"claims":    claims,
		"sources":   sources,
	}

	featureInventory := map[string]any{
		"spec_version": specVersion,
		"generated_at": scan.GeneratedAt,
		"features": []map[string]any{
			{
				"feature_id":         "feature:repository-baseline",
				"title":              fmt.Sprintf("%s baseline", scan.RepoName),
				"summary":            featureSummary(scan, opts),
				"status":             featureStatus(scan),
				"evidence_claim_ids": topFindingIDs,
				"confidence":         confidenceFromScan(scan),
				"mapped_components":  stringSliceOrEmpty(scan.Modules),
			},
		},
		"claims":  claims,
		"sources": sources,
	}

	archNodes := make([]map[string]any, 0, len(scan.Modules))
	for _, module := range scan.Modules {
		archNodes = append(archNodes, map[string]any{
			"node_id":  "module:" + module,
			"name":     module,
			"kind":     "module",
			"boundary": "repository",
			"repo":     scan.RepoName,
			"path":     modulePath(module),
		})
	}
	for _, integration := range scan.Integrations {
		nodeKind := "integration"
		if integration.Kind == "data_store" {
			nodeKind = "data_store"
		}
		archNodes = append(archNodes, map[string]any{
			"node_id":  "integration:" + integration.Name,
			"name":     integration.Name,
			"kind":     nodeKind,
			"boundary": "external",
			"repo":     scan.RepoName,
			"path":     strings.Join(integration.EvidencePaths, ", "),
		})
	}
	archEdges := inferredEdges(scan.Modules)
	archEdges = append(archEdges, inferredIntegrationEdges(scan)...)
	hotspots := make([]map[string]any, 0, len(scan.HotspotFiles))
	for _, h := range scan.HotspotFiles {
		hotspots = append(hotspots, map[string]any{
			"node_id":  "file:" + h.Path,
			"reason":   fmt.Sprintf("%d lines", h.Lines),
			"severity": "medium",
		})
	}
	architectureMap := map[string]any{
		"spec_version": specVersion,
		"generated_at": scan.GeneratedAt,
		"nodes":        archNodes,
		"edges":        archEdges,
		"hotspots":     hotspots,
		"claims":       claims,
		"sources":      sources,
	}

	integrationMap := map[string]any{
		"spec_version": specVersion,
		"generated_at": scan.GeneratedAt,
		"integrations": integrations,
		"claims":       claims,
		"sources":      sources,
	}

	qualityFindings := qualityFindings(scan, claimIDs, opts)
	qualityReport := map[string]any{
		"spec_version":    specVersion,
		"generated_at":    scan.GeneratedAt,
		"analysis_mode":   string(opts.Mode),
		"analysis_focus":  opts.Focus,
		"findings":        qualityFindings,
		"hotspot_ranking": hotspotRanking(scan),
		"claims":          claims,
		"sources":         sources,
	}

	driftReport := map[string]any{
		"spec_version":         specVersion,
		"generated_at":         scan.GeneratedAt,
		"contradictions":       driftFindings,
		"unresolved_questions": unresolvedQuestions(scan, opts),
		"claims":               claims,
		"sources":              sources,
	}

	readinessReport := map[string]any{
		"spec_version":            specVersion,
		"generated_at":            scan.GeneratedAt,
		"verdict":                 verdict,
		"dimensions":              dim,
		"justification_claim_ids": topFindingIDs,
		"constraints":             constraints,
		"suggested_workstreams":   bootstrapRecommendations,
		"claims":                  claims,
		"sources":                 sources,
	}

	jsonOutputs := map[string]any{
		".sdp/reality/reality-summary.json":   realitySummary,
		".sdp/reality/feature-inventory.json": featureInventory,
		".sdp/reality/architecture-map.json":  architectureMap,
		".sdp/reality/integration-map.json":   integrationMap,
		".sdp/reality/quality-report.json":    qualityReport,
		".sdp/reality/drift-report.json":      driftReport,
		".sdp/reality/readiness-report.json":  readinessReport,
	}

	for rel, payload := range jsonOutputs {
		if err := writeJSON(filepath.Join(projectRoot, rel), payload); err != nil {
			return nil, err
		}
	}

	mdOutputs := map[string]string{
		"docs/reality/summary.md":      renderSummaryMD(scan, verdict, topFindingIDs, opts),
		"docs/reality/architecture.md": renderArchitectureMD(scan, opts),
		"docs/reality/quality.md":      renderQualityMD(scan, opts),
		"docs/reality/bootstrap.md":    renderBootstrapMD(scan, verdict, constraints, bootstrapRecommendations, opts),
	}
	for rel, body := range mdOutputs {
		if err := writeText(filepath.Join(projectRoot, rel), body); err != nil {
			return nil, err
		}
	}

	return artifactPaths, nil
}

func scanProject(root string) (scanResult, error) {
	scan := scanResult{
		Root:     root,
		RepoName: filepath.Base(root),
	}
	scan.GeneratedAt = deterministicGeneratedAt(root)

	skipDirs := map[string]bool{
		".git":         true,
		".sdp":         true,
		".beads":       true,
		"node_modules": true,
		"vendor":       true,
	}
	sourceExt := map[string]bool{
		".go":   true,
		".py":   true,
		".js":   true,
		".ts":   true,
		".java": true,
		".rs":   true,
		".sh":   true,
	}

	modules := map[string]bool{}
	entrypoints := map[string]bool{}
	hotspots := make([]hotspot, 0)
	integrations := map[string]integrationObservation{}
	drifts := make([]docReferenceDrift, 0)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if path == root {
			return nil
		}
		name := d.Name()
		if d.IsDir() && skipDirs[name] {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil || !info.Mode().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return nil
		}
		rel = filepath.ToSlash(rel)
		if strings.HasPrefix(rel, "docs/reality/") {
			// Ignore generated human-readable artifacts so reruns are deterministic.
			return nil
		}
		scan.TotalFiles++

		ext := strings.ToLower(filepath.Ext(rel))
		if strings.HasPrefix(rel, "docs/") && ext == ".md" {
			scan.DocFiles++
			for _, drift := range scanDocReferenceDrift(root, rel, path) {
				drifts = append(drifts, drift)
			}
		}
		if kind, ok := classifyConfigOrManifest(rel); ok {
			if kind == "config" {
				scan.ConfigFiles++
			} else {
				scan.ManifestFiles++
			}
		}
		if sourceExt[ext] {
			scan.SourceFiles++
			module := topModule(rel)
			modules[module] = true
			if strings.HasSuffix(rel, "/main.go") || rel == "main.go" {
				entrypoints[rel] = true
			}
			if isTestFile(rel) {
				scan.TestsFiles++
			}
			lines := countLines(path)
			if lines >= hotspotLineThreshold {
				hotspots = append(hotspots, hotspot{Path: rel, Lines: lines})
			}
		}
		if shouldScanForIntegrations(rel, ext) {
			content := readFileForScan(path)
			for _, detected := range detectIntegrations(rel, content) {
				mergeIntegration(integrations, detected)
			}
		}
		return nil
	})
	if err != nil {
		return scan, err
	}

	for m := range modules {
		scan.Modules = append(scan.Modules, m)
	}
	sort.Strings(scan.Modules)
	for ep := range entrypoints {
		scan.Entrypoints = append(scan.Entrypoints, ep)
	}
	sort.Strings(scan.Entrypoints)
	sort.Slice(hotspots, func(i, j int) bool {
		if hotspots[i].Lines == hotspots[j].Lines {
			return hotspots[i].Path < hotspots[j].Path
		}
		return hotspots[i].Lines > hotspots[j].Lines
	})
	scan.HotspotFiles = hotspots
	scan.Integrations = flattenIntegrations(integrations)
	scan.DocReferenceDrifts = drifts
	scan.RepoType = classifyRepoType(scan)

	return scan, nil
}

func deterministicGeneratedAt(root string) string {
	cmd := exec.Command("git", "-C", root, "log", "-1", "--format=%cI")
	out, err := cmd.Output()
	if err != nil {
		return fallbackGeneratedAtUTC
	}
	value := strings.TrimSpace(string(out))
	if value == "" {
		return fallbackGeneratedAtUTC
	}
	return value
}

func readinessVerdict(scan scanResult) string {
	if scan.SourceFiles == 0 {
		return "not_ready"
	}
	if scan.TestsFiles == 0 {
		return "ready_with_constraints"
	}
	if len(scan.DocReferenceDrifts) > 0 {
		return "ready_with_constraints"
	}
	return "ready"
}

func readinessConstraints(scan scanResult) []string {
	constraints := make([]string, 0)
	if scan.SourceFiles == 0 {
		constraints = append(constraints, "No source files detected for analysis.")
	}
	if scan.TestsFiles == 0 {
		constraints = append(constraints, "No test files detected; verification surface is weak.")
	}
	if len(scan.DocReferenceDrifts) > 0 {
		constraints = append(constraints, fmt.Sprintf("%d documented path references do not resolve in the repository.", len(scan.DocReferenceDrifts)))
	}
	if len(scan.Integrations) > 0 {
		constraints = append(constraints, fmt.Sprintf("%d external integration or data-store surfaces need explicit boundary review.", len(scan.Integrations)))
	}
	if len(scan.HotspotFiles) > 0 {
		constraints = append(constraints, "Large hotspot files detected; prefer narrow scoped changes.")
	}
	return constraints
}

func readinessDimensions(scan scanResult, opts Options) map[string]any {
	moduleCount := float64(len(scan.Modules))
	sourceCount := float64(max(scan.SourceFiles, 1))
	hotspotCount := float64(len(scan.HotspotFiles))
	integrationCount := float64(len(scan.Integrations))
	driftCount := float64(len(scan.DocReferenceDrifts))
	verificationScore := clamp(float64(scan.TestsFiles) / sourceCount)
	docTrustScore := clamp(1 - (driftCount / float64(max(scan.DocFiles+1, 1))))
	integrationScore := clamp(1 - (integrationCount / float64(max(len(scan.Modules)+1, 1))))
	if opts.Focus == "testing" && scan.TestsFiles > 0 {
		verificationScore = clamp(verificationScore + 0.1)
	}
	if opts.Focus == "docs" && len(scan.DocReferenceDrifts) == 0 && scan.DocFiles > 0 {
		docTrustScore = clamp(docTrustScore + 0.1)
	}

	return map[string]any{
		"boundary_clarity": map[string]any{
			"score": clamp(moduleCount / 8),
			"note":  fmt.Sprintf("%d top-level modules", len(scan.Modules)),
		},
		"verification_coverage": map[string]any{
			"score": verificationScore,
			"note":  fmt.Sprintf("%d test files / %d source files", scan.TestsFiles, scan.SourceFiles),
		},
		"hotspot_concentration": map[string]any{
			"score": clamp(1 - (hotspotCount / sourceCount)),
			"note":  fmt.Sprintf("%d hotspot files (>= %d lines)", len(scan.HotspotFiles), hotspotLineThreshold),
		},
		"integration_fragility": map[string]any{
			"score": integrationScore,
			"note":  fmt.Sprintf("%d detected integration or data-store surfaces", len(scan.Integrations)),
		},
		"documentation_trust_level": map[string]any{
			"score": docTrustScore,
			"note":  fmt.Sprintf("%d markdown docs under docs/, %d unresolved path references", scan.DocFiles, len(scan.DocReferenceDrifts)),
		},
	}
}

func buildClaims(scan scanResult, verdict string, opts Options) []claim {
	claims := []claim{
		{
			ClaimID:     "claim:source-footprint",
			Title:       "Source footprint detected",
			Statement:   fmt.Sprintf("Detected %d source files across %d modules.", scan.SourceFiles, len(scan.Modules)),
			Status:      "observed",
			Confidence:  0.95,
			SourceIDs:   []string{"source:local-tree"},
			ReviewState: "cross_checked",
			Tags:        []string{"baseline", "structure"},
		},
		{
			ClaimID:     "claim:test-posture",
			Title:       "Test posture baseline",
			Statement:   fmt.Sprintf("Detected %d test files.", scan.TestsFiles),
			Status:      "observed",
			Confidence:  0.9,
			SourceIDs:   []string{"source:local-tree"},
			ReviewState: "cross_checked",
			Tags:        []string{"testing"},
		},
		{
			ClaimID:     "claim:readiness-verdict",
			Title:       "Readiness verdict",
			Statement:   fmt.Sprintf("Emitter baseline classifies repository as %s.", verdict),
			Status:      "inferred",
			Confidence:  0.75,
			SourceIDs:   []string{"source:local-tree", "source:reality-contract"},
			ReviewState: "cross_checked",
			Tags:        []string{"readiness"},
		},
	}
	if len(scan.Integrations) > 0 {
		claims = append(claims, claim{
			ClaimID:     "claim:integration-surface",
			Title:       "Integration surface recovered",
			Statement:   fmt.Sprintf("Detected %d integration or data-store surfaces from local code, config, manifest, and docs evidence.", len(scan.Integrations)),
			Status:      "observed",
			Confidence:  0.8,
			SourceIDs:   []string{"source:local-tree", "source:config-inventory"},
			ReviewState: "cross_checked",
			Tags:        []string{"integration", "architecture"},
		})
	}
	if len(scan.DocReferenceDrifts) > 0 {
		claims = append(claims, claim{
			ClaimID:     "claim:documentation-drift",
			Title:       "Documentation drift detected",
			Statement:   fmt.Sprintf("Detected %d unresolved documentation path references during local cross-check.", len(scan.DocReferenceDrifts)),
			Status:      "conflicted",
			Confidence:  0.78,
			SourceIDs:   []string{"source:docs-inventory", "source:local-tree"},
			ReviewState: "challenged",
			Tags:        []string{"docs", "drift"},
		})
	}
	if opts.Focus != "" {
		claims = append(claims, claim{
			ClaimID:     "claim:analysis-focus",
			Title:       "Focused analysis requested",
			Statement:   fmt.Sprintf("OSS reality emitter applied additional reporting emphasis for %s.", opts.Focus),
			Status:      "observed",
			Confidence:  0.9,
			SourceIDs:   []string{"source:local-tree"},
			ReviewState: "cross_checked",
			Tags:        []string{"focus", opts.Focus},
		})
	}
	return claims
}

func buildSources(scan scanResult) []source {
	sources := []source{
		{
			SourceID: "source:local-tree",
			Kind:     "code",
			Locator:  scan.Root,
			Revision: scan.GeneratedAt,
			Repo:     scan.RepoName,
		},
		{
			SourceID: "source:reality-contract",
			Kind:     "doc",
			Locator:  "docs/specs/reality/ARTIFACT-CONTRACT.md",
			Revision: scan.GeneratedAt,
			Repo:     scan.RepoName,
			Path:     "docs/specs/reality/ARTIFACT-CONTRACT.md",
		},
	}
	if scan.DocFiles > 0 {
		sources = append(sources, source{
			SourceID: "source:docs-inventory",
			Kind:     "doc",
			Locator:  "docs/",
			Revision: scan.GeneratedAt,
			Repo:     scan.RepoName,
			Path:     "docs/",
		})
	}
	if scan.ConfigFiles > 0 || scan.ManifestFiles > 0 {
		sources = append(sources, source{
			SourceID: "source:config-inventory",
			Kind:     "config",
			Locator:  "config+manifest inventory",
			Revision: scan.GeneratedAt,
			Repo:     scan.RepoName,
		})
	}
	return sources
}

func topFindingClaimIDs(scan scanResult, claims []claim) []string {
	if scan.TestsFiles == 0 {
		return []string{"claim:test-posture", "claim:readiness-verdict"}
	}
	if len(scan.DocReferenceDrifts) > 0 {
		return []string{"claim:documentation-drift", "claim:readiness-verdict"}
	}
	if len(claims) == 0 {
		return []string{}
	}
	return []string{claims[0].ClaimID, "claim:readiness-verdict"}
}

func claimIDList(claims []claim) []string {
	ids := make([]string, 0, len(claims))
	for _, c := range claims {
		ids = append(ids, c.ClaimID)
	}
	return ids
}

func qualityFindings(scan scanResult, claimIDs []string, opts Options) []map[string]any {
	findings := make([]map[string]any, 0)
	if scan.TestsFiles == 0 {
		findings = append(findings, map[string]any{
			"finding_id":  "finding:testing:missing",
			"title":       "Verification surface is thin",
			"severity":    "high",
			"description": "No test files were detected in the repository baseline.",
			"claim_ids":   []string{"claim:test-posture", "claim:readiness-verdict"},
		})
	}
	if len(scan.DocReferenceDrifts) > 0 {
		finding := map[string]any{
			"finding_id":  "finding:docs:drift",
			"title":       "Documented paths drift from repository reality",
			"severity":    "medium",
			"description": fmt.Sprintf("%d documented path references do not resolve locally.", len(scan.DocReferenceDrifts)),
			"claim_ids":   []string{"claim:documentation-drift"},
		}
		findings = append(findings, finding)
	}
	if len(scan.Integrations) > 0 {
		findings = append(findings, map[string]any{
			"finding_id":  "finding:integration:review",
			"title":       "External boundaries need explicit review",
			"severity":    "medium",
			"description": fmt.Sprintf("%d integration or data-store surfaces were inferred from local evidence.", len(scan.Integrations)),
			"claim_ids":   []string{"claim:integration-surface", "claim:readiness-verdict"},
		})
	}
	if len(scan.HotspotFiles) == 0 {
		findings = append(findings, map[string]any{
			"finding_id":  "finding:hotspot:none",
			"title":       "No large source hotspots detected",
			"severity":    "low",
			"description": "No source files exceeded hotspot threshold.",
			"claim_ids":   []string{"claim:source-footprint"},
		})
	} else {
		limit := len(scan.HotspotFiles)
		if opts.Mode == ModeQuick && limit > 3 {
			limit = 3
		}
		for i, h := range scan.HotspotFiles[:limit] {
			findings = append(findings, map[string]any{
				"finding_id":  fmt.Sprintf("finding:hotspot:%d", i+1),
				"title":       fmt.Sprintf("Large file hotspot: %s", h.Path),
				"severity":    "medium",
				"description": fmt.Sprintf("%s has %d lines.", h.Path, h.Lines),
				"claim_ids":   claimIDs,
			})
		}
	}
	if opts.Focus == "architecture" && len(scan.Modules) > 0 {
		findings = append(findings, map[string]any{
			"finding_id":  "finding:architecture:modules",
			"title":       "Architecture focus enabled",
			"severity":    "low",
			"description": fmt.Sprintf("Analysis emphasized %d top-level modules and %d entrypoints.", len(scan.Modules), len(scan.Entrypoints)),
			"claim_ids":   []string{"claim:source-footprint"},
		})
	}
	if opts.Focus == "docs" && len(scan.DocReferenceDrifts) == 0 {
		findings = append(findings, map[string]any{
			"finding_id":  "finding:docs:stable",
			"title":       "No obvious documentation path drift detected",
			"severity":    "low",
			"description": "Local docs path references resolved during the OSS cross-check pass.",
			"claim_ids":   []string{"claim:source-footprint"},
		})
	}
	return findings
}

func hotspotRanking(scan scanResult) []map[string]any {
	result := make([]map[string]any, 0, len(scan.HotspotFiles))
	for _, h := range scan.HotspotFiles {
		result = append(result, map[string]any{
			"path":   h.Path,
			"score":  float64(h.Lines),
			"reason": fmt.Sprintf("line_count=%d", h.Lines),
		})
	}
	return result
}

func unresolvedQuestions(scan scanResult, opts Options) []string {
	questions := []string{
		"Which integrations need contract-level documentation or dedicated boundary tests next?",
	}
	if scan.TestsFiles == 0 {
		questions = append(questions, "Where should baseline verification tests be added first?")
	}
	if len(scan.DocReferenceDrifts) > 0 {
		questions = append(questions, "Which documentation references should be corrected or deleted to remove stale paths?")
	}
	if opts.Mode == ModeBootstrapSDP {
		questions = append(questions, "Which narrow first workstream creates the safest agent-executable scope?")
	}
	return questions
}

func renderSummaryMD(scan scanResult, verdict string, topFindingIDs []string, opts Options) string {
	var b strings.Builder
	b.WriteString("# Reality Summary\n\n")
	b.WriteString(fmt.Sprintf("- Repository: `%s`\n", scan.RepoName))
	b.WriteString(fmt.Sprintf("- Repo Type: `%s`\n", scan.RepoType))
	b.WriteString(fmt.Sprintf("- Generated At: `%s`\n", scan.GeneratedAt))
	b.WriteString(fmt.Sprintf("- Analysis Mode: `%s`\n", opts.Mode))
	if opts.Focus != "" {
		b.WriteString(fmt.Sprintf("- Analysis Focus: `%s`\n", opts.Focus))
	}
	b.WriteString(fmt.Sprintf("- Readiness Verdict: `%s`\n", verdict))
	b.WriteString(fmt.Sprintf("- Source Files: `%d`\n", scan.SourceFiles))
	b.WriteString(fmt.Sprintf("- Test Files: `%d`\n", scan.TestsFiles))
	b.WriteString(fmt.Sprintf("- Config + Manifest Files: `%d`\n", scan.ConfigFiles+scan.ManifestFiles))
	b.WriteString(fmt.Sprintf("- Modules: `%d`\n", len(scan.Modules)))
	b.WriteString(fmt.Sprintf("- Integrations: `%d`\n", len(scan.Integrations)))
	b.WriteString(fmt.Sprintf("- Doc Drift References: `%d`\n", len(scan.DocReferenceDrifts)))
	b.WriteString("\n## Top Finding Claims\n\n")
	for _, id := range topFindingIDs {
		b.WriteString(fmt.Sprintf("- `%s`\n", id))
	}
	return b.String()
}

func renderArchitectureMD(scan scanResult, opts Options) string {
	var b strings.Builder
	b.WriteString("# Reality Architecture\n\n")
	b.WriteString(fmt.Sprintf("- Analysis Mode: `%s`\n", opts.Mode))
	if opts.Focus != "" {
		b.WriteString(fmt.Sprintf("- Analysis Focus: `%s`\n", opts.Focus))
	}
	b.WriteString("## Modules\n\n")
	if len(scan.Modules) == 0 {
		b.WriteString("- none\n")
	} else {
		for _, m := range scan.Modules {
			b.WriteString(fmt.Sprintf("- `%s`\n", m))
		}
	}
	b.WriteString("\n## Entrypoints\n\n")
	if len(scan.Entrypoints) == 0 {
		b.WriteString("- none detected\n")
	} else {
		for _, ep := range scan.Entrypoints {
			b.WriteString(fmt.Sprintf("- `%s`\n", ep))
		}
	}
	b.WriteString("\n## Integrations\n\n")
	if len(scan.Integrations) == 0 {
		b.WriteString("- none detected\n")
	} else {
		for _, integration := range scan.Integrations {
			b.WriteString(fmt.Sprintf("- `%s` (%s via %s)\n", integration.Name, integration.Kind, strings.Join(integration.EvidencePaths, ", ")))
		}
	}
	return b.String()
}

func renderQualityMD(scan scanResult, opts Options) string {
	var b strings.Builder
	b.WriteString("# Reality Quality\n\n")
	b.WriteString(fmt.Sprintf("- Analysis Mode: `%s`\n", opts.Mode))
	if opts.Focus != "" {
		b.WriteString(fmt.Sprintf("- Analysis Focus: `%s`\n", opts.Focus))
	}
	b.WriteString(fmt.Sprintf("- Source Files: `%d`\n", scan.SourceFiles))
	b.WriteString(fmt.Sprintf("- Test Files: `%d`\n", scan.TestsFiles))
	b.WriteString(fmt.Sprintf("- Docs Files: `%d`\n", scan.DocFiles))
	b.WriteString(fmt.Sprintf("- Config Files: `%d`\n", scan.ConfigFiles))
	b.WriteString(fmt.Sprintf("- Manifest Files: `%d`\n", scan.ManifestFiles))
	b.WriteString("\n## Hotspots\n\n")
	if len(scan.HotspotFiles) == 0 {
		b.WriteString("- none\n")
	} else {
		for _, h := range scan.HotspotFiles {
			b.WriteString(fmt.Sprintf("- `%s` (%d lines)\n", h.Path, h.Lines))
		}
	}
	b.WriteString("\n## Documentation Drift\n\n")
	if len(scan.DocReferenceDrifts) == 0 {
		b.WriteString("- none detected\n")
	} else {
		for _, drift := range scan.DocReferenceDrifts {
			b.WriteString(fmt.Sprintf("- `%s` references missing `%s`\n", drift.DocPath, drift.ReferencedPath))
		}
	}
	return b.String()
}

func renderBootstrapMD(scan scanResult, verdict string, constraints []string, recommendations []string, opts Options) string {
	var b strings.Builder
	b.WriteString("# Reality Bootstrap\n\n")
	b.WriteString(fmt.Sprintf("- Current Verdict: `%s`\n", verdict))
	b.WriteString(fmt.Sprintf("- Analysis Mode: `%s`\n", opts.Mode))
	if opts.Focus != "" {
		b.WriteString(fmt.Sprintf("- Analysis Focus: `%s`\n", opts.Focus))
	}
	b.WriteString("\n## Constraints\n\n")
	if len(constraints) == 0 {
		b.WriteString("- none\n")
	} else {
		for _, c := range constraints {
			b.WriteString(fmt.Sprintf("- %s\n", c))
		}
	}
	b.WriteString("\n## Suggested First Workstreams\n\n")
	if len(recommendations) == 0 {
		b.WriteString("- No immediate bootstrap workstreams inferred.\n")
	} else {
		for _, recommendation := range recommendations {
			b.WriteString(fmt.Sprintf("- %s\n", recommendation))
		}
	}
	if opts.Mode == ModeBootstrapSDP {
		b.WriteString("\n## Agent Readiness Notes\n\n")
		b.WriteString(fmt.Sprintf("- Safe modules for first slices: `%s`\n", strings.Join(preferredBootstrapModules(scan), "`, `")))
	}
	return b.String()
}

func writeJSON(path string, payload any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func writeText(path, body string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	body = strings.TrimRight(body, "\n") + "\n"
	return os.WriteFile(path, []byte(body), 0o644)
}

func featureStatus(scan scanResult) string {
	if scan.SourceFiles == 0 {
		return "candidate"
	}
	if scan.TestsFiles == 0 || len(scan.DocReferenceDrifts) > 0 {
		return "partial"
	}
	return "implemented"
}

func confidenceFromScan(scan scanResult) float64 {
	if scan.SourceFiles == 0 {
		return 0.4
	}
	if scan.TestsFiles == 0 || len(scan.DocReferenceDrifts) > 0 {
		return 0.7
	}
	return 0.9
}

func inferredEdges(modules []string) []map[string]any {
	moduleSet := map[string]bool{}
	for _, m := range modules {
		moduleSet[m] = true
	}
	edges := make([]map[string]any, 0)
	if moduleSet["cmd"] && moduleSet["internal"] {
		edges = append(edges, map[string]any{
			"from":       "module:cmd",
			"to":         "module:internal",
			"relation":   "depends_on",
			"confidence": 0.6,
		})
	}
	return edges
}

func inferredIntegrationEdges(scan scanResult) []map[string]any {
	edges := make([]map[string]any, 0, len(scan.Integrations))
	for _, integration := range scan.Integrations {
		module := "root"
		if len(integration.EvidencePaths) > 0 {
			module = topModule(integration.EvidencePaths[0])
		}
		edges = append(edges, map[string]any{
			"from":       "module:" + module,
			"to":         "integration:" + integration.Name,
			"relation":   "depends_on",
			"confidence": integration.Confidence,
		})
	}
	return edges
}

func runID(repoName string) string {
	return fmt.Sprintf("reality-oss-%s", repoName)
}

func modulePath(module string) string {
	if module == "root" {
		return "."
	}
	return module
}

func stringSliceOrEmpty(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}

func topModule(rel string) string {
	parts := strings.Split(rel, "/")
	if len(parts) <= 1 {
		return "root"
	}
	return parts[0]
}

func isTestFile(rel string) bool {
	return strings.HasSuffix(rel, "_test.go") ||
		strings.HasSuffix(rel, ".test.js") ||
		strings.HasSuffix(rel, ".test.ts") ||
		strings.HasSuffix(rel, "_spec.py")
}

func countLines(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	lines := 0
	for scanner.Scan() {
		lines++
	}
	return lines
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func normalizeOptions(opts Options) Options {
	if opts.Mode == "" {
		opts.Mode = ModeDeep
	}
	opts.Focus = strings.TrimSpace(strings.ToLower(opts.Focus))
	return opts
}

func validateOptions(opts Options) error {
	switch opts.Mode {
	case ModeQuick, ModeDeep, ModeBootstrapSDP:
	default:
		return fmt.Errorf("unsupported reality mode %q", opts.Mode)
	}
	if !validFocuses[opts.Focus] {
		return fmt.Errorf("unsupported focus %q", opts.Focus)
	}
	return nil
}

func summaryScope(scan scanResult, opts Options) map[string]any {
	scope := map[string]any{
		"repos": []string{scan.RepoName},
		"mode":  string(opts.Mode),
	}
	if opts.Focus != "" {
		scope["focus"] = opts.Focus
	}
	return scope
}

func featureSummary(scan scanResult, opts Options) string {
	summary := "Repository baseline reconstructed from local source tree, tests, configs, manifests, and docs."
	if opts.Mode == ModeQuick {
		return summary + " Quick mode keeps the same artifact families with shallower evidence detail."
	}
	if opts.Mode == ModeBootstrapSDP {
		return summary + " Bootstrap mode prioritizes first SDP-safe workstream recommendations."
	}
	return summary
}

func classifyConfigOrManifest(rel string) (string, bool) {
	lower := strings.ToLower(rel)
	base := strings.ToLower(filepath.Base(rel))
	ext := strings.ToLower(filepath.Ext(rel))
	if base == "dockerfile" || strings.Contains(lower, "docker-compose") || strings.Contains(lower, "compose.yaml") {
		return "manifest", true
	}
	if strings.Contains(lower, "deploy/") || strings.Contains(lower, "k8s/") || strings.Contains(lower, "helm/") {
		return "manifest", true
	}
	switch ext {
	case ".yaml", ".yml":
		if strings.Contains(lower, "deploy") || strings.Contains(lower, "k8s") || strings.Contains(lower, "helm") {
			return "manifest", true
		}
		return "config", true
	case ".json", ".toml", ".ini", ".env":
		return "config", true
	}
	if strings.HasSuffix(base, ".conf") {
		return "config", true
	}
	return "", false
}

func shouldScanForIntegrations(rel, ext string) bool {
	if strings.HasPrefix(rel, "docs/reality/") {
		return false
	}
	if _, ok := classifyConfigOrManifest(rel); ok {
		return true
	}
	if ext == ".md" && strings.HasPrefix(rel, "docs/") {
		return true
	}
	switch ext {
	case ".go", ".py", ".js", ".ts", ".java", ".rs", ".sh":
		return true
	default:
		return false
	}
}

func readFileForScan(path string) string {
	info, err := os.Stat(path)
	if err != nil || info.Size() > 1<<20 {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.ToLower(string(data))
}

func detectIntegrations(rel, content string) []integrationObservation {
	type pattern struct {
		Name         string
		Kind         string
		ContractType string
		Keywords     []string
	}

	patterns := []pattern{
		{Name: "postgres", Kind: "data_store", ContractType: "sql", Keywords: []string{"postgres", "postgresql", "pgx"}},
		{Name: "mysql", Kind: "data_store", ContractType: "sql", Keywords: []string{"mysql"}},
		{Name: "sqlite", Kind: "data_store", ContractType: "sql", Keywords: []string{"sqlite"}},
		{Name: "redis", Kind: "data_store", ContractType: "cache", Keywords: []string{"redis"}},
		{Name: "kafka", Kind: "integration", ContractType: "event_stream", Keywords: []string{"kafka"}},
		{Name: "nats", Kind: "integration", ContractType: "event_stream", Keywords: []string{"nats"}},
		{Name: "rabbitmq", Kind: "integration", ContractType: "message_queue", Keywords: []string{"rabbitmq"}},
		{Name: "grpc", Kind: "integration", ContractType: "rpc", Keywords: []string{"grpc"}},
		{Name: "http", Kind: "integration", ContractType: "http", Keywords: []string{"http://", "https://", "httpclient", "resty"}},
		{Name: "s3", Kind: "integration", ContractType: "object_store", Keywords: []string{"s3", "minio"}},
		{Name: "github", Kind: "integration", ContractType: "api", Keywords: []string{"github", "gh "}},
		{Name: "slack", Kind: "integration", ContractType: "api", Keywords: []string{"slack"}},
		{Name: "stripe", Kind: "integration", ContractType: "api", Keywords: []string{"stripe"}},
	}

	haystack := strings.ToLower(rel) + "\n" + content
	detected := make([]integrationObservation, 0)
	for _, pattern := range patterns {
		for _, keyword := range pattern.Keywords {
			if strings.Contains(haystack, keyword) {
				confidence := 0.65
				if kind, ok := classifyConfigOrManifest(rel); ok {
					if kind == "manifest" {
						confidence = 0.85
					} else {
						confidence = 0.8
					}
				}
				if strings.HasSuffix(rel, ".md") {
					confidence = 0.55
				}
				detected = append(detected, integrationObservation{
					Name:          pattern.Name,
					Kind:          pattern.Kind,
					ContractType:  pattern.ContractType,
					EvidencePaths: []string{rel},
					Confidence:    confidence,
				})
				break
			}
		}
	}
	return detected
}

func mergeIntegration(target map[string]integrationObservation, detected integrationObservation) {
	current, ok := target[detected.Name]
	if !ok {
		target[detected.Name] = detected
		return
	}
	current.EvidencePaths = append(current.EvidencePaths, detected.EvidencePaths...)
	current.EvidencePaths = dedupeStrings(current.EvidencePaths)
	if detected.Confidence > current.Confidence {
		current.Confidence = detected.Confidence
		current.Kind = detected.Kind
		current.ContractType = detected.ContractType
	}
	target[detected.Name] = current
}

func flattenIntegrations(m map[string]integrationObservation) []integrationObservation {
	items := make([]integrationObservation, 0, len(m))
	for _, item := range m {
		sort.Strings(item.EvidencePaths)
		items = append(items, item)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items
}

func integrationEntries(scan scanResult) []map[string]any {
	entries := make([]map[string]any, 0, len(scan.Integrations))
	for _, integration := range scan.Integrations {
		entries = append(entries, map[string]any{
			"integration_id":   fmt.Sprintf("integration:%s", integration.Name),
			"name":             integration.Name,
			"integration_type": integration.Kind,
			"producer":         "repo:" + scan.RepoName,
			"consumer":         "external:" + integration.Name,
			"contract_type":    integration.ContractType,
			"confidence":       integration.Confidence,
			"evidence_paths":   integration.EvidencePaths,
			"risk_notes": []string{
				"Confirm runtime ownership and failure behavior for this boundary.",
			},
		})
	}
	return entries
}

func scanDocReferenceDrift(root, rel, absPath string) []docReferenceDrift {
	content := readFileForScan(absPath)
	if content == "" {
		return nil
	}
	pathPattern := regexp.MustCompile("`([A-Za-z0-9_./-]+)`")
	matches := pathPattern.FindAllStringSubmatch(content, -1)
	drifts := make([]docReferenceDrift, 0)
	for _, match := range matches {
		referenced := normalizeDocPathReference(match[1])
		if !looksLikeTrackedPath(referenced) {
			continue
		}
		if _, err := os.Stat(filepath.Join(root, referenced)); err == nil {
			continue
		}
		drifts = append(drifts, docReferenceDrift{
			DocPath:         rel,
			ReferencedPath:  referenced,
			ObservationNote: "Documented path does not resolve in the current repository tree.",
		})
	}
	return drifts
}

func normalizeDocPathReference(ref string) string {
	ref = strings.TrimSpace(ref)
	ref = strings.TrimPrefix(ref, "./")
	ref = strings.TrimSuffix(ref, "/")
	return filepath.ToSlash(ref)
}

func looksLikeTrackedPath(ref string) bool {
	if ref == "" {
		return false
	}
	prefixes := []string{"cmd/", "internal/", "pkg/", "api/", "deploy/", "configs/", "schema/", "scripts/"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(ref, prefix) {
			return true
		}
	}
	suffixes := []string{".go", ".py", ".js", ".ts", ".java", ".rs", ".sh", ".md", ".yaml", ".yml", ".json", ".toml"}
	for _, suffix := range suffixes {
		if strings.HasSuffix(ref, suffix) {
			return true
		}
	}
	return false
}

func driftFindings(scan scanResult) []map[string]any {
	findings := make([]map[string]any, 0, len(scan.DocReferenceDrifts))
	for i, drift := range scan.DocReferenceDrifts {
		findings = append(findings, map[string]any{
			"contradiction_id": fmt.Sprintf("drift:%d", i+1),
			"title":            fmt.Sprintf("Missing documented path: %s", drift.ReferencedPath),
			"severity":         "medium",
			"doc_path":         drift.DocPath,
			"referenced_path":  drift.ReferencedPath,
			"note":             drift.ObservationNote,
		})
	}
	return findings
}

func bootstrapRecommendations(scan scanResult, verdict string, opts Options) []string {
	recommendations := make([]string, 0)
	if scan.TestsFiles == 0 {
		recommendations = append(recommendations, "Add a first verification slice around a narrow entrypoint or core module.")
	}
	if len(scan.DocReferenceDrifts) > 0 {
		recommendations = append(recommendations, "Reconcile stale documentation paths before delegating agent work from docs.")
	}
	if len(scan.Integrations) > 0 {
		recommendations = append(recommendations, "Fence external integration boundaries with explicit contracts and failure notes.")
	}
	if len(scan.HotspotFiles) > 0 {
		recommendations = append(recommendations, "Split hotspot files into smaller scopes before broad autonomous changes.")
	}
	if verdict == "ready" {
		recommendations = append(recommendations, "Start with a small SDP workstream in the safest tested module.")
	}
	if opts.Mode == ModeBootstrapSDP && len(recommendations) == 0 {
		recommendations = append(recommendations, "Seed the first SDP workstream from the lowest-coupling module and keep scope single-boundary.")
	}
	return dedupeStrings(recommendations)
}

func preferredBootstrapModules(scan scanResult) []string {
	if len(scan.Modules) == 0 {
		return []string{"root"}
	}
	limit := len(scan.Modules)
	if limit > 3 {
		limit = 3
	}
	return scan.Modules[:limit]
}

func classifyRepoType(scan scanResult) string {
	hasCmd := false
	hasSchema := false
	hasDeploy := false
	for _, module := range scan.Modules {
		switch module {
		case "cmd":
			hasCmd = true
		case "schema":
			hasSchema = true
		case "deploy":
			hasDeploy = true
		}
	}
	switch {
	case hasCmd && hasDeploy:
		return "app"
	case hasSchema && !hasCmd:
		return "protocol"
	case hasDeploy && !hasCmd:
		return "infra"
	case hasCmd:
		return "service"
	default:
		return "mixed"
	}
}

func dedupeStrings(values []string) []string {
	seen := map[string]bool{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		result = append(result, value)
	}
	return result
}
