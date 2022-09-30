package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/imgtxt"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/paste"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/redis"
	"github.com/eric2788/MiraiValBot/valorant"
	"github.com/eric2788/common-utils/datetime"
)

func info(args []string, source *command.MessageSource) error {
	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}
	info, err := valorant.GetAccountDetails(name, tag)
	if err != nil {
		return err
	}
	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextfLn("%s 的账户资讯:", fmt.Sprintf("%s#%s", info.Name, info.Tag)))
	msg.Append(qq.NewTextfLn("用户ID: %s", info.PUuid))
	msg.Append(qq.NewTextfLn("区域: %s", info.Region))
	msg.Append(qq.NewTextfLn("等级: %d", info.AccountLevel))
	msg.Append(qq.NewTextfLn("最新API刷取时间: %s (%s)", datetime.FormatSeconds(info.LastUpdateRaw), info.LastUpdate))
	img, err := qq.NewImageByUrl(info.Card["small"])
	if err != nil {
		logger.Errorf("无法获取用户卡片: %v", err)
	} else {
		msg.Append(img)
	}
	return qq.SendWithRandomRiskyStrategy(msg)
}

func forceUpdate(args []string, source *command.MessageSource) error {
	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}
	msg := qq.CreateReply(source.Message)
	if err := valorant.UpdateAccountDetails(name, tag); err == nil {
		msg.Append(qq.NewTextf("强制更新用户资讯成功。"))
	} else {
		msg.Append(qq.NewTextf("强制更新用户资讯失败: %v", err))
	}

	return qq.SendGroupMessage(msg)
}

func status(args []string, source *command.MessageSource) error {
	status, err := valorant.GetGameStatus(valorant.AsiaSpecific)
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	if len(status.Data.Incidents) == 0 && len(status.Data.Maintenances) == 0 {
		msg.Append(qq.NewTextLn("目前没有任何维护或者故障。"))
	} else {
		for i, incident := range status.Data.Incidents {
			msg.Append(qq.NewTextfLn("=========== 事故 (%d) ===========", i))
			appendDetails(msg, incident)
		}
		for i, maintenance := range status.Data.Maintenances {
			msg.Append(qq.NewTextfLn("=========== 维护 (%d) ===========", i))
			appendDetails(msg, maintenance)
		}
	}
	return qq.SendWithRandomRiskyStrategy(msg)
}

func matches(args []string, source *command.MessageSource) error {

	if err := qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取对战资料..."))); err != nil {
		logger.Errorf("發送預備索取消息失敗: %v", err)
	}

	info, err := valorant.GetAccountInfo(args[0])
	if err != nil {
		return err
	}
	matches, err := valorant.GetMatchHistories(info.Name, info.Tag, valorant.AsiaSpecific)
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 最近的对战:", info.Display))
	for _, match := range matches {
		// empty match id
		if match.MetaData.MatchId == "" {
			continue
		}
		msg.Append(qq.NewTextLn("===================="))
		msg.Append(qq.NewTextfLn("对战ID: %s", match.MetaData.MatchId))
		msg.Append(qq.NewTextfLn("对战模式: %s", match.MetaData.Mode))
		msg.Append(qq.NewTextfLn("对战开始时间: %s", datetime.FormatSeconds(match.MetaData.GameStart)))
		msg.Append(qq.NewTextfLn("对战时长: %s", formatDuration(match.MetaData.GameLength)))
		msg.Append(qq.NewTextfLn("对战地图: %s", match.MetaData.Map))
		msg.Append(qq.NewTextfLn("回合总数: %d", match.MetaData.RoundsPlayed))
		msg.Append(qq.NewTextfLn("服务器节点: %s", match.MetaData.Cluster))
		msg.Append(qq.NewTextfLn("对战结果: %s", formatResult(match, info.PUuid)))
	}

	msg.Append(qq.NewTextLn("===================="))
	msg.Append(qq.NewTextfLn("输入 !val leaderboard <对战ID> 查看排行榜"))
	msg.Append(qq.NewTextfLn("输入 !val players <对战ID> 查看对战玩家"))
	msg.Append(qq.NewTextfLn("输入 !val rounds <对战ID> 查看对战回合"))

	return qq.SendWithRandomRiskyStrategy(msg)
}

