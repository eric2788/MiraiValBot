package valorant

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/eric2788/MiraiValBot/redis"
	"github.com/google/uuid"
)

const shortenPuuidKey = "valorant:puuid_short_list"

var uuidCache = make(map[string]*AccountInfo)

type (
	AccountInfo struct {
		Name    string
		Tag     string
		PUuid   string
		Display string
	}

	Statistics struct {
		KDRatio             float64
		HeadshotRate        float64
		AvgScore            float64
		DamagePerRounds     float64
		KillsPerRounds      float64
		MostUsedWeapon      string
		TotalFriendlyDamage int64
		TotalFriendlyKills  int
		WinRate             float64
		TotalMatches        int
	}

	FriendlyFireInfo struct {
		FriendlyFire
		Deaths int
		Kills  int
	}

	Performance struct {
		UserName    string
		Character   string
		CurrentTier string
		Killed      int
		Deaths      int
		Assists     int
	}
)

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

			for _, killEvent := range playerStats.KillEvents {

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

func GetStatistics(name, tag, filter string, region Region) (*Statistics, error) {

	if filter != "" && !AllowedModes.Contains(strings.ToLower(filter)) {
		return nil, fmt.Errorf("无效的模式: %s, 可用模式: %s", filter, strings.Join(AllowedModes.ToArr(), ", "))
	}

	ac, err := GetAccountDetails(name, tag)
	if err != nil {
		return nil, err
	}

	matches, err := GetMatchHistories(name, tag, region)
	if err != nil {
		return nil, err
	}

	totalShots, totalHeadShots := 0, 0
	totalKills, totalDeaths := 0, 0
	totalScores := 0
	totalDamage, totalRounds := int64(0), 0
	ffDamage, ffKills := int64(0), 0
	usedWeapons := make(map[string]int)
	wins := 0

	totalMatches := 0
	for _, match := range matches {

		// have filter and the mode of match data is not matched
		if filter != "" && !strings.EqualFold(match.MetaData.Mode, filter) {
			continue
		}

		// 👇 Friendly fire

		info := GetFriendlyFireInfo(&match)

		if ff, ok := info[ac.PUuid]; ok {
			ffDamage += int64(ff.Outgoing)
			ffKills += ff.Kills
		}

		totalRounds += match.MetaData.RoundsPlayed

		// 👇 kda and hs rate

		players := match.Players["all_players"]

		for _, player := range players {
			if player.PUuid == ac.PUuid {

				totalShots += player.Stats.BodyShots + player.Stats.LegShots + player.Stats.Headshots
				totalHeadShots += player.Stats.Headshots

				totalKills += player.Stats.Kills
				totalDeaths += player.Stats.Deaths

				totalScores += player.Stats.Score

				totalDamage += player.DamageMade
				break
			}
		}

		// 👇 most used weapons

		for _, rounds := range match.Rounds {

			for _, player := range rounds.PlayerStats {

				if player.PlayerPUuid == ac.PUuid {

					if count, ok := usedWeapons[player.Economy.Weapon.Name]; ok {
						usedWeapons[player.Economy.Weapon.Name] = count + 1
					} else {
						usedWeapons[player.Economy.Weapon.Name] = 1
					}

					break
				}

			}
		}

		// 👇 win rates

		// if deathmatch and rank 1
		if strings.ToLower(match.MetaData.Mode) == "deathmatch" && GetDeathMatchRanking(&match)[0].PUuid == ac.PUuid {
			wins += 1
		} else { // check team wins
			red := match.Teams["red"]
			blue := match.Teams["blue"]

			if team, err := FoundPlayerInTeam(ac.PUuid, &match); err != nil {
				return nil, fmt.Errorf("寻找 %s#%s 的队伍所属时出现错误: %v", ac.Name, ac.Tag, err)
			} else if (red.HasWon && strings.ToLower(team) == "red") || (blue.HasWon && strings.ToLower(team) == "blue") {
				wins += 1
			}
		}

		totalMatches++
	}

	mostUsedWeapon, t := "无", 0
	for weapon, times := range usedWeapons {
		if times > t {
			mostUsedWeapon, t = weapon, times
		}
	}

	return &Statistics{
		KDRatio:             float64(totalKills) / float64(totalDeaths),
		HeadshotRate:        float64(totalHeadShots) / float64(totalShots) * 100,
		AvgScore:            float64(totalScores) / float64(len(matches)),
		DamagePerRounds:     float64(totalDamage) / float64(totalRounds),
		KillsPerRounds:      float64(totalKills) / float64(totalRounds),
		WinRate:             float64(wins) / float64(len(matches)) * 100,
		TotalFriendlyDamage: ffDamage,
		TotalFriendlyKills:  ffKills,
		MostUsedWeapon:      mostUsedWeapon,
	}, nil
}

var seasonRegex = regexp.MustCompile(`^[e](\d+)[a](\d+)$`)

func findEposideAct(season string) (ep int, act int) {
	finds := seasonRegex.FindStringSubmatch(season)

	if len(finds) != 3 {
		ep, act = 0, 0
		return
	}

	ep, err := strconv.Atoi(finds[1])
	if err != nil {
		ep = 0
	}

	act, err = strconv.Atoi(finds[2])
	if err != nil {
		act = 0
	}
	return
}

func ShortenUUID(puuid string) (int64, error) {
	if _, err := uuid.Parse(puuid); err != nil {
		return -1, err
	}
	if err := redis.ListAdd(shortenPuuidKey, puuid); err != nil && err != redis.ListExists {
		return -1, err
	}
	return redis.ListPos(shortenPuuidKey, puuid)
}

func ShortenUUIDs(puuids []string) (map[string]int64, map[string]error) {

	var errMap = make(map[string]error)
	var uuidMap = make(map[string]int64)
	for _, puuid := range puuids {

		if _, err := uuid.Parse(puuid); err != nil {
			errMap[puuid] = err
			continue
		}

		if err := redis.ListAdd(shortenPuuidKey, puuid); err != nil && err != redis.ListExists {
			errMap[puuid] = err
			continue
		}

	}

	for _, puuid := range puuids {

		if index, err := redis.ListPos(shortenPuuidKey, puuid); err != nil {
			errMap[puuid] = err
			continue
		} else {
			uuidMap[puuid] = index
		}
	}

	return uuidMap, errMap
}

func GetRealId(id string) (string, error) {
	// already is uuid
	if _, err := uuid.Parse(id); err == nil {
		return id, nil
	}
	pos, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return "", err
	}
	return redis.ListIndex(shortenPuuidKey, pos)
}

