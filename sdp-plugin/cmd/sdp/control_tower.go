package main

import (
	"github.com/fall-out-bug/sdp/internal/controltower"
)

type controlTowerData = controltower.Data

func collectControlTowerData() (*controlTowerData, error) {
	return controltower.Collect()
}
