package riot

import "github.com/galchammat/kadeem/internal/models"

func (r *RiotClient) GetMatchByID(region, matchID string) (models.MatchSummary, error) {
	return models.MatchSummary{}, nil
}
