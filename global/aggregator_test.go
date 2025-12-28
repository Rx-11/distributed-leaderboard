package global

import (
	"testing"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

func TestGlobalTopK(t *testing.T) {
	r1 := leaderboard.New("us-east")
	r2 := leaderboard.New("us-west")

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

func TestRankEstimateNoDoubleCount(t *testing.T) {
	r1 := leaderboard.New("us-east")
	r2 := leaderboard.New("eu-west")

	r1.UpdateScore("alice", 100)
	r2.UpdateScore("bob", 200)

	s1 := r1.RegionSummary(10)
	s2 := r2.RegionSummary(10)

	est := EstimateGlobalRank(
		leaderboard.Entry{UserID: "alice", Score: 100},
		1,
		"us-east",
		[]leaderboard.RegionSummary{s1, s2},
	)

	if est.Rank != 2 {
		t.Fatalf("expected rank 2, got %d", est.Rank)
	}
}
