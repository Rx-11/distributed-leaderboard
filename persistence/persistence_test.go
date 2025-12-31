package persistence

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

func TestExportSnapshotToCSV(t *testing.T) {

	now := time.Now()
	snapshot := &leaderboard.GlobalSnapshot{
		FinalizedAt: now,
		TopK: []leaderboard.Entry{
			{UserID: "alice", Score: 500},
			{UserID: "bob", Score: 450},
			{UserID: "charlie", Score: 450},
		},
		IncludedRegions: []leaderboard.RegionID{"us-east", "eu-west"},
		TotalUserCount:  1000,
	}

	var buf bytes.Buffer
	err := ExportSnapshotToCSV(&buf, snapshot)
	if err != nil {
		t.Fatalf("ExportSnapshotToCSV failed: %v", err)
	}

	csvOutput := buf.String()
	lines := strings.Split(strings.TrimSpace(csvOutput), "\n")

	if len(lines) != 4 {
		t.Fatalf("expected 4 lines (1 header + 3 rows), got %d", len(lines))
	}

	expectedHeader := "Rank,UserID,Score,FinalizedAt"
	if lines[0] != expectedHeader {
		t.Errorf("header mismatch.\nGot: %s\nWant: %s", lines[0], expectedHeader)
	}

	if !strings.Contains(lines[1], "1,alice,500") {
		t.Errorf("row 1 incorrect: %s", lines[1])
	}
}
