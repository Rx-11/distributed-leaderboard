package global

import (
	"testing"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

func TestGlobalTopK(t *testing.T) {
	r1 := leaderboard.New()
	r2 := leaderboard.New()

	r1.UpdateScore("alice", 100)
	r1.UpdateScore("bob", 200)

	r2.UpdateScore("charlie", 150)
	r2.UpdateScore("dave", 300)

	summaries := []leaderboard.RegionSummary{
		r1.RegionSummary(2),
		r2.RegionSummary(2),
	}

	result, err := ComputeGlobalTopK(summaries, 3, Fast)

	if err != nil {
		t.Fatalf("compute global Top K returned error")
	}

	if result.Entries[0].UserID != "dave" {
		t.Fatalf("expected dave as global #1")
	}
}