func match(args []string, source *command.MessageSource) error {
	match, err := valorant.GetMatchDetails(args[0])
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("对战ID: %s", match.MetaData.MatchId))
	msg.Append(qq.NewTextfLn("对战模式: %s", match.MetaData.Mode))
	msg.Append(qq.NewTextfLn("对战开始时间: %s", datetime.FormatSeconds(match.MetaData.GameStart)))
	msg.Append(qq.NewTextfLn("对战时长: %s", formatDuration(match.MetaData.GameLength)))
	msg.Append(qq.NewTextfLn("对战地图: %s", match.MetaData.Map))
	msg.Append(qq.NewTextfLn("回合总数: %d", match.MetaData.RoundsPlayed))
	msg.Append(qq.NewTextfLn("服务器节点: %s", match.MetaData.Cluster))
	msg.Append(qq.NewTextfLn("对战结果: %s", formatResultObjective(match)))
	msg.Append(qq.NewTextfLn("输入 !val leaderboard %s 查看排行榜", match.MetaData.MatchId))
	msg.Append(qq.NewTextfLn("输入 !val players %s 查看对战玩家", match.MetaData.MatchId))
	msg.Append(qq.NewTextfLn("输入 !val rounds %s 查看对战回合", match.MetaData.MatchId))
	return qq.SendWithRandomRiskyStrategy(msg)
}

func matchPlayers(args []string, source *command.MessageSource) error {

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取对战玩家的资料..")))

	match, err := valorant.GetMatchDetails(args[0])
	if err != nil {
		return err
	}

	img, err := generateMatchPlayersImage(match)
	if err != nil {
		return err
	}

	sending := message.NewSendingMessage().Append(img)
	return qq.SendWithRandomRiskyStrategy(sending)
}

