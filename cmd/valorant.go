package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/valorant"
	"github.com/eric2788/common-utils/datetime"
)

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
	name, tag, err := valorant.ParseNameTag(args[0])
	if err != nil {
		return err
	}
	matches, err := valorant.GetMatchHistories(name, tag, valorant.AsiaSpecific)
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 最近的比赛:", args[0]))
	for _, match := range matches {
		// empty match id
		if match.MetaData.MatchId == "" {
			continue
		}
		msg.Append(qq.NewTextLn("===================="))
		msg.Append(qq.NewTextfLn("比赛ID: %s", match.MetaData.MatchId))
		msg.Append(qq.NewTextfLn("比赛模式: %s", match.MetaData.Mode))
		msg.Append(qq.NewTextfLn("比赛开始时间: %s", datetime.FormatSeconds(match.MetaData.GameStart)))
		msg.Append(qq.NewTextfLn("比赛时长: %s", formatDuration(match.MetaData.GameLength)))
		msg.Append(qq.NewTextfLn("比赛地图: %s", match.MetaData.Map))
		msg.Append(qq.NewTextfLn("回合总数: %d", match.MetaData.RoundsPlayed))
		msg.Append(qq.NewTextfLn("服务器: %s", match.MetaData.Cluster))
		msg.Append(qq.NewTextfLn("比赛结果: %s", formatResult(match, args[0])))
		msg.Append(qq.NewTextfLn("输入 !val players %s 查看详细玩家信息", match.MetaData.MatchId))
		msg.Append(qq.NewTextfLn("输入 !val rounds %s 查看详细回合信息", match.MetaData.MatchId))
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

func match(args []string, source *command.MessageSource) error {
	match, err := valorant.GetMatchDetails(args[0])
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("比赛ID: %s", match.MetaData.MatchId))
	msg.Append(qq.NewTextfLn("比赛模式: %s", match.MetaData.Mode))
	msg.Append(qq.NewTextfLn("比赛开始时间: %s", datetime.FormatSeconds(match.MetaData.GameStart)))
	msg.Append(qq.NewTextfLn("比赛时长: %s", formatDuration(match.MetaData.GameLength)))
	msg.Append(qq.NewTextfLn("比赛地图: %s", match.MetaData.Map))
	msg.Append(qq.NewTextfLn("回合总数: %d", match.MetaData.RoundsPlayed))
	msg.Append(qq.NewTextfLn("服务器: %s", match.MetaData.Cluster))
	msg.Append(qq.NewTextfLn("比赛结果: %s", formatResult(*match, args[0])))
	msg.Append(qq.NewTextfLn("输入 !val players %s 查看详细玩家信息", match.MetaData.MatchId))
	msg.Append(qq.NewTextfLn("输入 !val rounds %s 查看详细回合信息", match.MetaData.MatchId))
	return qq.SendWithRandomRiskyStrategy(msg)
}

