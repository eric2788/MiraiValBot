package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/file"
	"strconv"
)

var (
	Client *discordgo.Session
	config file.DiscordConfig
)

func Start() {
	config = file.ApplicationYaml.Discord
	discord, err := discordgo.New("Bot " + file.ApplicationYaml.Discord.Token)
	if err != nil {
		fmt.Printf("啟動 discord 機器人時出現錯誤: %v\n", err)
		return
	}
	Client = discord
	Log("Discord 機器人已成功啟動。")
}

func Log(msg string, arg ...interface{}) {
	line := fmt.Sprintf(msg, arg)
	fmt.Printf("發送 Discord Log 訊息 => %s\n", line)
	_, err := Client.ChannelMessageSend(strconv.FormatInt(config.LogChannel, 10), line)
	if err != nil {
		fmt.Printf("發送 Discord Log 訊息時出現錯誤: %v\n", err)
	}
}
