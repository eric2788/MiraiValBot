package valorant

import (
	"fmt"
	"testing"

	"github.com/eric2788/MiraiValBot/redis"
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

func TestGetStatistics(t *testing.T) {
	name, tag := "麻將", "4396"
	stats, err := GetStatistics(name, tag, AsiaSpecific)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	t.Logf("%s#%s 在最近五场对战中的统计数据: ", name, tag)
	t.Logf("爆头率: %.2f%%", stats.HeadshotRate)
	t.Logf("胜率: %.f%%", stats.WinRate)
	t.Logf("KD比例: %.2f", stats.KDRatio)
	t.Logf("最常使用武器: %s", stats.MostUsedWeapon)
	t.Logf("平均分数: %.1f", stats.AvgScore)
	t.Logf("每回合平均伤害: %.1f", stats.DamagePerRounds)
	t.Logf("每回合平均击杀: %.1f", stats.KillsPerRounds)
	t.Logf("总队友伤害: %d", stats.TotalFriendlyDamage)
	t.Logf("总队友击杀: %d", stats.TotalFriendlyKills)
}

func TestGetPerformance(t *testing.T) {
	match, err := GetMatchDetailsAPI("c82f5416-a4b6-4720-be13-a05414049210")
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	performances, err := GetPerformances(match, "麻將", "4396")
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}

	if len(performances) == 0 {
		t.Log("target is not in this match.")
		return
	}

	for i, perfor := range performances {
		t.Logf("%d.\t%s\tK:%d\tD:%d\tA:%d\t(%s)\t(%s)", i+1,
			perfor.UserName,
			perfor.Killed,
			perfor.Deaths,
			perfor.Assists,
			perfor.CurrentTier,
			perfor.Character,
		)
	}
}

func TestShortenUUIDs(t *testing.T) {
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