func leaderboard(args []string, source *command.MessageSource) error {

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取对战排行榜的资料..")))

	match, err := valorant.GetMatchDetails(args[0])
	if err != nil {
		return err
	}

	msg := message.NewSendingMessage()

	msg.Append(qq.NewTextfLn("对战 %s 的玩家排行榜", match.MetaData.MatchId))
	if strings.ToLower(match.MetaData.Mode) == "deathmatch" {
		players := valorant.GetDeathMatchRanking(match)
		for i, player := range players {
			msg.Append(qq.NewTextLn("===================="))
			msg.Append(qq.NewTextfLn("%d. - %s", i+1, fmt.Sprintf("%s#%s", player.Name, player.Tag)))
			msg.Append(qq.NewTextfLn("均分: %d", player.Stats.Score))
			msg.Append(qq.NewTextfLn("K/D/A: %d/%d/%d", player.Stats.Kills, player.Stats.Deaths, player.Stats.Assists))
		}
	} else {
		players := valorant.GetMatchRanking(match)
		ffMap := valorant.GetFriendlyFireInfo(match)

		getFFDamage := func(player valorant.MatchPlayer) int {
			if info, ok := ffMap[player.PUuid]; ok {
				return info.Outgoing
			} else {
				return player.Behaviour.FriendlyFire.Outgoing
			}
		}

		getFFKills := func(player valorant.MatchPlayer) int {
			if info, ok := ffMap[player.PUuid]; ok {
				return info.Kills
			} else {
				return 0
			}
		}

		for i, player := range players {
			totalShots := player.Stats.BodyShots + player.Stats.LegShots + player.Stats.Headshots
			msg.Append(qq.NewTextLn("===================="))
			msg.Append(qq.NewTextfLn("%d. - %s", i+1, fmt.Sprintf("%s#%s", player.Name, player.Tag)))
			msg.Append(qq.NewTextfLn("均分: %d", player.Stats.Score))
			msg.Append(qq.NewTextfLn("K/D/A: %d/%d/%d", player.Stats.Kills, player.Stats.Deaths, player.Stats.Assists))
			msg.Append(qq.NewTextfLn("爆头率: %.1f%%", formatPercentageInt(player.Stats.Headshots, totalShots)))
			msg.Append(qq.NewTextfLn("队友伤害: %d", getFFDamage(player)))
			msg.Append(qq.NewTextfLn("队友误杀: %d", getFFKills(player)))
			msg.Append(qq.NewTextfLn("装包次数: %d", valorant.GetPlantCount(match, player.PUuid)))
			msg.Append(qq.NewTextfLn("拆包次数: %d", valorant.GetDefuseCount(match, player.PUuid)))
		}
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

func matchRounds(args []string, source *command.MessageSource) error {

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取对战回合的资料..")))

	match, err := valorant.GetMatchDetails(args[0])
	if err != nil {
		return err
	}

	// 过滤死斗
	if strings.ToLower(match.MetaData.Mode) == "deathmatch" {
		return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("死斗没有可以查看的对战回合资讯。")))
	}

	msg := message.NewSendingMessage()

	for i, round := range match.Rounds {
		msg.Append(qq.NewTextfLn("第 %d 回合 (胜者: %s, 胜利类型: %s)", i+1, round.WinningTeam, round.EndType))

		if round.BombPlanted {
			msg.Append(qq.NewTextfLn("\t安装炸弹:"))
			msg.Append(qq.NewTextfLn("\t\t安装玩家: %s", round.PlantEvents.PlantedBy.DisplayName))
			msg.Append(qq.NewTextfLn("\t\t安装队伍: %s", round.PlantEvents.PlantedBy.Team))
			msg.Append(qq.NewTextfLn("\t\t安装地点: %s", round.PlantEvents.PlantSite))
		}

		if round.BombDefused {
			msg.Append(qq.NewTextLn("\t解除炸弹:"))
			msg.Append(qq.NewTextfLn("\t\t解除玩家: %s", round.DefuseEvents.DefusedBy.DisplayName))
			msg.Append(qq.NewTextfLn("\t\t解除队伍: %s", round.DefuseEvents.DefusedBy.Team))
		}

		for _, playerStats := range round.PlayerStats {
			msg.Append(qq.NewTextfLn("\t%s(队伍:%s) 在该回合的战绩:", playerStats.PlayerDisplayName, playerStats.PlayerTeam))

			msg.Append(qq.NewTextfLn("\t\tAFK: %t", playerStats.WasAfk))
			msg.Append(qq.NewTextfLn("\t\t被惩罚: %t", playerStats.WasPenalized))
			msg.Append(qq.NewTextfLn("\t\t回合花费: $%d (剩余 $%d)", playerStats.Economy.Spent, playerStats.Economy.Remaining))
			msg.Append(qq.NewTextfLn("\t\t武器: %s", playerStats.Economy.Weapon.Weapon.Name))
			msg.Append(qq.NewTextfLn("\t\t装备: %s", playerStats.Economy.Weapon.Armor.Name))

			if playerStats.Damage > 0 {
				msg.Append(qq.NewTextfLn("\t\t分别伤害:"))
				for _, damageEvent := range playerStats.DamageEvents {
					msg.Append(qq.NewTextfLn("\t\t\t%s:", damageEvent.ReceiverDisplayName))
					msg.Append(qq.NewTextfLn("\t\t\t\t伤害: %d (%.1f%%)", damageEvent.Damage, formatPercentageInt(damageEvent.Damage, playerStats.Damage)))
					msg.Append(qq.NewTextfLn("\t\t\t\t所在队伍: %s", damageEvent.ReceiverTeam))
					msg.Append(qq.NewTextfLn("\t\t\t\t伤害分布:"))
					total := damageEvent.BodyShots + damageEvent.HeadShots + damageEvent.LegShots
					msg.Append(qq.NewTextfLn("\t\t\t\t\t头部: %d (%.1f%%)", damageEvent.HeadShots, formatPercentageInt(damageEvent.HeadShots, total)))
					msg.Append(qq.NewTextfLn("\t\t\t\t\t身体: %d (%.1f%%)", damageEvent.BodyShots, formatPercentageInt(damageEvent.BodyShots, total)))
					msg.Append(qq.NewTextfLn("\t\t\t\t\t腿部: %d (%.1f%%)", damageEvent.LegShots, formatPercentageInt(damageEvent.LegShots, total)))
				}
			}

			if len(playerStats.KillsEvents) > 0 {
				msg.Append(qq.NewTextLn("\t\t分别击杀:"))
				for _, killEvent := range playerStats.KillsEvents {
					msg.Append(qq.NewTextfLn("\t\t\t%s:", killEvent.VictimDisplayName))
					msg.Append(qq.NewTextfLn("\t\t\t\t所在队伍: %s", killEvent.VictimTeam))
					msg.Append(qq.NewTextfLn("\t\t\t\t击杀使用武器: %s", killEvent.DamageWeaponName))
					msg.Append(qq.NewTextfLn("\t\t\t\t右键开火: %t", killEvent.SecondaryFireMode))

					assistantArr := make([]string, len(killEvent.Assistants))
					for i, assistant := range killEvent.Assistants {
						assistantArr[i] = fmt.Sprintf("%s(%s)", assistant.AssistantDisplayName, assistant.AssistantTeam)
					}

					if len(assistantArr) > 0 {
						msg.Append(qq.NewTextfLn("\t\t\t\t助攻者: %s", strings.Join(assistantArr, ", ")))
					}
				}
			}
		}
	}

	content := strings.Join(qq.ParseMsgContent(msg.Elements).Texts, "")

	key, err := paste.CreatePaste("plain", content)
	if err != nil {
		return err
	}

	sending := qq.CreateReply(source.Message).
		Append(qq.NewTextLn("链接只能使用一次，五分钟后过期。")).
		Append(qq.NewTextfLn("https://pasteme.cn#%s", key)).
		Append(qq.NewTextf("如过期，请重新输入指令生成。"))
	return qq.SendWithRandomRiskyStrategy(sending)
}

