package leaderboard

import "time"

type SeasonID string

type SeasonState int

const (
	SeasonLive SeasonState = iota
	SeasonFrozen
)

type Season struct {
	ID        SeasonID
	StartTime time.Time
	EndTime   time.Time
	State     SeasonState
}
