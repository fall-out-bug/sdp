package orchestrate

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// WorkstreamInfo holds parsed workstream metadata.
type WorkstreamInfo struct {
	ID         string
	FeatureID  string
	DependsOn  []string
}

// DiscoverWorkstreams finds workstream files for a feature and returns IDs in dependency order.
// Pattern: docs/workstreams/backlog/00-FFF-SS.md for feature FFFF.
func DiscoverWorkstreams(projectRoot, featureID string) ([]string, error) {
	fnum := strings.TrimPrefix(strings.ToUpper(featureID), "F")
	if fnum == "" {
		return nil, fmt.Errorf("invalid feature_id %q", featureID)
	}
	pattern := fmt.Sprintf("00-%s-*.md", strings.TrimLeft(fnum, "0"))
	dir := filepath.Join(projectRoot, "docs", "workstreams", "backlog")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read workstreams dir: %w", err)
	}

	var infos []WorkstreamInfo
	prefix := "00-" + fnum + "-"
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) && strings.HasSuffix(e.Name(), ".md") {
			path := filepath.Join(dir, e.Name())
			info, err := parseWorkstreamFrontmatter(path)
			if err != nil {
				continue
			}
			infos = append(infos, info)
		}
	}
	if len(infos) == 0 {
		return nil, fmt.Errorf("no workstreams found for %s (pattern %s)", featureID, pattern)
	}

	ordered, err := topologicalSort(infos)
	if err != nil {
		return nil, err
	}
	return ordered, nil
}

var (
	reWSID     = regexp.MustCompile(`(?m)^ws_id:\s*(\S+)`)
	reFeature  = regexp.MustCompile(`(?m)^feature_id:\s*(\S+)`)
	reDepends  = regexp.MustCompile(`(?m)^depends_on:\s*\[(.*?)\]`)
)

func parseWorkstreamFrontmatter(path string) (WorkstreamInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return WorkstreamInfo{}, err
	}
	content := string(data)
	info := WorkstreamInfo{}
	if m := reWSID.FindStringSubmatch(content); len(m) > 1 {
		info.ID = strings.Trim(m[1], `"`)
	}
	if m := reFeature.FindStringSubmatch(content); len(m) > 1 {
		info.FeatureID = strings.Trim(m[1], `"`)
	}
	if m := reDepends.FindStringSubmatch(content); len(m) > 1 {
		inner := m[1]
		for _, s := range strings.Split(inner, ",") {
			id := strings.Trim(strings.TrimSpace(s), `"`)
			if id != "" {
				info.DependsOn = append(info.DependsOn, id)
			}
		}
	}
	return info, nil
}

func topologicalSort(infos []WorkstreamInfo) ([]string, error) {
	idToInfo := make(map[string]WorkstreamInfo)
	for _, i := range infos {
		idToInfo[i.ID] = i
	}
	var order []string
	// 0=unvisited, 1=inProgress, 2=completed
	state := make(map[string]int)
	var visit func(id string) error
	visit = func(id string) error {
		switch state[id] {
		case 1:
			return fmt.Errorf("cycle detected in workstream dependencies: %s", id)
		case 2:
			return nil
		}
		state[id] = 1
		info, ok := idToInfo[id]
		if !ok {
			state[id] = 2
			return nil
		}
		for _, dep := range info.DependsOn {
			if _, ok := idToInfo[dep]; ok {
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		state[id] = 2
		order = append(order, id)
		return nil
	}
	sort.Slice(infos, func(i, j int) bool { return infos[i].ID < infos[j].ID })
	for _, info := range infos {
		if err := visit(info.ID); err != nil {
			return nil, err
		}
	}
	return order, nil
}
