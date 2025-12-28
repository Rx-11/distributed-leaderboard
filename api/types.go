package api

import (
	"time"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

type Coverage struct {
	TotalRegions    int
	IncludedRegions int
	CoverageRatio   float64
}

type Staleness struct {
	MaxAge time.Duration
	Oldest time.Duration
}

type GlobalTopKResponse struct {
	Entries   []leaderboard.Entry
	Coverage  Coverage
	Staleness Staleness
}

type GlobalRankResponse struct {
	LowerBound int
	UpperBound *int
	Coverage   Coverage
	Staleness  Staleness
}
