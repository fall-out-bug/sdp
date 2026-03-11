package reality

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

const (
	specVersion            = "v1.0"
	fallbackGeneratedAtUTC = "1970-01-01T00:00:00Z"
	hotspotLineThreshold   = 800
)

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
	Root         string
	RepoName     string
	GeneratedAt  string
	TotalFiles   int
	SourceFiles  int
	TestFiles    int
	DocFiles     int
	Modules      []string
	Entrypoints  []string
	HotspotFiles []hotspot
}

type hotspot struct {
	Path  string
	Lines int
}

// EmitOSS writes the OSS reality artifact set to docs/reality and .sdp/reality.
func EmitOSS(projectRoot string) ([]string, error) {
	scan, err := scanProject(projectRoot)
	if err != nil {
		return nil, err
	}

	verdict := readinessVerdict(scan)
	dim := readinessDimensions(scan)
	constraints := readinessConstraints(scan)
	claims := buildClaims(scan, verdict)
	sources := buildSources(scan)
	claimIDs := claimIDList(claims)
	topFindingIDs := topFindingClaimIDs(scan, claims)

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
		"scope":                 map[string]any{"repos": []string{scan.RepoName}, "mode": "emit-oss"},
		"readiness_verdict":     verdict,
		"top_finding_claim_ids": topFindingIDs,
		"artifacts":             artifactPaths,
		"claims":                claims,
		"sources":               sources,
	}

	featureInventory := map[string]any{
		"spec_version": specVersion,
		"generated_at": scan.GeneratedAt,
		"features": []map[string]any{
			{
				"feature_id":         "feature:repository-baseline",
				"title":              fmt.Sprintf("%s baseline", scan.RepoName),
				"summary":            "Repository baseline reconstructed from local source tree.",
				"status":             featureStatus(scan),
				"evidence_claim_ids": topFindingIDs,
				"confidence":         confidenceFromScan(scan),
				"mapped_components":  scan.Modules,
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
	archEdges := inferredEdges(scan.Modules)
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
		"integrations": []map[string]any{},
		"claims":       claims,
		"sources":      sources,
	}

	qualityFindings := qualityFindings(scan, claimIDs)
	qualityReport := map[string]any{
		"spec_version":    specVersion,
		"generated_at":    scan.GeneratedAt,
		"findings":        qualityFindings,
		"hotspot_ranking": hotspotRanking(scan),
		"claims":          claims,
		"sources":         sources,
	}

	driftReport := map[string]any{
		"spec_version":         specVersion,
		"generated_at":         scan.GeneratedAt,
		"contradictions":       []map[string]any{},
		"unresolved_questions": unresolvedQuestions(scan),
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
		"docs/reality/summary.md":      renderSummaryMD(scan, verdict, topFindingIDs),
		"docs/reality/architecture.md": renderArchitectureMD(scan),
		"docs/reality/quality.md":      renderQualityMD(scan),
		"docs/reality/bootstrap.md":    renderBootstrapMD(verdict, constraints),
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
		}
		if sourceExt[ext] {
			scan.SourceFiles++
			module := topModule(rel)
			modules[module] = true
			if strings.HasSuffix(rel, "/main.go") || rel == "main.go" {
				entrypoints[rel] = true
			}
			if isTestFile(rel) {
				scan.TestFiles++
			}
			lines := countLines(path)
			if lines >= hotspotLineThreshold {
				hotspots = append(hotspots, hotspot{Path: rel, Lines: lines})
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
	if scan.TestFiles == 0 {
		return "ready_with_constraints"
	}
	return "ready"
}

func readinessConstraints(scan scanResult) []string {
	constraints := make([]string, 0)
	if scan.SourceFiles == 0 {
		constraints = append(constraints, "No source files detected for analysis.")
	}
	if scan.TestFiles == 0 {
		constraints = append(constraints, "No test files detected; verification surface is weak.")
	}
	if len(scan.HotspotFiles) > 0 {
		constraints = append(constraints, "Large hotspot files detected; prefer narrow scoped changes.")
	}
	return constraints
}

func readinessDimensions(scan scanResult) map[string]any {
	moduleCount := float64(len(scan.Modules))
	sourceCount := float64(max(scan.SourceFiles, 1))
	hotspotCount := float64(len(scan.HotspotFiles))
	docCount := float64(scan.DocFiles)

	return map[string]any{
		"boundary_clarity": map[string]any{
			"score": clamp(moduleCount / 8),
			"note":  fmt.Sprintf("%d top-level modules", len(scan.Modules)),
		},
		"verification_coverage": map[string]any{
			"score": clamp(float64(scan.TestFiles) / sourceCount),
			"note":  fmt.Sprintf("%d test files / %d source files", scan.TestFiles, scan.SourceFiles),
		},
		"hotspot_concentration": map[string]any{
			"score": clamp(1 - (hotspotCount / sourceCount)),
			"note":  fmt.Sprintf("%d hotspot files (>= %d lines)", len(scan.HotspotFiles), hotspotLineThreshold),
		},
		"integration_fragility": map[string]any{
			"score": 0.5,
			"note":  "No explicit integration extraction in OSS baseline emitter yet.",
		},
		"documentation_trust_level": map[string]any{
			"score": clamp(docCount / float64(max(scan.DocFiles+len(scan.Modules), 1))),
			"note":  fmt.Sprintf("%d markdown docs under docs/", scan.DocFiles),
		},
	}
}

func buildClaims(scan scanResult, verdict string) []claim {
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
			Statement:   fmt.Sprintf("Detected %d test files.", scan.TestFiles),
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
	return claims
}

func buildSources(scan scanResult) []source {
	return []source{
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
}

func topFindingClaimIDs(scan scanResult, claims []claim) []string {
	if scan.TestFiles == 0 {
		return []string{"claim:test-posture", "claim:readiness-verdict"}
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

func qualityFindings(scan scanResult, claimIDs []string) []map[string]any {
	findings := make([]map[string]any, 0)
	if len(scan.HotspotFiles) == 0 {
		findings = append(findings, map[string]any{
			"finding_id":  "finding:hotspot:none",
			"title":       "No large source hotspots detected",
			"severity":    "low",
			"description": "No source files exceeded hotspot threshold.",
			"claim_ids":   []string{"claim:source-footprint"},
		})
		return findings
	}
	for i, h := range scan.HotspotFiles {
		if i >= 5 {
			break
		}
		findings = append(findings, map[string]any{
			"finding_id":  fmt.Sprintf("finding:hotspot:%d", i+1),
			"title":       fmt.Sprintf("Large file hotspot: %s", h.Path),
			"severity":    "medium",
			"description": fmt.Sprintf("%s has %d lines.", h.Path, h.Lines),
			"claim_ids":   claimIDs,
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

func unresolvedQuestions(scan scanResult) []string {
	questions := []string{
		"Which integrations should be promoted from inferred to observed in the OSS emitter?",
	}
	if scan.TestFiles == 0 {
		questions = append(questions, "Where should baseline verification tests be added first?")
	}
	return questions
}

func renderSummaryMD(scan scanResult, verdict string, topFindingIDs []string) string {
	var b strings.Builder
	b.WriteString("# Reality Summary\n\n")
	b.WriteString(fmt.Sprintf("- Repository: `%s`\n", scan.RepoName))
	b.WriteString(fmt.Sprintf("- Generated At: `%s`\n", scan.GeneratedAt))
	b.WriteString(fmt.Sprintf("- Readiness Verdict: `%s`\n", verdict))
	b.WriteString(fmt.Sprintf("- Source Files: `%d`\n", scan.SourceFiles))
	b.WriteString(fmt.Sprintf("- Test Files: `%d`\n", scan.TestFiles))
	b.WriteString(fmt.Sprintf("- Modules: `%d`\n", len(scan.Modules)))
	b.WriteString("\n## Top Finding Claims\n\n")
	for _, id := range topFindingIDs {
		b.WriteString(fmt.Sprintf("- `%s`\n", id))
	}
	return b.String()
}

func renderArchitectureMD(scan scanResult) string {
	var b strings.Builder
	b.WriteString("# Reality Architecture\n\n")
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
	return b.String()
}

func renderQualityMD(scan scanResult) string {
	var b strings.Builder
	b.WriteString("# Reality Quality\n\n")
	b.WriteString(fmt.Sprintf("- Source Files: `%d`\n", scan.SourceFiles))
	b.WriteString(fmt.Sprintf("- Test Files: `%d`\n", scan.TestFiles))
	b.WriteString(fmt.Sprintf("- Docs Files: `%d`\n", scan.DocFiles))
	b.WriteString("\n## Hotspots\n\n")
	if len(scan.HotspotFiles) == 0 {
		b.WriteString("- none\n")
	} else {
		for _, h := range scan.HotspotFiles {
			b.WriteString(fmt.Sprintf("- `%s` (%d lines)\n", h.Path, h.Lines))
		}
	}
	return b.String()
}

func renderBootstrapMD(verdict string, constraints []string) string {
	var b strings.Builder
	b.WriteString("# Reality Bootstrap\n\n")
	b.WriteString(fmt.Sprintf("- Current Verdict: `%s`\n", verdict))
	b.WriteString("\n## Constraints\n\n")
	if len(constraints) == 0 {
		b.WriteString("- none\n")
	} else {
		for _, c := range constraints {
			b.WriteString(fmt.Sprintf("- %s\n", c))
		}
	}
	b.WriteString("\n## Suggested First Workstreams\n\n")
	b.WriteString("- Materialize artifact schemas (`schema/reality/*.schema.json`)\n")
	b.WriteString("- Implement deterministic artifact emitter (`sdp reality emit-oss`)\n")
	b.WriteString("- Add schema validation gate for `.sdp/reality/*.json`\n")
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
	if scan.TestFiles == 0 {
		return "partial"
	}
	return "implemented"
}

func confidenceFromScan(scan scanResult) float64 {
	if scan.SourceFiles == 0 {
		return 0.4
	}
	if scan.TestFiles == 0 {
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

func runID(repoName string) string {
	return fmt.Sprintf("reality-oss-%s", repoName)
}

func modulePath(module string) string {
	if module == "root" {
		return "."
	}
	return module
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
