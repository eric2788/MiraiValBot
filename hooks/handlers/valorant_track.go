package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/hooks/sites/valorant"
	v "github.com/eric2788/MiraiValBot/services/valorant"
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

	return valorantTrackRisky(displayName, shortHint, cmdId, metaData)
}

func init() {
	valorant.MessageHandler.AddHandler(valorant.MatchesUpdated, OnMatchesUpdated)
}
