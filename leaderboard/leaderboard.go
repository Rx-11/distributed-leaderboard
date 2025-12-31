package leaderboard

import (
	"fmt"
	"log"
	"sort"

	"github.com/Rx-11/distributed-leaderboard/wal"
)

type Entry struct {
	UserID string
	Score  int64
}

type Leaderboard struct {
	region RegionID
	scores map[string]int64
	order  []Entry
	epoch  uint64
	wal    *wal.WAL
}

func New(region RegionID, dataDir string) (*Leaderboard, error) {
	lb := &Leaderboard{
		region: region,
		scores: make(map[string]int64),
		order:  make([]Entry, 0),
		epoch:  0,
	}

	walPath := fmt.Sprintf("%s/%s.log", dataDir, region)
	log.Printf("Recovering state from %s...", walPath)
	entries, err := wal.Recover(walPath)
	if err != nil {
		return nil, fmt.Errorf("recovery failed: %w", err)
	}

	for _, e := range entries {
		// Apply directly to memory without writing to WAL again
		lb.scores[e.UserID] = e.Score
		lb.epoch++ // Epoch advances on recovery too
	}
	lb.rebuild() // Sort the recovered data
	log.Printf("Recovered %d entries.", len(entries))

	// 3. Open WAL for new writes
	wal, err := wal.OpenWAL(walPath)
	if err != nil {
		return nil, err
	}
	lb.wal = wal

	return lb, nil
}

func (lb *Leaderboard) Region() RegionID {
	return lb.region
}

func (lb *Leaderboard) Epoch() uint64 {
	return lb.epoch
}

func (lb *Leaderboard) UpdateScore(userID string, score int64) error {

	if err := lb.wal.Write(userID, score); err != nil {
		return fmt.Errorf("wal write error: %w", err)
	}

	lb.scores[userID] = score
	lb.epoch++
	lb.rebuild()

	return nil
}

func (lb *Leaderboard) GetRank(userID string) (int, bool) {
	for i, e := range lb.order {
		if e.UserID == userID {
			return i + 1, true
		}
	}
	return 0, false
}

func (lb *Leaderboard) GetScore(userID string) (int64, bool) {
	score, ok := lb.scores[userID]
	return score, ok
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

func (lb *Leaderboard) TopKSummary(k int) TopKSummary {
	if k > len(lb.order) {
		k = len(lb.order)
	}

	entries := make([]Entry, k)
	copy(entries, lb.order[:k])

	return TopKSummary{
		Epoch:   lb.epoch,
		Entries: entries,
	}
}

func (lb *Leaderboard) HistogramSummary() HistogramSummary {
	buckets := make(map[int64]int)

	for _, score := range lb.scores {
		bucket := score / HistogramBucketSize
		buckets[bucket]++
	}

	result := make([]HistogramBucket, 0, len(buckets))
	for bucket, count := range buckets {
		lower := bucket * HistogramBucketSize
		upper := lower + HistogramBucketSize - 1
		result = append(result, HistogramBucket{
			LowerBound: lower,
			UpperBound: upper,
			Count:      count,
		})
	}

	return HistogramSummary{
		Epoch:   lb.epoch,
		Buckets: result,
	}
}

func (lb *Leaderboard) RegionSummary(k int) RegionSummary {
	return RegionSummary{
		Region:    lb.region,
		Epoch:     lb.epoch,
		TopK:      lb.TopKSummary(k),
		Histogram: lb.HistogramSummary(),
		UserCount: len(lb.scores),
	}
}

func (lb *Leaderboard) GetFull() []Entry {
	return lb.order
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
