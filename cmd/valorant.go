package cmd

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/valorant"
	"github.com/eric2788/common-utils/datetime"
	"strings"
	"time"
)

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
	return fmt.Sprintf("%.2f 分钟", min)
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
	parts := strings.Split(args[0], "#")
	matches, err := valorant.GetMatchHistories(parts[0], parts[1], valorant.AsiaSpecific)
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 最近的比赛:", args[0]))
	for _, match := range matches {
		msg.Append(qq.NewTextLn("===================="))
		msg.Append(qq.NewTextfLn("比赛ID: %s", match.MetaData.MatchId))
		msg.Append(qq.NewTextfLn("比赛模式: %s", match.MetaData.Mode))
		msg.Append(qq.NewTextfLn("比赛开始时间: %s", datetime.FormatMillis(match.MetaData.GameStart)))
		msg.Append(qq.NewTextfLn("比赛时长: %s", formatDuration(match.MetaData.GameLength)))
		msg.Append(qq.NewTextfLn("比赛地图: %s", match.MetaData.Map))
		msg.Append(qq.NewTextfLn("回合总数: %s", match.MetaData.RoundsPlayed))
		msg.Append(qq.NewTextfLn("服务器: %s", match.MetaData.Cluster))
		msg.Append(qq.NewTextfLn("输入 /valorant match %s 查看详细信息", match.MetaData.MatchId))
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

func match(args []string, source *command.MessageSource) error {
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
	msg.Append(qq.NewTextfLn("%s 的目前段位: %s", mmr.Name, mmr.CurrentTierPatched))
	msg.Append(qq.NewTextfLn("%s 的目前段位分数: %d/100", mmr.Name, mmr.RankingInTier))
	msg.Append(qq.NewTextfLn("%s 上一次的分数变更: %d", mmr.Name, mmr.MMRChangeToLastGame))
	msg.Append(qq.NewTextfLn("%s 的 ELO: %d", mmr.Name, mmr.Elo))
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
	mmrCommand          = command.NewNode([]string{"mmr", "段位"}, "查询段位", false, mmr, "<名称#Tag>")
	mmrHistoriesCommand = command.NewNode([]string{"mmrHistories", "段位历史"}, "查询段位历史", false, mmrHistories, "<名称#Tag>")
	mmrBySeasonCommand  = command.NewNode([]string{"mmrBySeason", "赛季段位"}, "查询赛季段位", false, mmrBySeason, "<名称#Tag>", "<赛季>")
	localizeCommand     = command.NewNode([]string{"localize", "本地化"}, "查询本地化字段", false, localize, "<英文字段>")
)

var valorantCommand = command.NewParent([]string{"valorant", "valorant", "瓦罗兰", "瓦"}, "valorant指令",
	statusCommand,
	matchesCommand,
	matchCommand,
	mmrCommand,
	mmrHistoriesCommand,
	mmrBySeasonCommand,
	localizeCommand,
)

func init() {
	command.AddCommand(valorantCommand)
}
