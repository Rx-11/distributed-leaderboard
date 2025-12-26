package leaderboard

import "testing"

func TestLeaderboardBasic(t *testing.T) {
	lb := New()

	lb.UpdateScore("alice", 100)
	lb.UpdateScore("bob", 200)
	lb.UpdateScore("charlie", 150)

	rank, _ := lb.GetRank("bob")
	if rank != 1 {
		t.Fatalf("expected bob rank 1, got %d", rank)
	}

	top := lb.GetTopK(2)
	if top[0].UserID != "bob" || top[1].UserID != "charlie" {
		t.Fatalf("unexpected top-K ordering")
	}

	neigh := lb.GetNeighborhood("charlie", 1)
	if len(neigh) != 3 {
		t.Fatalf("expected neighborhood size 3")
	}
}

func TestSnapshotImmutability(t *testing.T) {
	lb := New()

	lb.UpdateScore("alice", 100)
	snap := lb.Snapshot()

	lb.UpdateScore("bob", 200)

	if snap.Epoch != 1 {
		t.Fatalf("expected snapshot epoch 1, got %d", snap.Epoch)
	}

	if len(snap.Order) != 1 || snap.Order[0].UserID != "alice" {
		t.Fatalf("snapshot should not change after updates")
	}
}

func TestRegionSummary(t *testing.T) {
	lb := New()

	lb.UpdateScore("alice", 120)
	lb.UpdateScore("bob", 250)
	lb.UpdateScore("charlie", 180)

	summary := lb.RegionSummary(2)

	if summary.UserCount != 3 {
		t.Fatalf("expected user count 3")
	}

	if len(summary.TopK.Entries) != 2 {
		t.Fatalf("expected topK size 2")
	}

	if summary.TopK.Entries[0].UserID != "bob" {
		t.Fatalf("expected bob as top scorer")
	}

	if summary.Epoch != lb.Epoch() {
		t.Fatalf("epoch mismatch")
	}
}
