package leaderboard

import (
	"errors"
	"time"

	"github.com/Rx-11/distributed-leaderboard/config"
)

type SeasonBoard struct {
	season Season
	rank   *Leaderboard
}

func NewSeasonBoard(season Season) (*SeasonBoard, error) {
	lb, err := New(RegionID(config.GetConfig().RegionID), season.ID, config.GetConfig().DataDir)
	if err != nil {
		return nil, err
	}
	return &SeasonBoard{
		season: season,
		rank:   lb,
	}, nil
}

func (sb *SeasonBoard) Update(entry Entry) error {
	if sb.season.State != SeasonLive {
		return errors.New("season is frozen")
	}
	sb.rank.UpdateScore(entry.UserID, entry.Score)
	return nil
}

func (sb *SeasonBoard) Freeze() error {
	if sb.season.State != SeasonLive {
		return errors.New("season is already frozen")
	}
	sb.season.EndTime = time.Now()
	sb.season.State = SeasonFrozen
	return nil
}

func (sb *SeasonBoard) GetRank() *Leaderboard {
	return sb.rank
}

func (sb *SeasonBoard) GetSeason() *Season {
	return &sb.season
}
