package leaderboard

const DefaultTopK = 10
const HistogramBucketSize int64 = 100

type TopKSummary struct {
	Epoch   uint64
	Entries []Entry
}

type HistogramBucket struct {
	LowerBound int64
	UpperBound int64
	Count      int
}

type HistogramSummary struct {
	Epoch   uint64
	Buckets []HistogramBucket
}

type RegionSummary struct {
	Epoch     uint64
	TopK      TopKSummary
	Histogram HistogramSummary
	UserCount int
}
