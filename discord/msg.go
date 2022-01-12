package discord

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strconv"
	"time"
)

type ChannelType uint8

func SendNewsEmbed(embed *discordgo.MessageEmbed) {
	SendEmbed(config.NewsChannel, embed)
}

func SendNewsTxt(txt string) {
	SendText(config.NewsChannel, txt)
}

func SendLogText(txt string) {
	SendText(config.LogChannel, txt)
}

func SendLogEmbed(embed *discordgo.MessageEmbed) {
	SendEmbed(config.LogChannel, embed)
}

func SendText(channel int64, content string) {
	RunSafe(func(session *discordgo.Session) (err error) {
		_, err = session.ChannelMessageSend(strconv.FormatInt(channel, 10), content)
		return
	})
}

func randomColor() int {
	rand.Seed(time.Now().UnixMicro())
	return rand.Intn(16777216)
}

func SendEmbed(channel int64, embed *discordgo.MessageEmbed) {
	embed.Color = randomColor() // 隨機顏色
	RunSafe(func(session *discordgo.Session) (err error) {
		_, err = session.ChannelMessageSendEmbed(strconv.FormatInt(channel, 10), embed)
		return
	})
}
