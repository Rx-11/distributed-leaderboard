package leaderboard

import "time"

type Snapshot struct {
	Epoch uint64
	Order []Entry
}

type GlobalSnapshot struct {
	FinalizedAt     time.Time
	TopK            []Entry
	IncludedRegions []RegionID
	TotalUserCount  int
}

func (lb *Leaderboard) Snapshot() Snapshot {
	orderCopy := make([]Entry, len(lb.order))
	copy(orderCopy, lb.order)

	return Snapshot{
		Epoch: lb.epoch,
		Order: orderCopy,
	}
}
