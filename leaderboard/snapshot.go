package leaderboard

type Snapshot struct {
	Epoch uint64
	Order []Entry
}

func (lb *Leaderboard) Snapshot() Snapshot {
	orderCopy := make([]Entry, len(lb.order))
	copy(orderCopy, lb.order)

	return Snapshot{
		Epoch: lb.epoch,
		Order: orderCopy,
	}
}
