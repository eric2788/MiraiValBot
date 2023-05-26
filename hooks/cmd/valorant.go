package cmd

import (
	"fmt"

	"math"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/paste"
	"github.com/eric2788/MiraiValBot/services/valorant"
	"github.com/eric2788/common-utils/datetime"

	v "github.com/eric2788/MiraiValBot/hooks/sites/valorant"
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
	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
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
		msg.Append(qq.NewTextLn("目前没有任何维护或者事件。"))
	} else {
		for i, incident := range status.Data.Incidents {
			msg.Append(qq.NewTextfLn("=========== 事件 (%d) ===========", i+1))
			appendDetails(msg, incident)
		}
		for i, maintenance := range status.Data.Maintenances {
			msg.Append(qq.NewTextfLn("=========== 维护 (%d) ===========", i+1))
			appendDetails(msg, maintenance)
		}
	}
	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

func track(args []string, source *command.MessageSource) error {

	name, tag, err := valorant.ParseNameTag(args[0])

	if err != nil {
		return err
	}

	reply := qq.CreateReply(source.Message)
	success, err := v.StartListen(name, tag)

	if err != nil {
		reply.Append(qq.NewTextf("监听玩家时出现错误: %v", err))
	} else if success {
		reply.Append(qq.NewTextf("开始监听玩家 %s#%s", name, tag))
	} else {
		reply.Append(qq.NewTextf("该玩家(%s#%s) 已经启动监听。", name, tag))
	}

	return qq.SendGroupMessage(reply)
}

func untrack(args []string, source *command.MessageSource) error {
	name, tag, err := valorant.ParseNameTag(args[0])

	if err != nil {
		return err
	}

	reply := qq.CreateReply(source.Message)
	success, err := v.StopListen(name, tag)

	if err != nil {
		reply.Append(qq.NewTextf("中止监听玩家时出现错误: %v", err))
	} else if success {
		reply.Append(qq.NewTextf("已中止监听玩家 %s#%s", name, tag))
	} else {
		reply.Append(qq.NewTextf("该玩家(%s#%s) 尚未启动监听。", name, tag))
	}

	return qq.SendGroupMessage(reply)
}

