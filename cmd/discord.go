package cmd

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
	"strings"
)

func sendToDiscord(args []string, source *command.MessageSource) error {

	sender := source.Message.Sender
	content := qq.ParseMsgContent(source.Message.Elements)
	text := strings.Join(args, " ")

	embed := &discordgo.MessageEmbed{
		Title:       "有从 QQ 传来的新消息",
		Description: text,
		Author: &discordgo.MessageEmbedAuthor{
			Name:    sender.Nickname,
			IconURL: fmt.Sprintf("https://q.qlogo.cn/g?b=qq&s=640&nk=%v", sender.Uin),
		},
	}

	if len(content.Images) == 1 {
		embed.Image = &discordgo.MessageEmbedImage{URL: content.Images[0]}
	} else if len(content.Images) > 1 {
		embed.Fields = make([]*discordgo.MessageEmbedField, len(content.Images))
		for i, img := range content.Images {
			embed.Fields[i] = &discordgo.MessageEmbedField{
				Name:  fmt.Sprintf("图片%v", i+1),
				Value: img,
			}
		}
	}

	go discord.SendEmbed(file.ApplicationYaml.Discord.CrossPlatChannel, embed)
	return nil

}

var discordCommand = command.NewNode([]string{"discord", "dc", "跨平台"}, "往 discord 发送讯息", false, sendToDiscord, "<讯息>")

func init() {
	command.AddCommand(discordCommand)
}
