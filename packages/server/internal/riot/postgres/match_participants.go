package postgres

import (
	"context"
	"fmt"

	"github.com/galchammat/kadeem/internal/riot/models"
	"github.com/lib/pq"
)

func (s *DB) SaveMatchParticipantBatch(ctx context.Context, participants []models.MatchParticipantSummary) error {
	if len(participants) == 0 {
		return nil
	}

	matchIDs := make([]int64, len(participants))
	championIDs := make([]int, len(participants))
	champLevels := make([]int, len(participants))
	kills := make([]int, len(participants))
	deaths := make([]int, len(participants))
	assists := make([]int, len(participants))
	totalMinionsKilled := make([]int, len(participants))
	doubleKills := make([]int, len(participants))
	tripleKills := make([]int, len(participants))
	quadraKills := make([]int, len(participants))
	pentaKills := make([]int, len(participants))
	item0s := make([]int, len(participants))
	item1s := make([]int, len(participants))
	item2s := make([]int, len(participants))
	item3s := make([]int, len(participants))
	item4s := make([]int, len(participants))
	item5s := make([]int, len(participants))
	item6s := make([]int, len(participants))
	summoner1IDs := make([]int, len(participants))
	summoner2IDs := make([]int, len(participants))
	lanes := make([]string, len(participants))
	participantIDs := make([]int, len(participants))
	puuids := make([]string, len(participants))
	riotIDGameNames := make([]string, len(participants))
	riotIDTaglines := make([]string, len(participants))
	totalDamageDealtToChampions := make([]int, len(participants))
	totalDamageTaken := make([]int, len(participants))
	wins := make([]bool, len(participants))

	for i, participant := range participants {
		matchIDs[i] = participant.GameID
		championIDs[i] = participant.ChampionID
		champLevels[i] = participant.ChampLevel
		kills[i] = participant.Kills
		deaths[i] = participant.Deaths
		assists[i] = participant.Assists
		totalMinionsKilled[i] = participant.TotalMinionsKilled
		doubleKills[i] = participant.DoubleKills
		tripleKills[i] = participant.TripleKills
		quadraKills[i] = participant.QuadraKills
		pentaKills[i] = participant.PentaKills
		item0s[i] = participant.Item0
		item1s[i] = participant.Item1
		item2s[i] = participant.Item2
		item3s[i] = participant.Item3
		item4s[i] = participant.Item4
		item5s[i] = participant.Item5
		item6s[i] = participant.Item6
		summoner1IDs[i] = participant.Summoner1ID
		summoner2IDs[i] = participant.Summoner2ID
		lanes[i] = participant.Lane
		participantIDs[i] = participant.ParticipantID
		puuids[i] = participant.PUUID
		riotIDGameNames[i] = participant.RiotIDGameName
		riotIDTaglines[i] = participant.RiotIDTagline
		totalDamageDealtToChampions[i] = participant.TotalDamageDealtToChampions
		totalDamageTaken[i] = participant.TotalDamageTaken
		wins[i] = participant.Win
	}

	_, err := s.db.SQL.ExecContext(ctx, `
		INSERT INTO participants (
			match_id,
			champion_id,
			champ_level,
			kills,
			deaths,
			assists,
			total_minions_killed,
			double_kills,
			triple_kills,
			quadra_kills,
			penta_kills,
			item0,
			item1,
			item2,
			item3,
			item4,
			item5,
			item6,
			summoner1_id,
			summoner2_id,
			lane,
			participant_id,
			puuid,
			riot_id_game_name,
			riot_id_tagline,
			total_damage_dealt_to_champions,
			total_damage_taken,
			win
		)
		SELECT *
		FROM unnest(
			$1::bigint[],
			$2::integer[],
			$3::integer[],
			$4::integer[],
			$5::integer[],
			$6::integer[],
			$7::integer[],
			$8::integer[],
			$9::integer[],
			$10::integer[],
			$11::integer[],
			$12::integer[],
			$13::integer[],
			$14::integer[],
			$15::integer[],
			$16::integer[],
			$17::integer[],
			$18::integer[],
			$19::integer[],
			$20::integer[],
			$21::text[],
			$22::integer[],
			$23::text[],
			$24::text[],
			$25::text[],
			$26::integer[],
			$27::integer[],
			$28::boolean[]
		) AS batch(
			match_id,
			champion_id,
			champ_level,
			kills,
			deaths,
			assists,
			total_minions_killed,
			double_kills,
			triple_kills,
			quadra_kills,
			penta_kills,
			item0,
			item1,
			item2,
			item3,
			item4,
			item5,
			item6,
			summoner1_id,
			summoner2_id,
			lane,
			participant_id,
			puuid,
			riot_id_game_name,
			riot_id_tagline,
			total_damage_dealt_to_champions,
			total_damage_taken,
			win
		)
		ON CONFLICT (match_id, participant_id) DO UPDATE SET
			champion_id = EXCLUDED.champion_id,
			champ_level = EXCLUDED.champ_level,
			kills = EXCLUDED.kills,
			deaths = EXCLUDED.deaths,
			assists = EXCLUDED.assists,
			total_minions_killed = EXCLUDED.total_minions_killed,
			double_kills = EXCLUDED.double_kills,
			triple_kills = EXCLUDED.triple_kills,
			quadra_kills = EXCLUDED.quadra_kills,
			penta_kills = EXCLUDED.penta_kills,
			item0 = EXCLUDED.item0,
			item1 = EXCLUDED.item1,
			item2 = EXCLUDED.item2,
			item3 = EXCLUDED.item3,
			item4 = EXCLUDED.item4,
			item5 = EXCLUDED.item5,
			item6 = EXCLUDED.item6,
			summoner1_id = EXCLUDED.summoner1_id,
			summoner2_id = EXCLUDED.summoner2_id,
			lane = EXCLUDED.lane,
			puuid = EXCLUDED.puuid,
			riot_id_game_name = EXCLUDED.riot_id_game_name,
			riot_id_tagline = EXCLUDED.riot_id_tagline,
			total_damage_dealt_to_champions = EXCLUDED.total_damage_dealt_to_champions,
			total_damage_taken = EXCLUDED.total_damage_taken,
			win = EXCLUDED.win
	`,
		pq.Array(matchIDs),
		pq.Array(championIDs),
		pq.Array(champLevels),
		pq.Array(kills),
		pq.Array(deaths),
		pq.Array(assists),
		pq.Array(totalMinionsKilled),
		pq.Array(doubleKills),
		pq.Array(tripleKills),
		pq.Array(quadraKills),
		pq.Array(pentaKills),
		pq.Array(item0s),
		pq.Array(item1s),
		pq.Array(item2s),
		pq.Array(item3s),
		pq.Array(item4s),
		pq.Array(item5s),
		pq.Array(item6s),
		pq.Array(summoner1IDs),
		pq.Array(summoner2IDs),
		pq.Array(lanes),
		pq.Array(participantIDs),
		pq.Array(puuids),
		pq.Array(riotIDGameNames),
		pq.Array(riotIDTaglines),
		pq.Array(totalDamageDealtToChampions),
		pq.Array(totalDamageTaken),
		pq.Array(wins),
	)
	if err != nil {
		return fmt.Errorf("save match participant batch: %w", err)
	}

	return nil
}
