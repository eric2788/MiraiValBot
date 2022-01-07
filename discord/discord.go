package discord

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/file"
	"strconv"
)

var (
	client *discordgo.Session
	config file.DiscordConfig
	logger = utils.GetModuleLogger("discord.bot")
)

func Start() {
	config = file.ApplicationYaml.Discord
	discord, err := discordgo.New("Bot " + file.ApplicationYaml.Discord.Token)
	if err != nil {
		logger.Errorf("啟動 discord 機器人時出現錯誤: %v\n", err)
		return
	}
	client = discord
	Log("Discord 機器人已成功啟動。")
}

func Log(msg string, arg ...interface{}) {
	if client == nil {
		logger.Warnf("Discord 尚未啟動，無法發送 Log")
		return
	}
	line := msg
	if len(arg) > 0 {
		line = fmt.Sprintf(msg, arg)
	}
	logger.Infof("發送 Discord Log 訊息 => %s\n", line)
	_, err := client.ChannelMessageSend(strconv.FormatInt(config.LogChannel, 10), line)
	if err != nil {
		logger.Warnf("發送 Discord Log 訊息時出現錯誤: %v\n", err)
	}
}

func RunSafe(runner func(*discordgo.Session)) {
	if client == nil {
		logger.Warnf("Discord 尚未啟動，無法進行操作")
		return
	}
	runner(client)
}

func GoRunSafe(runner func(*discordgo.Session)) {
	if client == nil {
		logger.Warnf("Discord 尚未啟動，無法進行操作")
		return
	}
	go runner(client)
}
