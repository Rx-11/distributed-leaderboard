package cache

import (
	"testing"
	"time"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

func TestSummaryCacheFreshness(t *testing.T) {
	now := time.Now()

	cache := NewSummaryCache("us-east", time.Second)

	s := leaderboard.RegionSummary{
		Region: "eu-west",
		Epoch:  1,
	}

	cache.Update(s, now)

	if !cache.IsFresh("eu-west", now.Add(500*time.Millisecond)) {
		t.Fatalf("expected fresh")
	}

	if cache.IsFresh("eu-west", now.Add(2*time.Second)) {
		t.Fatalf("expected stale")
	}
}