func matchPlayers(args []string, source *command.MessageSource) error {

	match, err := valorant.GetMatchDetails(args[0])
	if err != nil {
		return err
	}

	ranking := valorant.GetMatchRanking(match)

	msg := message.NewSendingMessage()
	for i, player := range ranking {

		// 基本资料
		msg.Append(qq.NewTextfLn("%d. %s (K: %d, D: %d, A: %d, 分数: %d, 使用角色: %s, 所在队伍: %s)",
			i+1,
			fmt.Sprintf("%s#%s", player.Name, player.Tag),
			player.Stats.Kills,
			player.Stats.Deaths,
			player.Stats.Assists,
			player.Stats.Score,
			player.Character,
			player.Team,
		))

		// 击中分布
		total := player.Stats.BodyShots + player.Stats.Headshots + player.Stats.LegShots
		msg.Append(qq.NewTextfLn("\t击中次数分布: 头 %.1f% (%d) 身 %.1f (%d) 腿 %.1f (%d)",
			float64(player.Stats.Headshots/total), player.Stats.Headshots,
			float64(player.Stats.BodyShots/total), player.Stats.BodyShots,
			float64(player.Stats.LegShots/total), player.Stats.LegShots,
		))

		// 行为
		msg.Append(qq.NewTextfLn("\t行为: AFK回合次数 %d, 误击队友伤害 %d, 被误击队友伤害 %d, 拆包次数 %d, 装包次数 %d",
			player.Behaviour.AfkRounds,
			player.Behaviour.FriendlyFire.Outgoing,
			player.Behaviour.FriendlyFire.Incoming,
			valorant.GetDefuseCount(match, args[0]),
			valorant.GetPlantCount(match, args[0]),
		))

		//技能使用
		total = 0
		for _, times := range player.AbilityCasts {
			total += times
		}
		msg.Append(qq.NewTextfLn("\t技能使用次数分布: C %.1f (%d) Q %.1f (%d) E %.1f (%d) X %.1f (%d)",
			float64(player.AbilityCasts["c_cast"]/total), player.AbilityCasts["c_cast"],
			float64(player.AbilityCasts["q_cast"]/total), player.AbilityCasts["q_cast"],
			float64(player.AbilityCasts["e_cast"]/total), player.AbilityCasts["e_cast"],
			float64(player.AbilityCasts["x_cast"]/total), player.AbilityCasts["x_cast"],
		))

		// 经济
		msg.Append(qq.NewTextfLn("\t经济: 总支出 $%d, 平均支出 $%d", player.Economy.Spent.OverAll, player.Economy.Spent.Average))

		// 伤害
		totalDamage := player.DamageReceived + player.DamageMade
		msg.Append(qq.NewTextfLn("\t伤害分布: 总承受 %.1f (%d) 总伤害 %.1f (%d)",
			float64(player.DamageReceived/totalDamage), player.DamageReceived,
			float64(player.DamageMade/totalDamage), player.DamageMade,
		))
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

func matchRounds(args []string, source *command.MessageSource) error {
	return qq.SendGroupMessage(message.NewSendingMessage().Append(qq.NewTextLn("此指令暂不可用")))
}

// mmr get MMRV1Details
func mmr(args []string, source *command.MessageSource) error {
	parts := strings.Split(args[0], "#")
	mmr, err := valorant.GetMMRDetailsV1(parts[0], parts[1], valorant.AsiaSpecific)
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("======== %s 的 MMR 资料 ======="))
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
	statusCommand       = command.NewNode([]string{"status", "状态"}, "查询状态", false, status)
	matchesCommand      = command.NewNode([]string{"matches", "比赛历史"}, "查询比赛历史", false, matches)
	matchCommand        = command.NewNode([]string{"match", "比赛"}, "查询比赛详情", false, match, "<比赛ID>")
	matchPlayerscommand = command.NewNode([]string{"players", "玩家"}, "查询比赛玩家资讯", false, matchPlayers, "<比赛ID>")
	matchRoundsCommand  = command.NewNode([]string{"rounds", "回合"}, "查询比赛回合资讯", false, matchRounds, "<比赛ID>")
	mmrCommand          = command.NewNode([]string{"mmr", "段位"}, "查询段位", false, mmr, "<名称#Tag>")
	mmrHistoriesCommand = command.NewNode([]string{"mmrHistories", "段位历史"}, "查询段位历史", false, mmrHistories, "<名称#Tag>")
	mmrBySeasonCommand  = command.NewNode([]string{"mmrBySeason", "赛季段位"}, "查询赛季段位", false, mmrBySeason, "<名称#Tag>", "<赛季>")
	localizeCommand     = command.NewNode([]string{"localize", "本地化"}, "更新i18n内容", true, localize)
)

var valorantCommand = command.NewParent([]string{"valorant", "val", "瓦罗兰", "瓦"}, "valorant指令",
	statusCommand,
	matchesCommand,
	matchCommand,
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
	case "competitive":
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
