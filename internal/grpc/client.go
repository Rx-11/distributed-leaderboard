package grpc

import (
	"context"
	"log"
	"time"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
	pb "github.com/Rx-11/distributed-leaderboard/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conns map[string]pb.ReplicationServiceClient
}

func NewClient(peerAddresses []string) *Client {
	conns := make(map[string]pb.ReplicationServiceClient)

	for _, addr := range peerAddresses {
		conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Printf("Failed to dial peer %s: %v", addr, err)
			continue
		}
		conns[addr] = pb.NewReplicationServiceClient(conn)
	}

	return &Client{conns: conns}
}

func (c *Client) Broadcast(ctx context.Context, summary leaderboard.RegionSummary) {
	req := domainToProto(summary)

	for addr, client := range c.conns {
		go func(a string, cl pb.ReplicationServiceClient) {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()

			_, err := cl.PublishSummary(ctx, req)
			if err != nil {
				log.Printf("RPC failure to %s: %v", a, err)
			}
		}(addr, client)
	}
}

func domainToProto(s leaderboard.RegionSummary) *pb.RegionSummary {
	pbEntries := make([]*pb.Entry, len(s.TopK.Entries))
	for i, e := range s.TopK.Entries {
		pbEntries[i] = &pb.Entry{UserId: e.UserID, Score: e.Score}
	}

	pbBuckets := make([]*pb.HistogramBucket, len(s.Histogram.Buckets))
	for i, b := range s.Histogram.Buckets {
		pbBuckets[i] = &pb.HistogramBucket{
			LowerBound: b.LowerBound,
			UpperBound: b.UpperBound,
			Count:      int32(b.Count),
		}
	}

	return &pb.RegionSummary{
		RegionId:  string(s.Region),
		Epoch:     s.Epoch,
		IsFinal:   s.IsFinal,
		UserCount: int32(s.UserCount),
		TopK: &pb.TopKSummary{
			Epoch:   s.TopK.Epoch,
			Entries: pbEntries,
		},
		Histogram: &pb.HistogramSummary{
			Epoch:   s.Histogram.Epoch,
			Buckets: pbBuckets,
		},
	}
}