func SortSeason(seasons map[string]MMRV2SeasonDetails) []string {

	keys := make([]string, 0, len(seasons))
	for season := range seasons {
		keys = append(keys, season)
	}

	sort.Slice(keys, func(i, j int) bool {

		iep, iact := findEposideAct(keys[i])
		jep, jact := findEposideAct(keys[j])

		if iep == jep {
			return iact > jact
		} else {
			return iep > jep
		}
	})

	return keys
}

func GetPerformances(data *MatchData, name, tag string) ([]Performance, error) {

	ac, err := GetAccountDetails(name, tag)
	if err != nil {
		return nil, err
	}

	exist := false
	for _, player := range data.Players["all_players"] {
		if player.PUuid == ac.PUuid {
			exist = true
			break
		}
	}

	if !exist {
		return []Performance{}, nil
	}

	performanceMap := make(map[string]*Performance)

	getPerformanceMap := func(id, name string) *Performance {
		if d, ok := performanceMap[id]; ok {
			return d
		} else {
			d := &Performance{Killed: 0, CurrentTier: "Unrated", Deaths: 0, Assists: 0, UserName: name}
			performanceMap[id] = d
			return d
		}
	}

	for _, round := range data.Rounds {
		for _, pStats := range round.PlayerStats {
			for _, kEvent := range pStats.KillEvents {

				// ignore friendly kills
				if kEvent.VictimTeam == kEvent.KillerTeam {
					continue
				}

				// killer is target
				if kEvent.KillerPUuid == ac.PUuid {
					per := getPerformanceMap(kEvent.VictimPUuid, kEvent.VictimDisplayName)
					per.Killed += 1
				} else if kEvent.VictimPUuid == ac.PUuid { // victim target
					per := getPerformanceMap(kEvent.KillerPUuid, kEvent.KillerDisplayName)
					per.Deaths += 1
				} else {

					for _, assistant := range kEvent.Assistants {

						// target is assistant
						if assistant.AssistantPUuid == ac.PUuid {
							per := getPerformanceMap(kEvent.VictimPUuid, kEvent.VictimDisplayName)
							per.Assists += 1
							break
						}

					}
				}

			}
		}
	}

	for _, player := range data.Players["all_players"] {
		if per, ok := performanceMap[player.PUuid]; ok {
			per.CurrentTier = player.CurrentTierPatched
			per.Character = player.Character
		}
	}

	performances := make([]Performance, 0)

	for _, per := range performanceMap {
		performances = append(performances, *per)
	}

	sort.Slice(performances, func(i, j int) bool {
		pi, pj := performances[i], performances[j]
		return ((pi.Killed + pi.Assists/2) - pi.Deaths) > ((pj.Killed + pi.Assists/2) - pj.Deaths)
	})

	return performances, nil
}
