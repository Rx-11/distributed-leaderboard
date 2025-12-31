package persistence

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

func ExportSnapshotToCSV(w io.Writer, snapshot *leaderboard.GlobalSnapshot) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"Rank", "UserID", "Score", "FinalizedAt"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for i, entry := range snapshot.TopK {
		rank := i + 1

		row := []string{
			strconv.Itoa(rank),
			entry.UserID,
			strconv.FormatInt(entry.Score, 10),
			snapshot.FinalizedAt.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row for rank %d: %w", rank, err)
		}
	}

	return nil
}
