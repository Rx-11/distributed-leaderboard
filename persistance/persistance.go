package persistence

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/Rx-11/distributed-leaderboard/leaderboard"
)

func ExportSnapshotToCSV(w io.Writer, snapshot *leaderboard.GlobalSeasonSnapshot) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"Rank", "UserID", "Score", "SeasonID", "FinalizedAt"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	for i, entry := range snapshot.TopK {
		rank := i + 1

		row := []string{
			strconv.Itoa(rank),
			entry.UserID,
			strconv.FormatInt(entry.Score, 10),
			string(snapshot.SeasonID),
			snapshot.FinalizedAt.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row for rank %d: %w", rank, err)
		}
	}

	return nil
}

func ExportHistoricalLeaderBoardToCSV(w io.Writer, hl *leaderboard.HistoricalLeaderboards, season leaderboard.SeasonID) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"Rank", "UserID", "Score", "SeasonID", "FinalizedAt"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	seasonBoard, err := hl.GetSeasonLeaderboard(season)
	if err != nil {
		return fmt.Errorf("failed to get season leaderboard: %s", season)
	}

	for i, entry := range seasonBoard.GetRank().GetFull() {
		rank := i + 1

		row := []string{
			strconv.Itoa(rank),
			entry.UserID,
			strconv.FormatInt(entry.Score, 10),
			string(season),
			seasonBoard.GetSeason().EndTime.Format(time.RFC3339),
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row for rank %d: %w", rank, err)
		}
	}

	return nil
}
