package discord

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/services/waifu"
)

type ChannelType uint8

func SendNewsEmbed(embed *discordgo.MessageEmbed) {
	SendEmbed(config.NewsChannel, embed)
}

func SendNewsEmbedDouble(first, next *discordgo.MessageEmbed) {
	RunSafe(func(session *discordgo.Session) (err error) {
		news := strconv.FormatInt(config.NewsChannel, 10)
		msg, err := session.ChannelMessageSendEmbed(news, first)
		if err != nil {
			return
		}
		_, err = session.ChannelMessageSendComplex(news, &discordgo.MessageSend{
			Reference: &discordgo.MessageReference{
				MessageID: msg.ID,
				ChannelID: msg.ChannelID,
				GuildID:   msg.GuildID,
			},
			Embed: next,
		})
		return
	})
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

func SendNSFWImage(data *waifu.ImageData) {
	SendEmbed(config.NsfwChannel, &discordgo.MessageEmbed{
		Title: data.Title,
		Author: &discordgo.MessageEmbedAuthor{
			Name: data.Author,
			URL:  fmt.Sprintf("https://pixiv.net/users/%d", data.Uid),
		},
		URL: fmt.Sprintf("https://pixiv.net/artworks/%d", data.Pid),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "标签",
				Value: strings.Join(data.Tags, ","),
			},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: data.Url,
		},
		Provider: &discordgo.MessageEmbedProvider{
			Name: data.Author,
			URL:  fmt.Sprintf("https://pixiv.net/users/%d", data.Uid),
		},
	})
}

func SendEmbed(channel int64, embed *discordgo.MessageEmbed) {
	embed.Color = randomColor() // 隨機顏色
	RunSafe(func(session *discordgo.Session) (err error) {
		_, err = session.ChannelMessageSendEmbed(strconv.FormatInt(channel, 10), embed)
		return
	})
}
