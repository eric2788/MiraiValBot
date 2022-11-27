package discord_cmd

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/ai"
	"github.com/eric2788/MiraiValBot/services/discord"
	"github.com/eric2788/common-utils/request"
	"net/http"
	"strings"
)

type novelAI struct {
}

func (n *novelAI) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name: "novelai",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseTW: "畫圖",
			discordgo.ChineseCN: "画图",
		},
		Description: "AI画图",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "description",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.ChineseTW: "描述",
					discordgo.ChineseCN: "描述",
				},
				Description: "描述图画的内容",
				Required:    true,
			},
			{
				Type: discordgo.ApplicationCommandOptionBoolean,
				Name: "nsfw",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.ChineseTW: "涩涩",
					discordgo.ChineseCN: "涩涩",
				},
				Description: "是否包含R18色图",
				Required:    true,
			},
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "badprompt",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.ChineseTW: "不良标签",
					discordgo.ChineseCN: "不良标签",
				},
				Description: "选填, 用,分隔, 过滤不良标签",
				Required:    false,
			},
		},
	}
}

func (n *novelAI) Handler(session *discordgo.Session, interact *discordgo.InteractionCreate) (err error) {
	var config = file.ApplicationYaml.Discord
	data := interact.ApplicationCommandData()

	optionMap := discord.ToOptionMap(data)

	description := optionMap["description"].StringValue()
	nsfw := optionMap["nsfw"].BoolValue()
	badPrompt := ""
	if bp, ok := optionMap["badprompt"]; ok {
		badPrompt = bp.StringValue()
	}

	if interact.ChannelID != fmt.Sprint(config.NsfwChannel) && nsfw {
		err = session.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "请在 <#" + fmt.Sprint(config.NsfwChannel) + "> 频道使用此指令",
			},
		})
		return
	}

	err = session.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "正在生成图像....",
		},
	})

	t := ai.WithNSFW
	if nsfw {
		t = ai.WithR18
	}

	if err != nil {
		return
	}

	img, err := ai.GetNovelAI8zywImage(
		ai.New8zywPayload(
			description,
			t,
			strings.Split(badPrompt, ",")...,
		),
	)

	if err != nil {
		_, _ = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
			Content: "生成失败: " + err.Error(),
		})
		return
	}

	b, err := request.GetBytesByUrl(img)

	if err != nil {
		_, _ = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
			Content: "图片读取失败: " + err.Error(),
		})
		return
	}

	_, err = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
		Files: []*discordgo.File{
			{
				Name:        img,
				ContentType: http.DetectContentType(b),
				Reader:      bytes.NewReader(b),
			},
		},
		Content: description,
	})
	return
}

func init() {
	discord.RegisterCommand(&novelAI{})
}
