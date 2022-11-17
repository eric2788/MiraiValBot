package discord

import (
	"fmt"
	"strconv"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/internal/file"
)

var (
	client *discordgo.Session
	config *file.DiscordConfig
	logger = utils.GetModuleLogger("discord.bot")
)

func Start() {
	config = &file.ApplicationYaml.Discord
	discord, err := discordgo.New("Bot " + file.ApplicationYaml.Discord.Token)
	if err != nil {
		logger.Errorf("啟動 discord 機器人時出現錯誤: %v\n", err)
		return
	}
	client = discord
	client.Identify.Intents = discordgo.IntentsGuildMessages
	go StartChatListen() // 啟動跨平台聊天
	Log("Discord 機器人已成功啟動。")
}

func Close() error {
	if client != nil {
		return client.Close()
	}
	return nil
}

func Log(msg string, arg ...interface{}) {
	if client == nil {
		logger.Warnf("Discord 尚未啟動，無法發送 Log")
		return
	}
	line := msg
	if len(arg) > 0 {
		line = fmt.Sprintf(msg, arg...)
	}
	logger.Infof("發送 Discord Log 訊息 => %s\n", line)
	_, err := client.ChannelMessageSend(strconv.FormatInt(config.LogChannel, 10), line)
	if err != nil {
		logger.Warnf("發送 Discord Log 訊息時出現錯誤: %v\n", err)
	}
}

func RunSafe(runner func(*discordgo.Session) error) {
	if client == nil {
		logger.Warnf("Discord 尚未啟動，無法進行操作")
		return
	}
	if err := runner(client); err != nil {
		logger.Warnf("Discord 發送訊息時出現錯誤: %v", err)
	}
}
