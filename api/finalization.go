package api

import (
	"errors"
	"fmt"
	"time"

	"github.com/Rx-11/distributed-leaderboard/cache"
	"github.com/Rx-11/distributed-leaderboard/global"
	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

var (
	ErrSeasonNotReady = errors.New("season not ready for finalization")
	ErrRegionMissing  = errors.New("missing region data")
	ErrRegionNotFinal = errors.New("region has not finalized season")
)

func FinalizeSeason(now time.Time, seasonID leaderboard.SeasonID, expectedRegions []leaderboard.RegionID, cache *cache.SummaryCache, k int) (*leaderboard.GlobalSeasonSnapshot, error) {

	summaries := make([]leaderboard.RegionSummary, 0, len(expectedRegions))
	totalUsers := 0

	for _, regionID := range expectedRegions {
		entry, exists := cache.Get(regionID)
		if !exists {
			return nil, fmt.Errorf("%w: region %s", ErrRegionMissing, regionID)
		}

		if entry.Summary.Season != seasonID {
			return nil, fmt.Errorf("region %s has wrong season data (expected %s, got %s)", regionID, seasonID, entry.Summary.Season)
		}

		if !entry.Summary.IsFinal {
			return nil, fmt.Errorf("%w: region %s is still live", ErrRegionNotFinal, regionID)
		}

		summaries = append(summaries, entry.Summary)
		totalUsers += entry.Summary.UserCount
	}

	result, err := global.ComputeGlobalTopK(summaries, k, global.Strict)
	if err != nil {
		return nil, fmt.Errorf("merge failed: %w", err)
	}

	snapshot := &leaderboard.GlobalSeasonSnapshot{
		SeasonID:        seasonID,
		FinalizedAt:     now,
		TopK:            result.Entries,
		IncludedRegions: expectedRegions,
		TotalUserCount:  totalUsers,
	}

	return snapshot, nil
}
