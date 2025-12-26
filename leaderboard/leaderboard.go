package leaderboard

import (
	"sort"
)

type Entry struct {
	UserID string
	Score  int64
}

type Leaderboard struct {
	scores map[string]int64
	order  []Entry
	epoch  uint64
}

func New() *Leaderboard {
	return &Leaderboard{
		scores: make(map[string]int64),
		order:  make([]Entry, 0),
		epoch: 0,
	}
}

func (lb *Leaderboard) Epoch() uint64 {
    return lb.epoch
}

func (lb *Leaderboard) UpdateScore(userID string, score int64) {
	lb.scores[userID] = score
	lb.epoch++
	lb.rebuild()
}

func (lb *Leaderboard) GetRank(userID string) (int, bool) {
	for i, e := range lb.order {
		if e.UserID == userID {
			return i + 1, true
		}
	}
	return 0, false
}

func (lb *Leaderboard) GetTopK(k int) []Entry {
	if k > len(lb.order) {
		k = len(lb.order)
	}
	return append([]Entry(nil), lb.order[:k]...)
}

func (lb *Leaderboard) GetNeighborhood(userID string, n int) []Entry {
	idx := -1
	for i, e := range lb.order {
		if e.UserID == userID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return nil
	}

	start := idx - n
	if start < 0 {
		start = 0
	}
	end := idx + n + 1
	if end > len(lb.order) {
		end = len(lb.order)
	}

	return append([]Entry(nil), lb.order[start:end]...)
}

func (lb *Leaderboard) rebuild() {
	lb.order = lb.order[:0]
	for userID, score := range lb.scores {
		lb.order = append(lb.order, Entry{
			UserID: userID,
			Score:  score,
		})
	}

	sort.Slice(lb.order, func(i, j int) bool {
		if lb.order[i].Score != lb.order[j].Score {
			return lb.order[i].Score > lb.order[j].Score
		}
		return lb.order[i].UserID < lb.order[j].UserID
	})
}