// mmr get MMRV1Details
func mmr(args []string, source *command.MessageSource) error {
	parts := strings.Split(args[0], "#")
	mmr, err := valorant.GetMMRDetailsV1(parts[0], parts[1], valorant.AsiaSpecific)
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 的 MMR 资料:", args[0]))
	msg.Append(qq.NewTextfLn("目前段位: %s", mmr.CurrentTierPatched))
	msg.Append(qq.NewTextfLn("目前段位分数: %d/100", mmr.RankingInTier))
	msg.Append(qq.NewTextfLn("上一次的分数变更: %d", mmr.MMRChangeToLastGame))
	msg.Append(qq.NewTextfLn("ELO: %d", mmr.Elo))
	img, err := qq.NewImageByUrl(mmr.Images["small"])
	if err == nil {
		msg.Append(img)
	} else {
		logger.Errorf("无法获取段位图片: %v", err)
	}
	return qq.SendWithRandomRiskyStrategy(msg)
}

// mmrHistories get MMRHistories
func mmrHistories(args []string, source *command.MessageSource) error {
	return qq.SendGroupMessage(message.NewSendingMessage().Append(qq.NewTextLn("此指令暂不可用")))
}

// mmrBySeason GetMMRDetailsBySeason
func mmrBySeason(args []string, source *command.MessageSource) error {
	return qq.SendGroupMessage(message.NewSendingMessage().Append(qq.NewTextLn("此指令暂不可用")))
}

func localize(args []string, source *command.MessageSource) error {
	return qq.SendGroupMessage(message.NewSendingMessage().Append(qq.NewTextLn("此指令暂不可用")))
}

