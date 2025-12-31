package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Rx-11/distributed-leaderboard/api"
	"github.com/Rx-11/distributed-leaderboard/cache"
	"github.com/Rx-11/distributed-leaderboard/config"
	"github.com/Rx-11/distributed-leaderboard/global"
	internalgrpc "github.com/Rx-11/distributed-leaderboard/internal/grpc"
	"github.com/Rx-11/distributed-leaderboard/leaderboard"
	pb "github.com/Rx-11/distributed-leaderboard/proto"
	"google.golang.org/grpc"
)

var (
	cfg          *config.Config
	lb           *leaderboard.Leaderboard
	summaryCache *cache.SummaryCache
	grpcClient   *internalgrpc.Client
)

func main() {
	config.Load()
	cfg = config.GetConfig()

	if err := os.MkdirAll(cfg.DataDir, 0755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}

	var err error
	lb, err = leaderboard.New(leaderboard.RegionID(cfg.RegionID), cfg.DataDir)
	if err != nil {
		log.Fatalf("failed to init leaderboard: %v", err)
	}

	summaryCache = cache.NewSummaryCache(leaderboard.RegionID(cfg.RegionID), cfg.FreshTTL, cfg.StaleTTL)

	grpcClient = internalgrpc.NewClient(cfg.Peers)

	go startGRPCServer()

	go startReplicationLoop()

	startHTTPServer()
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen gRPC: %v", err)
	}

	s := grpc.NewServer()
	srv := internalgrpc.NewServer(summaryCache)
	pb.RegisterReplicationServiceServer(s, srv)

	log.Printf("gRPC Internal Server listening on :%d", cfg.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve gRPC: %v", err)
	}
}

func startHTTPServer() {
	http.HandleFunc("/score", handleUpdateScore)
	http.HandleFunc("/topk", handleGetTopK)
	http.HandleFunc("/rank", handleGetRank)

	addr := fmt.Sprintf(":%d", cfg.HTTPPort)
	log.Printf("HTTP Public Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func startReplicationLoop() {
	ticker := time.NewTicker(2 * time.Second)

	for range ticker.C {
		summary := lb.RegionSummary(20)
		grpcClient.Broadcast(context.Background(), summary)
	}
}

func handleUpdateScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user := r.URL.Query().Get("user")
	valStr := r.URL.Query().Get("val")
	val, _ := strconv.ParseInt(valStr, 10, 64)

	if err := lb.UpdateScore(user, val); err != nil {
		http.Error(w, "internal storage error", http.StatusInternalServerError)
		log.Printf("Write failed: %v", err)
		return
	}

	fmt.Fprintf(w, "updated %s to %d in region %s\n", user, val, cfg.RegionID)
}

func handleGetTopK(w http.ResponseWriter, r *http.Request) {
	k, _ := strconv.Atoi(r.URL.Query().Get("k"))
	if k == 0 {
		k = 10
	}

	mode := global.Fast
	if r.URL.Query().Get("mode") == "strict" {
		mode = global.Strict
	}

	resp, err := api.GetGlobalTopK(time.Now(), mode, lb, summaryCache, cfg.TotalRegions, k)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleGetRank(w http.ResponseWriter, r *http.Request) {
	user := r.URL.Query().Get("user")

	mode := global.Fast
	if r.URL.Query().Get("mode") == "strict" {
		mode = global.Strict
	}

	localRank, foundLocal := lb.GetRank(user)
	var score int64
	var found bool

	if foundLocal {
		s, _ := lb.GetScore(user)
		score = s
		found = true
	} else {
		score, found = summaryCache.FindUser(user)
		localRank = 1
	}

	if !found {
		http.Error(w, "user not found in any region (local or top-k cache)", http.StatusNotFound)
		return
	}

	entry := leaderboard.Entry{UserID: user, Score: score}

	targetRegion := lb.Region()
	if !foundLocal {
		targetRegion = ""
	}

	resp, err := api.GetGlobalRank(time.Now(), mode, entry, localRank, targetRegion, lb, summaryCache, cfg.TotalRegions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
