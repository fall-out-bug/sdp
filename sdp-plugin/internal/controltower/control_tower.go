package controltower

import (
	"fmt"
	"os"

	"github.com/fall-out-bug/sdp/internal/config"
	"github.com/fall-out-bug/sdp/internal/nextstep"
)

type Data struct {
	ProjectRoot string
	State       nextstep.ProjectState
	NextStep    *nextstep.Recommendation
	StatusView  *nextstep.StatusView
}

func Collect() (*Data, error) {
	projectRoot, err := config.FindProjectRoot()
	if err != nil {
		projectRoot, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	collector := nextstep.NewStateCollector(projectRoot)
	state, err := collector.Collect()
	if err != nil {
		return nil, fmt.Errorf("failed to collect project state: %w", err)
	}

	resolver := nextstep.NewResolver()
	rec, err := resolver.Recommend(state)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendation: %w", err)
	}

	return &Data{
		ProjectRoot: projectRoot,
		State:       state,
		NextStep:    rec,
		StatusView:  nextstep.BuildStatusView(projectRoot, state, rec),
	}, nil
}