var (
	infoCommand         = command.NewNode([]string{"info", "资讯"}, "查询玩家账户资讯", false, info, "<名称#Tag>")
	forceUpdateCommand  = command.NewNode([]string{"update", "更新"}, "强制更新玩家资讯", false, forceUpdate, "<名称#Tag>")
	statusCommand       = command.NewNode([]string{"status", "状态"}, "查询状态", false, status)
	matchesCommand      = command.NewNode([]string{"matches", "对战历史"}, "查询对战历史", false, matches)
	matchCommand        = command.NewNode([]string{"match", "对战"}, "查询对战详情", false, match, "<对战ID>")
	leaderboardCommand  = command.NewNode([]string{"leaderboard", "排行榜"}, "查询对战排行榜", false, leaderboard, "<对战ID>")
	matchPlayerscommand = command.NewNode([]string{"players", "玩家"}, "查询对战玩家资讯", false, matchPlayers, "<对战ID>")
	matchRoundsCommand  = command.NewNode([]string{"rounds", "回合"}, "查询对战回合资讯", false, matchRounds, "<对战ID>")
	mmrCommand          = command.NewNode([]string{"mmr", "段位"}, "查询段位", false, mmr, "<名称#Tag>")
	mmrHistoriesCommand = command.NewNode([]string{"mmrHistories", "段位历史"}, "查询段位历史", false, mmrHistories, "<名称#Tag>")
	mmrBySeasonCommand  = command.NewNode([]string{"mmrBySeason", "赛季段位"}, "查询赛季段位", false, mmrBySeason, "<名称#Tag>", "<赛季>")
	localizeCommand     = command.NewNode([]string{"localize", "本地化"}, "更新i18n内容", true, localize)
)

var valorantCommand = command.NewParent([]string{"valorant", "val", "瓦罗兰", "瓦"}, "valorant指令",
	infoCommand,
	forceUpdateCommand,
	statusCommand,
	matchesCommand,
	matchCommand,
	leaderboardCommand,
	matchPlayerscommand,
	matchRoundsCommand,
	mmrCommand,
	mmrHistoriesCommand,
	mmrBySeasonCommand,
	localizeCommand,
)

func init() {
	command.AddCommand(valorantCommand)
}

// ===================================
//
//            Util Functions
//
// ===================================

func formatResultObjective(data *valorant.MatchData) string {
	mode := strings.ToLower(data.MetaData.Mode)

	switch mode {
	case "deathmatch":
		ranking := valorant.GetDeathMatchRanking(data)
		if len(ranking) == 0 {
			return "没有名次"
		}
		player := ranking[0]
		return fmt.Sprintf("胜出者: %s (K %d | D %d | A %d, 分数: %d)",
			fmt.Sprintf("%s#%s", player.Name, player.Tag),
			player.Stats.Kills,
			player.Stats.Deaths,
			player.Stats.Assists,
			player.Stats.Score,
		)
	case "unrated":
		fallthrough
	case "competitive":
		fallthrough
	case "custom game":
		red := data.Teams["red"]
		blue := data.Teams["blue"]
		return fmt.Sprintf("Red %d : %d Blue", red.RoundsWon, blue.RoundsWon)
	}
	return fmt.Sprintf("(错误: 不支援的模式 %s)", data.MetaData.Mode)
}

func formatResult(data valorant.MatchData, name string) string {
	mode := strings.ToLower(data.MetaData.Mode)

	switch mode {
	case "deathmatch":
		ranking := valorant.GetDeathMatchRanking(&data)
		rank, player := valorant.GetRankingFromPlayers(ranking, name)
		if rank == -1 {
			return fmt.Sprintf("在该排名中找不到玩家: %s", name)
		}
		return fmt.Sprintf("第 %d 名 (K %d | D %d | A %d, 分数: %d)",
			rank,
			player.Stats.Kills,
			player.Stats.Deaths,
			player.Stats.Assists,
			player.Stats.Score,
		)
	case "unrated":
		fallthrough
	case "competitive":
		fallthrough
	case "custom game":
		red := data.Teams["red"]
		blue := data.Teams["blue"]
		team, err := valorant.FoundPlayerInTeam(name, &data)
		if err != nil {
			return fmt.Sprintf("(错误: %s)", err.Error())
		}
		return fmt.Sprintf("Red %d : %d Blue (用户所在队伍: %s)", red.RoundsWon, blue.RoundsWon, team)
	}
	return fmt.Sprintf("(错误: 不支援的模式 %s)", data.MetaData.Mode)
}

