package matchsync

import (
	riotmodels "github.com/galchammat/kadeem/internal/riot/models"
)

func mapMatchDetails(matchDetails riotmodels.MatchDetails) (riotmodels.MatchSummary, []riotmodels.MatchParticipantSummary) {
	summary := riotmodels.MatchSummary{
		ID:        matchDetails.Info.ID,
		Region:    matchDetails.Info.Region,
		StartedAt: matchDetails.Info.StartedAt,
		Duration:  matchDetails.Info.Duration,
		QueueID:   matchDetails.Info.QueueID,
	}
	return summary, matchDetails.Info.Participants
}
