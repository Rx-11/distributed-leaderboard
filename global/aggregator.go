package global

import (
	"sort"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

type GlobalTopKResult struct {
	Epochs  []uint64
	Entries []leaderboard.Entry
}

type GlobalRankEstimate struct {
	Rank  int
	Exact bool
}

func ComputeGlobalTopK(summaries []leaderboard.RegionSummary, k int, mode ConsistencyMode) (GlobalTopKResult, error) {

	all := make([]leaderboard.Entry, 0)

	epochs := make([]uint64, 0, len(summaries))

	for _, s := range summaries {
		epochs = append(epochs, s.Epoch)
		all = append(all, s.TopK.Entries...)
	}

	if err := CheckEpochAlignment(epochs, mode); err != nil {
		return GlobalTopKResult{}, err
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Score != all[j].Score {
			return all[i].Score > all[j].Score
		}
		return all[i].UserID < all[j].UserID
	})

	if k > len(all) {
		k = len(all)
	}

	return GlobalTopKResult{
		Epochs:  epochs,
		Entries: all[:k],
	}, nil
}

func EstimateGlobalRank(user leaderboard.Entry, localRank int, summaries []leaderboard.RegionSummary) GlobalRankEstimate {
	rank := localRank
	for _, s := range summaries {
		for _, b := range s.Histogram.Buckets {
			if b.LowerBound > user.Score {
				rank += b.Count
			}
		}
	}

	return GlobalRankEstimate{
		Rank:  rank,
		Exact: false,
	}
}
