package cache

import (
	"time"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

type CachedSummary struct {
	Summary    leaderboard.RegionSummary
	ReceivedAt time.Time
}

type SummaryCache struct {
	self    leaderboard.RegionID
	entries map[leaderboard.RegionID]CachedSummary
	maxAge  time.Duration
}

func NewSummaryCache(self leaderboard.RegionID, maxAge time.Duration) *SummaryCache {

	return &SummaryCache{
		self:    self,
		entries: make(map[leaderboard.RegionID]CachedSummary),
		maxAge:  maxAge,
	}
}

func (c *SummaryCache) Update(summary leaderboard.RegionSummary, now time.Time) {
	if summary.Region == c.self {
		return
	}

	c.entries[summary.Region] = CachedSummary{
		Summary:    summary,
		ReceivedAt: now,
	}
}

func (c *SummaryCache) IsFresh(region leaderboard.RegionID, now time.Time) bool {

	entry, ok := c.entries[region]
	if !ok {
		return false
	}

	return now.Sub(entry.ReceivedAt) <= c.maxAge
}

func (c *SummaryCache) ActiveSummaries(now time.Time) []leaderboard.RegionSummary {
	out := make([]leaderboard.RegionSummary, 0)

	for _, entry := range c.entries {
		if now.Sub(entry.ReceivedAt) <= c.maxAge {
			out = append(out, entry.Summary)
		}
	}

	return out
}
