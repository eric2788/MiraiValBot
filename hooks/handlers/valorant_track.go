package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/hooks/sites/valorant"
	"github.com/eric2788/MiraiValBot/internal/qq"
	v "github.com/eric2788/MiraiValBot/services/valorant"
	"github.com/eric2788/common-utils/datetime"
)

func OnMatchesUpdated(_ *bot.Bot, data *valorant.MatchMetaDataSub) error {

	displayName, metaData := data.DisplayName, data.Data

	if metaData.MatchId == "" || len(metaData.MatchId) == 0 {
		logger.Warnf("收到空的對戰ID: %q, 已略過。", metaData.MatchId)
		return nil
	}

	shortHint := ""
	short, err := v.ShortenUUID(metaData.MatchId)
	if err != nil {
		logger.Warnf("无法缩短 UUID: %v", err)
	} else {
		shortHint = fmt.Sprintf(" (短号: %d)", short)
	}

	cmdId := metaData.MatchId

	if short > -1 {
		cmdId = fmt.Sprintf("%d", short)
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 的最新对战信息已更新👇", displayName))
	msg.Append(qq.NewTextfLn("对战ID: %s%s", metaData.MatchId, shortHint))
	msg.Append(qq.NewTextfLn("对战模式: %s", metaData.Mode))
	msg.Append(qq.NewTextfLn("对战开始时间: %s", datetime.FormatSeconds(metaData.GameStart)))
	msg.Append(qq.NewTextfLn("对战地图: %s", metaData.Map))
	msg.Append(qq.NewTextfLn("输入 !val match %s 查看更详细资讯。", cmdId))

	return qq.SendWithRandomRiskyStrategy(msg)
}

func init() {
	valorant.MessageHandler.AddHandler(valorant.MatchesUpdated, OnMatchesUpdated)
}
