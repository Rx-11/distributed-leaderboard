package leaderboard

import "errors"

type HistoricalLeaderboards struct {
	seasons map[SeasonID]*SeasonBoard
}

func NewHistoricalLeadboards() *HistoricalLeaderboards {
	return &HistoricalLeaderboards{
		seasons: make(map[SeasonID]*SeasonBoard),
	}
}

func (hl *HistoricalLeaderboards) CreateSeason(s Season) error {
	if _, ok := hl.seasons[s.ID]; ok {
		return errors.New("season already exists")
	}

	hl.seasons[s.ID] = NewSeasonBoard(s)
	return nil
}

func (hl *HistoricalLeaderboards) UpdateScore(seasonID SeasonID, entry Entry) error {

	sb, ok := hl.seasons[seasonID]
	if !ok {
		return errors.New("unknown season")
	}

	return sb.Update(entry)
}

func (hl *HistoricalLeaderboards) FreezeSeason(seasonID SeasonID) error {
	sb, ok := hl.seasons[seasonID]
	if !ok {
		return errors.New("unknown season")
	}
	return sb.Freeze()
}

func (hl *HistoricalLeaderboards) GetSeasonLeaderboard(seasonID SeasonID) (SeasonBoard, error) {

	sb, ok := hl.seasons[seasonID]
	if !ok {
		return SeasonBoard{}, errors.New("unknown season")
	}

	if sb.season.State != SeasonFrozen {
		return SeasonBoard{}, errors.New("season not frozen")
	}

	return *sb, nil
}
