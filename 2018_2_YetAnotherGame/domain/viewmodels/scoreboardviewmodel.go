package viewmodels

import "2018_2_YetAnotherGame/domain/models"

type ScoreboardPageViewModel struct {
	Scoreboard models.Scoreboard
	CanNext    bool
}
