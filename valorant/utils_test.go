package valorant

import (
	"fmt"
	"github.com/eric2788/MiraiValBot/redis"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDeathMatchRanking(t *testing.T) {
	match, err := GetMatchDetailsAPI("a4e99fec-647d-4a15-9015-967c8e29355a")
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	players := GetDeathMatchRanking(match)
	for i, player := range players {
		t.Logf("%d. %s (%d kills)(score: %d)",
			i+1,
			fmt.Sprintf("%s#%s", player.Name, player.Tag),
			player.Stats.Kills,
			player.Stats.Score,
		)
	}

	assert.Equal(t, 12, len(players))
}

func aTestShortenUUIDs(t *testing.T) {
	redis.Init()
	matches, err := GetMatchHistories("suou", "9035", AsiaSpecific)
	if err != nil {
		t.Log(err)
	} else {
		var ids = make([]string, len(matches))
		for i, match := range matches {
			ids[i] = match.MetaData.MatchId
		}
		results, errs := ShortenUUIDs(ids)
		if len(errs) > 0 {
			for _, err := range errs {
				t.Log(err)
			}
		}
		for id, result := range results {
			t.Logf("%s -> %d", id, result)
		}
	}
}

func TestPercentageDisplay(t *testing.T) {
	total, a, b := 23, 11, 12
	t.Logf("A %.1f%% (%d) B %.1f%% (%d) - %.0f", float64(a)/float64(total)*100, a, float64(b)/float64(total)*100, b, float64(b)/float64(total))
}

func TestSortSeason(t *testing.T) {
	seasonKeys := []string{
		"e1a1", "e2a1", "e2a3", "e5a1", "e4a1", "e4a2", "e4a3", "e5a2", "e5a3", "e1a2", "e1a3", "e2a2", "e3a1", "e3a3",
	}
	seasons := make(map[string]MMRV2SeasonDetails)

	for _, key := range seasonKeys {
		seasons[key] = MMRV2SeasonDetails{}
	}

	t.Log("before:")
	for season := range seasons {
		t.Log(season)
	}
	sorted := SortSeason(seasons)
	t.Log("after:")
	for _, season := range sorted {
		t.Log(season)
	}
}
