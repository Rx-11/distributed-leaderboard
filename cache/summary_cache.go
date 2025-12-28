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
	self     leaderboard.RegionID
	entries  map[leaderboard.RegionID]CachedSummary
	freshTTL time.Duration
	staleTTL time.Duration
}

func NewSummaryCache(self leaderboard.RegionID, freshTTL time.Duration, staleTTL time.Duration) *SummaryCache {
	return &SummaryCache{
		self:     self,
		entries:  make(map[leaderboard.RegionID]CachedSummary),
		freshTTL: freshTTL,
		staleTTL: staleTTL,
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

	return now.Sub(entry.ReceivedAt) <= c.freshTTL
}

func (c *SummaryCache) ActiveSummaries(now time.Time) []leaderboard.RegionSummary {
	out := make([]leaderboard.RegionSummary, 0)

	for _, entry := range c.entries {
		if now.Sub(entry.ReceivedAt) <= c.freshTTL {
			out = append(out, entry.Summary)
		}
	}

	return out
}

func (c *SummaryCache) Entries() map[leaderboard.RegionID]CachedSummary {
	return c.entries
}

func (c *SummaryCache) AllSummaries() []CachedSummary {
	summaries := []CachedSummary{}
	for _, entry := range c.entries {
		summaries = append(summaries, entry)
	}

	return summaries
}

func (c *SummaryCache) FreshTTL() time.Duration {
	return c.freshTTL
}

func (c *SummaryCache) StaleTTL() time.Duration {
	return c.staleTTL
}

func (c *SummaryCache) Get(region leaderboard.RegionID) (CachedSummary, bool) {
	entry, ok := c.entries[region]
	return entry, ok
}