func tracking(args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)
	listening := v.GetListening()
	if len(listening) > 0 {
		reply.Append(qq.NewTextf("正在监听的玩家: %v", strings.Join(listening, ", ")))
	} else {
		reply.Append(message.NewText("没有正在监听的玩家"))
	}

	return qq.SendWithRandomRiskyStrategyRemind(reply, source.Message)
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

	var uuidsToShort = make([]string, len(matches))
	for i, match := range matches {
		uuidsToShort[i] = match.MetaData.MatchId
	}

	shorts := getShortIdsHint(uuidsToShort)

	image := false
	if len(args) > 1 {
		image = args[1] == "image" || strings.HasPrefix(args[1], "图") || args[1] == "true"
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 最近的对战:", info.Display))
	for _, match := range matches {
		// empty match id
		if match.MetaData.MatchId == "" {
			continue
		}

		shortHint := shorts[match.MetaData.MatchId]
		msg.Append(qq.NewTextLn("===================="))
		msg.Append(qq.NewTextfLn("对战ID: %s%s", match.MetaData.MatchId, shortHint))
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
	msg.Append(qq.NewTextfLn("输入 !val performance <对战ID> <名称#Tag> 查看对战玩家表现"))

	if image {
		if err = qq.SendGroupImageText(msg); err == nil {
			return nil
		} else {
			logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
		}
	}

	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

func match(args []string, source *command.MessageSource) error {

	id, err := valorant.GetRealId(args[0])
	if err != nil {
		return fmt.Errorf("id 解析失败: %v", err)
	}

	match, err := valorant.GetMatchDetails(id)
	if err != nil {
		return err
	}

	shortHint, short := getShortIdHint(match.MetaData.MatchId)

	cmdId := match.MetaData.MatchId

	if short > -1 {
		cmdId = fmt.Sprintf("%d", short)
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("对战ID: %s%s", match.MetaData.MatchId, shortHint))
	msg.Append(qq.NewTextfLn("对战模式: %s", match.MetaData.Mode))
	msg.Append(qq.NewTextfLn("对战开始时间: %s", datetime.FormatSeconds(match.MetaData.GameStart)))
	msg.Append(qq.NewTextfLn("对战时长: %s", formatDuration(match.MetaData.GameLength)))
	msg.Append(qq.NewTextfLn("对战地图: %s", match.MetaData.Map))
	msg.Append(qq.NewTextfLn("回合总数: %d", match.MetaData.RoundsPlayed))
	msg.Append(qq.NewTextfLn("服务器节点: %s", match.MetaData.Cluster))
	msg.Append(qq.NewTextfLn("对战结果: %s", formatResultObjective(match)))
	msg.Append(qq.NewTextfLn("输入 !val leaderboard %s 查看排行榜", cmdId))
	msg.Append(qq.NewTextfLn("输入 !val players %s 查看对战玩家", cmdId))
	msg.Append(qq.NewTextfLn("输入 !val rounds %s 查看对战回合", cmdId))
	msg.Append(qq.NewTextfLn("输入 !val performance %s <名称#Tag> 查看对战玩家表现", cmdId))
	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

func matchPlayers(args []string, source *command.MessageSource) error {

	id, err := valorant.GetRealId(args[0])
	if err != nil {
		return fmt.Errorf("id 解析失败: %v", err)
	}

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取对战玩家的资料..")))

	match, err := valorant.GetMatchDetails(id)
	if err != nil {
		return err
	}

	sending := generateMatchPlayersLines(match)
	return qq.SendGroupImageText(sending)
}

func leaderboard(args []string, source *command.MessageSource) error {

	id, err := valorant.GetRealId(args[0])
	if err != nil {
		return fmt.Errorf("id 解析失败: %v", err)
	}

	image := false
	if len(args) > 1 {
		image = args[1] == "image" || strings.HasPrefix(args[1], "图") || args[1] == "true"
	}

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取对战排行榜的资料..")))

	match, err := valorant.GetMatchDetails(id)
	if err != nil {
		return err
	}

	go qq.SendRiskyMessageWithFunc(5, 60, func(currentTry int) error {

		msg := message.NewSendingMessage()
		alts := qq.GetRandomMessageByTry(currentTry)

		msg.Append(qq.NewTextfLn("对战 %s 的玩家排行榜", match.MetaData.MatchId))
		if strings.ToLower(match.MetaData.Mode) == "deathmatch" {
			players := valorant.GetDeathMatchRanking(match)
			for i, player := range players {
				msg.Append(qq.NewTextLn("===================="))
				msg.Append(qq.NewTextfLn("%d. - %s (%s)", i+1, fmt.Sprintf("%s#%s", player.Name, player.Tag), player.Character))
				msg.Append(qq.NewTextfLn("均分: %d", player.Stats.Score))
				msg.Append(qq.NewTextfLn("K/D/A: %d/%d/%d (%.2f)", player.Stats.Kills, player.Stats.Deaths, player.Stats.Assists, float64(player.Stats.Kills)/float64(player.Stats.Deaths)))
			}
		} else {
			players := valorant.GetMatchRanking(match)
			ffMap := valorant.GetFriendlyFireInfo(match)

			getFFDamage := func(player valorant.MatchPlayer) int {
				if info, ok := ffMap[player.PUuid]; ok {
					return int(math.Round(info.Outgoing))
				} else {
					return int(math.Round(player.Behaviour.FriendlyFire.Outgoing))
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
				msg.Append(qq.NewTextfLn("%d. - %s (%s)", i+1, fmt.Sprintf("%s#%s", player.Name, player.Tag), player.Character))

				// 如果是競技模式，則顯示段位
				if strings.ToLower(match.MetaData.Mode) == "competitive" {
					msg.Append(qq.NewTextfLn("段位: %s", player.CurrentTierPatched))
				}

				msg.Append(qq.NewTextfLn("队伍: %s", player.Team))
				msg.Append(qq.NewTextfLn("均分: %d", player.Stats.Score))
				msg.Append(qq.NewTextfLn("K/D/A: %d/%d/%d (%.2f)", player.Stats.Kills, player.Stats.Deaths, player.Stats.Assists, float64(player.Stats.Kills)/float64(player.Stats.Deaths)))
				if currentTry <= 4 {
					msg.Append(qq.NewTextfLn("爆头率: %.1f%%", formatPercentageInt(player.Stats.Headshots, totalShots)))
				}
				if currentTry <= 2 {
					msg.Append(qq.NewTextfLn("队友伤害: %d", getFFDamage(player)))
					msg.Append(qq.NewTextfLn("队友误杀: %d", getFFKills(player)))
				}
				if currentTry <= 3 {
					msg.Append(qq.NewTextfLn("装包次数: %d", valorant.GetPlantCount(match, player.PUuid)))
					msg.Append(qq.NewTextfLn("拆包次数: %d", valorant.GetDefuseCount(match, player.PUuid)))
				}
			}
		}

		if len(alts) > 0 {
			msg.Append(qq.NextLn())
		}
		for _, ele := range alts {
			msg.Append(ele)
			msg.Append(qq.NextLn())
		}

		if image {
			if err = qq.SendGroupImageText(msg); err == nil {
				return nil
			} else {
				logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
			}
		}

		return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)

	}, func() {
		// 重试失败后，提示信息被风控
		remind := qq.CreateAtReply(source.Message)
		remind.Append(message.NewText("回应发送失败，可能被风控咯 😔"))
		_ = qq.SendGroupMessageByGroup(source.Message.GroupCode, remind)
	})

	return nil
}

func performances(args []string, source *command.MessageSource) error {

	name, tag, err := valorant.ParseNameTag(args[1])
	if err != nil {
		return err
	}

	id, err := valorant.GetRealId(args[0])
	if err != nil {
		return fmt.Errorf("id 解析失败: %v", err)
	}

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(qq.NewTextf("正在索取玩家 %s 在该对战的表现..", args[1])))

	match, err := valorant.GetMatchDetails(id)
	if err != nil {
		return err
	}

	msg := message.NewSendingMessage()

	performances, err := valorant.GetPerformances(match, name, tag)
	if err != nil {
		return err
	} else if len(performances) == 0 {
		msg.Append(qq.NewTextf("%s 并不在对战 %s 之中。", args[1], args[0]))
		return qq.SendGroupMessage(msg)
	}

	image := false
	if len(args) > 2 {
		image = args[2] == "image" || strings.HasPrefix(args[2], "图") || args[2] == "true"
	}

	msg.Append(qq.NewTextfLn("玩家 %s 在对战 %s 中的击杀表现 (由高到低):", args[1], match.MetaData.MatchId))

	for i, performance := range performances {
		msg.Append(qq.NewTextLn("==================="))
		msg.Append(qq.NewTextfLn("%d. %s (%s)", i+1, performance.UserName, performance.Character))

		if strings.ToLower(match.MetaData.Mode) == "competitive" {
			msg.Append(qq.NewTextfLn("段位: %s", performance.CurrentTier))
		}

		msg.Append(qq.NewTextfLn("K/D/A: %d/%d/%d", performance.Killed, performance.Deaths, performance.Assists))
	}

	if image {
		if err = qq.SendGroupImageText(msg); err == nil {
			return nil
		} else {
			logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
		}
	}

	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

func stats(args []string, source *command.MessageSource) error {

	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取最近十场对战资料的统计数据...")))

	filter := ""
	if len(args) > 1 {
		filter = args[1]
	}

	stats, err := valorant.GetStatistics(name, tag, filter, valorant.AsiaSpecific)
	if err != nil {
		return err
	}

	image := false
	if len(args) > 2 {
		image = args[2] == "image" || strings.HasPrefix(args[2], "图") || args[2] == "true"
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 在最近 %d 场对战中的统计数据: ", args[0], stats.TotalMatches))
	msg.Append(qq.NewTextfLn("爆头率: %.2f%%", stats.HeadshotRate))
	msg.Append(qq.NewTextfLn("胜率: %.f%%", stats.WinRate))
	msg.Append(qq.NewTextfLn("KD比例: %.2f", stats.KDRatio))
	msg.Append(qq.NewTextfLn("最常使用武器: %s", stats.MostUsedWeapon))
	msg.Append(qq.NewTextfLn("平均分数: %.1f", stats.AvgScore))
	msg.Append(qq.NewTextfLn("每回合平均伤害: %.1f", stats.DamagePerRounds))
	msg.Append(qq.NewTextfLn("每回合平均击杀: %.1f", stats.KillsPerRounds))
	msg.Append(qq.NewTextfLn("总队友伤害: %d", stats.TotalFriendlyDamage))
	msg.Append(qq.NewTextfLn("总队友击杀: %d", stats.TotalFriendlyKills))

	if image {
		if err = qq.SendGroupImageText(msg); err == nil {
			return nil
		} else {
			logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
		}
	}

	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

func matchRounds(args []string, source *command.MessageSource) error {

	id, err := valorant.GetRealId(args[0])
	if err != nil {
		return fmt.Errorf("id 解析失败: %v", err)
	}

	go qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("正在索取对战回合的资料..")))

	match, err := valorant.GetMatchDetails(id)
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
			msg.Append(qq.NewTextfLn("\t%s(队伍:%s) 在该回合的表现:", playerStats.PlayerDisplayName, playerStats.PlayerTeam))

			msg.Append(qq.NewTextfLn("\t\tAFK: %t", playerStats.WasAfk))
			msg.Append(qq.NewTextfLn("\t\t被惩罚: %t", playerStats.WasPenalized))
			msg.Append(qq.NewTextfLn("\t\t回合花费: $%d (剩余 $%d)", playerStats.Economy.Spent, playerStats.Economy.Remaining))
			msg.Append(qq.NewTextfLn("\t\t武器: %s", playerStats.Economy.Weapon.Name))
			msg.Append(qq.NewTextfLn("\t\t装备: %s", playerStats.Economy.Weapon.Name))

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

			if playerStats.Kills > 0 {
				msg.Append(qq.NewTextLn("\t\t分别击杀:"))
				for _, killEvent := range playerStats.KillEvents {
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

	pmUrl, err := paste.CreatePasteMe("plain", content)
	if err != nil {
		pmUrl = fmt.Sprintf("(错误: %v)", err)
	}

	pbUrl, err := paste.CreatePasteBin(fmt.Sprintf("%s 的对战回合资讯", match.MetaData.MatchId), content, "yaml")
	if err != nil {
		pbUrl = fmt.Sprintf("(错误: %v)", err)
	}

	sending := qq.CreateReply(source.Message)

	sending.Append(qq.NewTextfLn("PasteMe(国内): %s (五分钟过期 / 阅后即焚)", pmUrl))
	sending.Append(qq.NewTextfLn("PasteBin(国外): %s (一天后过期)", pbUrl))

	return qq.SendWithRandomRiskyStrategyRemind(sending, source.Message)
}

// mmr get MMRV1Details
func mmr(args []string, source *command.MessageSource) error {

	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}

	mmr, err := valorant.GetMMRDetailsV1(name, tag, valorant.AsiaSpecific)
	if err != nil {
		return err
	}

	image := false
	if len(args) > 1 {
		image = args[1] == "image" || strings.HasPrefix(args[1], "图") || args[1] == "true"
	}

	msg := message.NewSendingMessage()

	plus := ""
	if mmr.MMRChangeToLastGame > 0 {
		plus = "+"
	}

	msg.Append(qq.NewTextfLn("%s 的 MMR 资料:", args[0]))
	msg.Append(qq.NewTextfLn("目前段位: %s", mmr.CurrentTierPatched))
	msg.Append(qq.NewTextfLn("目前段位分数: %d/100", mmr.RankingInTier))
	msg.Append(qq.NewTextfLn("上一次的分数变更: %s%d", plus, mmr.MMRChangeToLastGame))
	msg.Append(qq.NewTextfLn("ELO: %d", mmr.Elo))
	img, err := qq.NewImageByUrl(mmr.Images["small"])
	if err == nil {
		msg.Append(img)
	} else {
		logger.Errorf("无法获取段位图片: %v", err)
	}

	if image {
		if err = qq.SendGroupImageText(msg); err == nil {
			return nil
		} else {
			logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
		}
	}

	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

// mmrHistories get MMRHistories
func mmrHistories(args []string, source *command.MessageSource) error {
	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}

	info, err := valorant.GetMMRHistories(name, tag, valorant.AsiaSpecific)
	if err != nil {
		return err
	}

	image := false
	if len(args) > 1 {
		image = args[1] == "image" || strings.HasPrefix(args[1], "图") || args[1] == "true"
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 的 MMR 变更记录: ", fmt.Sprintf("%s#%s", info.Name, info.Tag)))

	for _, data := range info.Data {
		symbol := ""

		if data.MMRChangeToLastGame > 0 {
			symbol = "+"
		}

		msg.Append(qq.NewTextLn("===================="))
		msg.Append(qq.NewTextfLn("对战时间: %s", datetime.FormatSeconds(data.DateRaw)))
		msg.Append(qq.NewTextfLn("段位: %s", data.CurrentTierPatched))
		msg.Append(qq.NewTextfLn("ELO: %d", data.Elo))
		msg.Append(qq.NewTextfLn("分数变更: %s%d", symbol, data.MMRChangeToLastGame))
	}

	if image {
		if err = qq.SendGroupImageText(msg); err == nil {
			return nil
		} else {
			logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
		}
	}

	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

// mmrBySeason GetMMRDetailsBySeason
func mmrBySeason(args []string, source *command.MessageSource) error {
	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}

	data, err := valorant.GetMMRDetailsBySeason(name, tag, args[1], valorant.AsiaSpecific)
	if err != nil {
		return err
	}

	if data.NumberOfGames == 0 {
		msg := qq.CreateReply(source.Message)
		msg.Append(qq.NewTextfLn("找不到 %s 在赛季 %s 的记录。", args[0], args[1]))
		return qq.SendWithRandomRiskyStrategy(msg)
	}

	image := false
	if len(args) > 1 {
		image = args[1] == "image" || strings.HasPrefix(args[1], "图") || args[1] == "true"
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 在赛季 %s 的 MMR 资料如下:", args[0], args[1]))
	msg.Append(qq.NewTextfLn("最终段位: %s", data.FinalRankPatched))
	msg.Append(qq.NewTextfLn("总场数: %d", data.NumberOfGames))
	msg.Append(qq.NewTextfLn("总胜利次数: %d", data.Wins))
	msg.Append(qq.NewTextfLn("胜率: %.1f", formatPercentageInt(data.Wins, data.NumberOfGames)))
	msg.Append(qq.NewTextfLn("段位变更记录(每次胜利): "))

	for i, act := range data.ActRankWins {
		msg.Append(qq.NewTextfLn("\t%d. %s", i+1, act.PatchedTier))
	}

	if image {
		if err = qq.SendGroupImageText(msg); err == nil {
			return nil
		} else {
			logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
		}
	}

	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

// mmrActs GetMMRDetailsV2
func mmrActs(args []string, source *command.MessageSource) error {
	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}

	acts, err := valorant.GetMMRDetailsV2(name, tag, valorant.AsiaSpecific)
	if err != nil {
		return err
	}

	image := false
	if len(args) > 1 {
		image = args[1] == "image" || strings.HasPrefix(args[1], "图") || args[1] == "true"
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s#%s 的赛季段位资料如下: ", acts.Name, acts.Tag))
	msg.Append(qq.NewTextfLn("目前段位: %s", acts.CurrentData.CurrentTierPatched))
	msg.Append(qq.NewTextfLn("赛季段位:"))

	for _, season := range valorant.SortSeason(acts.BySeason) {

		data := acts.BySeason[season]

		if data.Error != "" {
			msg.Append(qq.NewTextfLn("	%s: 没有记录", season))
		} else {
			msg.Append(qq.NewTextfLn("	%s: %s", season, data.FinalRankPatched))
		}
	}

	if image {
		if err = qq.SendGroupImageText(msg); err == nil {
			return nil
		} else {
			logger.Errorf("發送圖片消息失敗: %v，将改回文本发送", err)
		}
	}

	return qq.SendWithRandomRiskyStrategyRemind(msg, source.Message)
}

func weapons(args []string, source *command.MessageSource) error {

	if !valorant.LangAvailable.Contains(args[0]) {
		msg := qq.CreateReply(source.Message)
		msg.Append(qq.NewTextf("未知语言，目前支援的语言: %s", strings.Join(valorant.LangAvailable.ToSlice(), ", ")))
		return qq.SendGroupMessage(msg)
	}

	presend := qq.CreateReply(source.Message)
	presend.Append(qq.NewTextf("正在索取瓦武器列表资料..."))
	_ = qq.SendGroupMessage(presend)

	weapons, err := valorant.GetWeapons(valorant.AllWeapons, valorant.Language(args[0]))
	if err != nil {
		return err
	}

	for _, weapon := range weapons {
		msg := message.NewSendingMessage()

		msg.Append(qq.NewTextfLn("武器名称: %s", weapon.DisplayName))
		msg.Append(qq.NewTextfLn("武器类型: %s", weapon.ShopData.CategoryText))
		msg.Append(qq.NewTextfLn("武器价格: $%d", weapon.ShopData.Cost))
		img, err := qq.NewImageByUrl(weapon.DisplayIcon)
		if err != nil {
			logger.Errorf("获取武器 %s 图片时出现错误: %v", weapon.DisplayName, err)
			msg.Append(qq.NewTextfLn("[图片]"))
		} else {
			msg.Append(img)
		}

		_ = qq.SendWithRandomRiskyStrategy(msg)
	}

	return nil
}

func agents(args []string, source *command.MessageSource) error {

	if !valorant.LangAvailable.Contains(args[0]) {
		msg := qq.CreateReply(source.Message)
		msg.Append(qq.NewTextf("未知语言，目前支援的语言: %s", strings.Join(valorant.LangAvailable.ToSlice(), ", ")))
		return qq.SendGroupMessage(msg)
	}

	presend := qq.CreateReply(source.Message)
	presend.Append(qq.NewTextf("正在索取瓦角色列表资料..."))
	_ = qq.SendGroupMessage(presend)

	agents, err := valorant.GetAgents(valorant.AllAgents, valorant.Language(args[0]))
	if err != nil {
		return err
	}

	for _, agent := range agents {

		msg := message.NewSendingMessage()
		msg.Append(qq.NewTextfLn("角色名称: %s", agent.DisplayName))
		msg.Append(qq.NewTextfLn("角色类型: %s", agent.Role.DisplayName))
		msg.Append(qq.NewTextfLn("简介: %s", agent.Description))
		if agent.CharacterTags != nil {
			msg.Append(qq.NewTextfLn("标签: %s", strings.Join(*agent.CharacterTags, ", ")))
		}

		skills := make([]string, 0)

		for _, skill := range agent.Abilities {
			skills = append(skills, skill.DisplayName)
		}

		msg.Append(qq.NewTextfLn("技能: %s", strings.Join(skills, ", ")))
		img, err := qq.NewImageByUrl(agent.KillfeedPortrait)
		if err != nil {
			logger.Errorf("获取角色 %s 图片时出现错误: %v", agent.DisplayName, err)
			msg.Append(qq.NewTextfLn("[图片]"))
		} else {
			msg.Append(img)
		}

		_ = qq.SendWithRandomRiskyStrategy(msg)

	}

	return nil
}

var (
	infoCommand         = command.NewNode([]string{"info", "资讯"}, "查询玩家账户资讯", false, info, "<名称#Tag>")
	forceUpdateCommand  = command.NewNode([]string{"update", "更新"}, "强制更新玩家资讯", false, forceUpdate, "<名称#Tag>")
	statusCommand       = command.NewNode([]string{"status", "状态"}, "查询状态", false, status)
	trackCommand        = command.NewNode([]string{"track", "追踪玩家"}, "追踪玩家最新对战", true, track, "<名称#Tag>")
	untrackCommand      = command.NewNode([]string{"untrack", "取消追踪玩家"}, "取消追踪玩家最新对战", true, untrack, "<名称#Tag>")
	trackingCommand     = command.NewNode([]string{"tracking", "追踪中"}, "查询追踪中的玩家", false, tracking)
	matchesCommand      = command.NewNode([]string{"matches", "对战历史"}, "查询对战历史", false, matches, "<名称#Tag>")
	matchCommand        = command.NewNode([]string{"match", "对战"}, "查询对战详情", false, match, "<对战ID>")
	leaderboardCommand  = command.NewNode([]string{"leaderboard", "排行榜"}, "查询对战排行榜", false, leaderboard, "<对战ID>", "[图片]")
	performanceCommand  = command.NewNode([]string{"performance", "表现", "击杀表现"}, "查询对战玩家的击杀表现", false, performances, "<对战ID>", "<名称#Tag>", "[图片]")
	statsCommand        = command.NewNode([]string{"stats", "统计数据"}, "查询该玩家在最近五场的统计数据", false, stats, "<名称#Tag>", "[对战模式]", "[图片]")
	matchPlayerscommand = command.NewNode([]string{"players", "玩家"}, "查询对战玩家资讯", false, matchPlayers, "<对战ID>")
	matchRoundsCommand  = command.NewNode([]string{"rounds", "回合"}, "查询对战回合资讯", false, matchRounds, "<对战ID>")
	mmrCommand          = command.NewNode([]string{"mmr", "段位"}, "查询段位", false, mmr, "<名称#Tag>")
	mmrHistoriesCommand = command.NewNode([]string{"mmrhist", "段位历史"}, "查询段位历史", false, mmrHistories, "<名称#Tag>")
	mmrBySeasonCommand  = command.NewNode([]string{"season", "赛季段位"}, "查询赛季段位", false, mmrBySeason, "<名称#Tag>", "<赛季>")
	mmrActsCommand      = command.NewNode([]string{"mmracts", "赛季段位历史"}, "查询赛季段位历史", false, mmrActs, "<名称#Tag>")
	weaponsCommand      = command.NewNode([]string{"weapons", "武器", "武器列表"}, "查询武器名称", false, weapons, "<语言区域>")
	agentCommand        = command.NewNode([]string{"agents", "角色", "特务", "角色列表", "特务列表"}, "查询角色名称", false, agents, "<语言区域>")
)

var valorantCommand = command.NewParent([]string{"valorant", "val", "瓦罗兰", "瓦"}, "valorant指令",
	infoCommand,
	forceUpdateCommand,
	statusCommand,
	trackCommand,
	untrackCommand,
	trackingCommand,
	matchesCommand,
	matchCommand,
	leaderboardCommand,
	performanceCommand,
	statsCommand,
	matchPlayerscommand,
	matchRoundsCommand,
	mmrCommand,
	mmrHistoriesCommand,
	mmrBySeasonCommand,
	mmrActsCommand,
	weaponsCommand,
	agentCommand,
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
	default:
		red := data.Teams["red"]
		blue := data.Teams["blue"]
		return fmt.Sprintf("Red %d : %d Blue", red.RoundsWon, blue.RoundsWon)
	}
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
	default:
		red := data.Teams["red"]
		blue := data.Teams["blue"]
		team, err := valorant.FoundPlayerInTeam(name, &data)
		if err != nil {
			return fmt.Sprintf("(错误: %s)", err.Error())
		}
		return fmt.Sprintf("Red %d : %d Blue (用户所在队伍: %s)", red.RoundsWon, blue.RoundsWon, team)
	}
}

func formatPercentage(part, total int64) float64 {
	return float64(part) / float64(total) * 100
}

func formatPercentageInt(part, total int) float64 {
	return float64(part) / float64(total) * 100
}

func formatTime(timeStr string) string {
	if timeStr == "" {
		return "无"
	}
	ti, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		logger.Errorf("无法解析时间: %s, 将返回厡讯息", timeStr)
		return timeStr
	}
	return ti.Format(datetime.TimeFormat)
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

func generateMatchPlayersLines(match *valorant.MatchData) *message.SendingMessage {

	ffInfo := valorant.GetFriendlyFireInfo(match)
	ranking := valorant.GetMatchRanking(match)

	msg := message.NewSendingMessage()
	for i, player := range ranking {
		msg.Append(qq.NewTextfLn("\t第 %d 名: %s", i+1, fmt.Sprintf("%s#%s", player.Name, player.Tag)))

		// 基本资料
		msg.Append(qq.NewTextLn("\t基本资料:"))
		msg.Append(qq.NewTextfLn("\t\tK/D/A: %d/%d/%d (%.2f)", player.Stats.Kills, player.Stats.Deaths, player.Stats.Assists, float64(player.Stats.Kills)/float64(player.Stats.Deaths)))
		msg.Append(qq.NewTextfLn("\t\t分数: %d", player.Stats.Score))
		msg.Append(qq.NewTextfLn("\t\t使用角色: %s", player.Character))

		// 如果不是死鬥模式，则显示所在队伍
		if strings.ToLower(match.MetaData.Mode) != "deathmatch" {
			msg.Append(qq.NewTextfLn("\t\t所在队伍: %s", player.Team))
		}

		// 如果是競技模式，則顯示段位
		if strings.ToLower(match.MetaData.Mode) == "competitive" {
			msg.Append(qq.NewTextfLn("\t\t段位: %s", player.CurrentTierPatched))
		}

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
		msg.Append(qq.NewTextfLn("\t\t误击队友伤害: %.f", friendlyFire.Outgoing))
		msg.Append(qq.NewTextfLn("\t\t误杀队友次数: %d", friendlyFire.Kills))
		msg.Append(qq.NewTextfLn("\t\t被误击队友伤害: %.f", friendlyFire.Incoming))
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

	return msg
}

func getShortIdHint(uuid string) (string, int64) {
	shortHint := ""
	short, err := valorant.ShortenUUID(uuid)
	if err != nil {
		logger.Warnf("无法缩短 UUID: %v", err)
	} else {
		shortHint = fmt.Sprintf(" (短号: %d)", short)
	}
	return shortHint, short
}

func getShortIdsHint(uuids []string) map[string]string {
	shortHints := make(map[string]string)
	shorts, errs := valorant.ShortenUUIDs(uuids)
	if len(errs) > 0 {
		for uuid, err := range errs {
			logger.Warnf("无法缩短 UUID %s: %v", uuid, err)
		}
	} else {
		for uuid, short := range shorts {
			shortHints[uuid] = fmt.Sprintf(" (短号: %d)", short)
		}
	}
	return shortHints
}