func formatPercentage(part, total int64) float64 {
	return float64(part) / float64(total) * 100
}

func formatPercentageInt(part, total int) float64 {
	return float64(part) / float64(total) * 100
}

func formatTime(timeStr string) string {
	ti, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		logger.Errorf("无法解析时间: %s, 将返回厡讯息", timeStr)
		return timeStr
	}
	return datetime.FormatSeconds(int64(ti.Second()))
}

func formatDuration(milis int64) string {
	min := float64(milis / 1000 / 60)
	return fmt.Sprintf("%.1f 分钟", min)
}

var localePriorities = []string{"zh-CN", "zh-TW", "en-US"}

func formatTranslate(translates []valorant.I18NContent) string {
	for _, translate := range translates {
		for _, locale := range localePriorities {
			if translate.Locale == locale {
				return translate.Content
			}
		}
	}
	logger.Errorf("无法找到任何翻译，将返回原始内容")
	return translates[0].Content
}

func appendDetails(msg *message.SendingMessage, maintenance valorant.MaintainInfo) {
	msg.Append(qq.NewTextfLn("开始时间: %s", formatTime(maintenance.CreatedAt)))
	msg.Append(qq.NewTextfLn("预计完成时间: %s", formatTime(maintenance.ArchiveAt)))
	msg.Append(qq.NewTextfLn("目前状态: %s", maintenance.MaintenanceStatus))
	msg.Append(qq.NewTextfLn("等级: %s", maintenance.IncidentSeverity))
	msg.Append(qq.NewTextfLn("标题: %s", formatTranslate(maintenance.Titles)))
	msg.Append(qq.NewTextLn("内容:"))
	for _, update := range maintenance.Updates {
		msg.Append(qq.NewTextfLn("> %s", formatTranslate(update.Translations)))
		msg.Append(qq.NewTextfLn("	创建于: %s", formatTime(update.CreatedAt)))
		msg.Append(qq.NewTextfLn("	更新于: %s", formatTime(update.UpdatedAt)))
		msg.Append(qq.NewTextfLn("	发布者: %s", formatTime(update.Author)))
	}
}

