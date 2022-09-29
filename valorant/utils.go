package valorant

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

func GetDeathMatchRanking(data *MatchData) []MatchPlayer {
	players := data.Players["all_players"]
	sort.Slice(players, func(i, j int) bool {
		return players[i].Stats.Kills > players[j].Stats.Kills
	})
	return players
}

func GetMatchRanking(data *MatchData) []MatchPlayer {
	players := data.Players["all_players"]
	sort.Slice(players, func(i, j int) bool {
		return players[i].Stats.Score > players[j].Stats.Score
	})
	return players
}

func GetRankingFromPlayers(players []MatchPlayer, id string) (int, *MatchPlayer) {
	for i, player := range players {
		if _, _, err := ParseNameTag(id); (err == nil && fmt.Sprintf("%s#%s", player.Name, player.Tag) == id) || player.PUuid == id {
			return i + 1, &player
		}
	}
	return -1, nil
}

func ParseNameTag(nameTag string) (name string, tag string, err error) {
	parts := strings.Split(nameTag, "#")
	if len(parts) != 2 {
		return "", "", errors.New(fmt.Sprintf("名称格式不正确: %s", nameTag))
	}
	return parts[0], parts[1], nil
}

func FoundPlayerInTeam(nameTag string, data *MatchData) (string, error){
	if len(data.Teams) == 0 {
		return "", nil
	}
	_, player := GetRankingFromPlayers(data.Players["all_players"], nameTag)
	if player == nil {
		return "", errors.New(fmt.Sprintf("在该排名中找不到玩家: %s", nameTag))
	}
	return player.Team, nil
}

func GetDefuseCount(data *MatchData, id string) int {
	defuse := 0
	for _, round := range data.Rounds {
		if _, _, err := ParseNameTag(id); err == nil && round.DefuseEvents.DefusedBy.DisplayName == id {
			defuse += 1
		} else if round.DefuseEvents.DefusedBy.PUuid == id {
			defuse += 1
		}
	}
	return defuse
}

func GetPlantCount(data *MatchData, id string) int {
	plant := 0
	for _, round := range data.Rounds {
		if _, _, err := ParseNameTag(id); err == nil && round.PlantEvents.PlantedBy.DisplayName == id {
			plant += 1
		} else if round.PlantEvents.PlantedBy.PUuid == id {
			plant += 1
		}
	}
	return plant
}
