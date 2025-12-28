package api

import (
	"errors"
	"time"

	"github.com/Rx-11/distributed-leaderboard/cache"
	"github.com/Rx-11/distributed-leaderboard/global"
	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

type SelectedSummaries struct {
	Summaries []leaderboard.RegionSummary
	Included  map[leaderboard.RegionID]bool
}

func selectSummaries(now time.Time, mode global.ConsistencyMode, cache *cache.SummaryCache, local leaderboard.RegionSummary, totalRegions int) (SelectedSummaries, Coverage, Staleness, error) {
	summaries := SelectedSummaries{
		Summaries: []leaderboard.RegionSummary{local},
		Included:  map[leaderboard.RegionID]bool{local.Region: true},
	}
	stalest := time.Duration(0)

	for region, entry := range cache.Entries() {
		age := now.Sub(entry.ReceivedAt)
		summaries.Included[region] = false

		if mode == global.Strict && age > cache.FreshTTL() {
			return SelectedSummaries{}, Coverage{}, Staleness{}, errors.New("strict mode: stale region")
		}

		if age <= cache.StaleTTL() {
			summaries.Summaries = append(summaries.Summaries, entry.Summary)
			summaries.Included[region] = true
			if age > stalest {
				stalest = age
			}
		}
	}

	included := len(summaries.Summaries)

	coverage := Coverage{
		TotalRegions:    totalRegions,
		IncludedRegions: included,
		CoverageRatio:   float64(included) / float64(totalRegions),
	}

	if mode == global.Strict {
		return summaries, coverage, Staleness{
			MaxAge: cache.FreshTTL(),
			Oldest: stalest,
		}, nil
	}

	return summaries, coverage, Staleness{
		MaxAge: cache.StaleTTL(),
		Oldest: stalest,
	}, nil

}

func GetGlobalTopK(now time.Time, mode global.ConsistencyMode, local *leaderboard.Leaderboard, cache *cache.SummaryCache, totalRegions int, k int) (*GlobalTopKResponse, error) {

	localSummary := local.RegionSummary(k)

	summaries, coverage, staleness, err := selectSummaries(now, mode, cache, localSummary, totalRegions)
	if err != nil {
		return nil, err
	}

	result, err := global.ComputeGlobalTopK(summaries.Summaries, k, mode)
	if err != nil {
		return nil, err
	}

	return &GlobalTopKResponse{
		Entries:   result.Entries,
		Coverage:  coverage,
		Staleness: staleness,
	}, nil
}

func GetGlobalRank(now time.Time, mode global.ConsistencyMode, user leaderboard.Entry, localRank int, localRegion leaderboard.RegionID, local *leaderboard.Leaderboard, cache *cache.SummaryCache, totalRegions int) (*GlobalRankResponse, error) {
	localSummary := local.RegionSummary(0)

	summaries, coverage, staleness, err := selectSummaries(now, mode, cache, localSummary, totalRegions)
	if err != nil {
		return nil, err
	}

	est := global.EstimateGlobalRank(
		user,
		localRank,
		localRegion,
		summaries.Summaries,
	)

	resp := &GlobalRankResponse{
		LowerBound: est.Rank,
		Coverage:   coverage,
		Staleness:  staleness,
	}

	if mode == global.Fast {
		upper := est.Rank
		for _, entry := range cache.AllSummaries() {
			if !summaries.Included[entry.Summary.Region] {
				upper += entry.Summary.UserCount
			}
		}

		resp.UpperBound = &upper
	}

	return resp, nil
}