func generateMatchPlayersImage(match *valorant.MatchData) (*message.GroupImageElement, error) {

	key := fmt.Sprintf("valorant:match_player:%s", match.MetaData.MatchId)

	var imgCache = &message.GroupImageElement{}
	if exist, err := redis.Get(key, imgCache); err == nil && exist {
		return imgCache, nil
	} else if err != nil {
		logger.Warnf("从 redis 获取对战玩家图片(%s)时出现错误: %v, 将重新生成。", match.MetaData.MatchId, err)
	}

	ffInfo := valorant.GetFriendlyFireInfo(match)
	ranking := valorant.GetMatchRanking(match)

	msg, err := imgtxt.NewPrependMessage()
	if err != nil {
		return nil, err
	}
	for i, player := range ranking {
		msg.Append(qq.NewTextfLn("\t第 %d 名: %s", i+1, fmt.Sprintf("%s#%s", player.Name, player.Tag)))

		// 基本资料
		msg.Append(qq.NewTextLn("\t基本资料:"))
		msg.Append(qq.NewTextfLn("\t\tKDA: %d | %d | %d", player.Stats.Kills, player.Stats.Deaths, player.Stats.Assists))
		msg.Append(qq.NewTextfLn("\t\t分数: %d", player.Stats.Score))
		msg.Append(qq.NewTextfLn("\t\t使用角色: %s", player.Character))
		msg.Append(qq.NewTextfLn("\t\t所在队伍: %s", player.Team))

		// 击中分布
		total := player.Stats.BodyShots + player.Stats.Headshots + player.Stats.LegShots
		msg.Append(qq.NewTextLn("\t击中次数分布"))
		msg.Append(qq.NewTextfLn("\t\t头部: %.1f%% (%d次)", formatPercentageInt(player.Stats.Headshots, total), player.Stats.Headshots))
		msg.Append(qq.NewTextfLn("\t\t身体: %.1f%% (%d次)", formatPercentageInt(player.Stats.BodyShots, total), player.Stats.BodyShots))
		msg.Append(qq.NewTextfLn("\t\t腿部: %.1f%% (%d次)", formatPercentageInt(player.Stats.LegShots, total), player.Stats.LegShots))

		// 行为
		friendlyFire := &valorant.FriendlyFireInfo{FriendlyFire: player.Behaviour.FriendlyFire}
		if ff, ok := ffInfo[player.PUuid]; ok {
			friendlyFire = ff
		} else {
			logger.Warnf("找不到 %s#%s 的隊友傷害行為資訊。", player.Name, player.Tag)
		}
		msg.Append(qq.NewTextLn("\t行为:"))
		msg.Append(qq.NewTextfLn("\t\tAFK回合次数: %.2f", player.Behaviour.AfkRounds))
		msg.Append(qq.NewTextfLn("\t\t误击队友伤害: %d", friendlyFire.Outgoing))
		msg.Append(qq.NewTextfLn("\t\t误杀队友次数: %d", friendlyFire.Kills))
		msg.Append(qq.NewTextfLn("\t\t被误击队友伤害: %d", friendlyFire.Incoming))
		msg.Append(qq.NewTextfLn("\t\t被误杀队友次数: %d", friendlyFire.Deaths))
		msg.Append(qq.NewTextfLn("\t\t拆包次数: %d", valorant.GetDefuseCount(match, player.PUuid)))
		msg.Append(qq.NewTextfLn("\t\t装包次数: %d", valorant.GetPlantCount(match, player.PUuid)))

		//技能使用
		total = 0
		for _, times := range player.AbilityCasts {
			total += times
		}

		msg.Append(qq.NewTextLn("\t技能使用次数分布:"))
		msg.Append(qq.NewTextfLn("\t\t技能 Q: %d次 (%.1f%%)", player.AbilityCasts["q_cast"], formatPercentageInt(player.AbilityCasts["q_cast"], total)))
		msg.Append(qq.NewTextfLn("\t\t技能 E: %d次 (%.1f%%)", player.AbilityCasts["e_cast"], formatPercentageInt(player.AbilityCasts["e_cast"], total)))
		msg.Append(qq.NewTextfLn("\t\t技能 C: %d次 (%.1f%%)", player.AbilityCasts["c_cast"], formatPercentageInt(player.AbilityCasts["c_cast"], total)))
		msg.Append(qq.NewTextfLn("\t\t技能 X: %d次 (%.1f%%)", player.AbilityCasts["x_cast"], formatPercentageInt(player.AbilityCasts["x_cast"], total)))

		// 经济
		msg.Append(qq.NewTextLn("\t经济:"))
		msg.Append(qq.NewTextfLn("\t\t总支出 $%d", player.Economy.Spent.OverAll))
		msg.Append(qq.NewTextfLn("\t\t平均支出 $%d", player.Economy.Spent.Average))

		// 伤害
		totalDamage := player.DamageReceived + player.DamageMade
		msg.Append(qq.NewTextLn("\t伤害分布:"))
		msg.Append(qq.NewTextfLn("\t\t总承受 %d (%.1f%%)", player.DamageReceived, formatPercentage(player.DamageReceived, totalDamage)))
		msg.Append(qq.NewTextfLn("\t\t总伤害 %d (%.1f%%)", player.DamageMade, formatPercentage(player.DamageMade, totalDamage)))
	}

	img, err := msg.ToGroupImageElement()

	if err != nil {
		return nil, err
	}

	if err = redis.Store(key, img); err != nil {
		logger.Warnf("储存对战玩家图片(%s)到redis时出现错误: %v", match.MetaData.MatchId, err)
	} else {
		logger.Infof("储存对战玩家图片(%s)到redis成功", match.MetaData.MatchId)
	}

	return img, nil
}
