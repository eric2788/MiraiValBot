package valorant

import (
	"fmt"
	"sort"
	"strings"
)

var uuidCache = make(map[string]*AccountInfo)

type AccountInfo struct {
	Name    string
	Tag     string
	PUuid   string
	Display string
}

type FriendlyFireInfo struct {
	FriendlyFire
	Deaths int
	Kills  int
}

func GetFriendlyFireInfo(data *MatchData) map[string]*FriendlyFireInfo {
	var infoMap = make(map[string]*FriendlyFireInfo)

	getInfo := func(id string) *FriendlyFireInfo {
		if value, ok := infoMap[id]; ok {
			return value
		} else {
			info := &FriendlyFireInfo{}
			infoMap[id] = info
			return info
		}
	}

	for _, round := range data.Rounds {
		for _, playerStats := range round.PlayerStats {

			info := getInfo(playerStats.PlayerPUuid)

			for _, damageEvent := range playerStats.DamageEvents {

				victimInfo := getInfo(damageEvent.ReceiverPUuid)

				// friendly fire damage! and not himself
				if damageEvent.ReceiverTeam == playerStats.PlayerTeam && playerStats.PlayerPUuid != damageEvent.ReceiverPUuid {
					info.Outgoing += damageEvent.Damage
					victimInfo.Incoming += damageEvent.Damage
				}

			}

			for _, killEvent := range playerStats.KillsEvents {

				victimInfo := getInfo(killEvent.KillerPUuid)

				// friendly kill!
				if killEvent.VictimTeam == playerStats.PlayerTeam && playerStats.PlayerPUuid != killEvent.KillerPUuid {
					info.Kills += 1
					victimInfo.Deaths += 1
				}
			}
		}
	}

	return infoMap
}

func GetAccountInfo(id string) (*AccountInfo, error) {
	name, tag, err := ParseNameTag(id)
	if err != nil {
		return nil, err
	}
	if cache, ok := uuidCache[fmt.Sprintf("%s#%s", name, tag)]; ok {
		return cache, nil
	} else {
		details, err := GetAccountDetails(name, tag)
		if err != nil {
			return nil, err
		}
		info := &AccountInfo{
			Name:    details.Name,
			Tag:     details.Tag,
			PUuid:   details.PUuid,
			Display: fmt.Sprintf("%s#%s", details.Name, details.Tag),
		}
		uuidCache[fmt.Sprintf("%s#%s", name, tag)] = info
		return info, nil
	}
}

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
		return "", "", fmt.Errorf("名称格式不正确: %s", nameTag)
	}
	return parts[0], parts[1], nil
}

func FoundPlayerInTeam(nameTag string, data *MatchData) (string, error) {
	if len(data.Teams) == 0 {
		return "", nil
	}
	_, player := GetRankingFromPlayers(data.Players["all_players"], nameTag)
	if player == nil {
		return "", fmt.Errorf("在该排名中找不到玩家: %s", nameTag)
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
