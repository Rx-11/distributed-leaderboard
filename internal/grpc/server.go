package grpc

import (
	"context"
	"time"

	"github.com/Rx-11/distributed-leaderboard/cache"
	"github.com/Rx-11/distributed-leaderboard/leaderboard"
	pb "github.com/Rx-11/distributed-leaderboard/proto"
)

type Server struct {
	pb.UnimplementedReplicationServiceServer
	cache *cache.SummaryCache
}

func NewServer(cache *cache.SummaryCache) *Server {
	return &Server{cache: cache}
}

func (s *Server) PublishSummary(ctx context.Context, req *pb.RegionSummary) (*pb.PublishResponse, error) {
	domainSummary := protoToDomain(req)
	s.cache.Update(domainSummary, time.Now())
	return &pb.PublishResponse{Success: true}, nil
}

func protoToDomain(req *pb.RegionSummary) leaderboard.RegionSummary {
	entries := make([]leaderboard.Entry, len(req.TopK.Entries))
	for i, e := range req.TopK.Entries {
		entries[i] = leaderboard.Entry{UserID: e.UserId, Score: e.Score}
	}

	buckets := make([]leaderboard.HistogramBucket, len(req.Histogram.Buckets))
	for i, b := range req.Histogram.Buckets {
		buckets[i] = leaderboard.HistogramBucket{
			LowerBound: b.LowerBound,
			UpperBound: b.UpperBound,
			Count:      int(b.Count),
		}
	}

	return leaderboard.RegionSummary{
		Region:    leaderboard.RegionID(req.RegionId),
		Epoch:     req.Epoch,
		IsFinal:   req.IsFinal,
		UserCount: int(req.UserCount),
		TopK: leaderboard.TopKSummary{
			Epoch:   req.TopK.Epoch,
			Entries: entries,
		},
		Histogram: leaderboard.HistogramSummary{
			Epoch:   req.Histogram.Epoch,
			Buckets: buckets,
		},
	}
}
