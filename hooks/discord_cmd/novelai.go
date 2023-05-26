package discord_cmd

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/aidraw"
	"github.com/eric2788/MiraiValBot/services/discord"
	"github.com/eric2788/common-utils/request"
	"net/http"
	"time"
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
				// model with choices
				Type: discordgo.ApplicationCommandOptionString,
				Name: "model",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.ChineseTW: "模型",
					discordgo.ChineseCN: "模型",
				},
				Description: "选填, 模型名称, 默认为 anime",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "anime",
						Value: "anime",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.ChineseTW: "動漫",
							discordgo.ChineseCN: "动漫",
						},
					},
					{
						Name:  "real1",
						Value: "real1",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.ChineseTW: "寫實1",
							discordgo.ChineseCN: "写实1",
						},
					},
					{
						Name:  "real2",
						Value: "real2",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.ChineseTW: "寫實2",
							discordgo.ChineseCN: "写实2",
						},
					},
					{
						Name:  "real3",
						Value: "real3",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.ChineseTW: "寫實3",
							discordgo.ChineseCN: "写实3",
						},
					},
					{
						Name:  "real (Random)",
						Value: "real",
						NameLocalizations: map[discordgo.Locale]string{
							discordgo.ChineseTW: "寫實 (隨機)",
							discordgo.ChineseCN: "写实 (随机)",
						},
					},
				},
			},
		},
	}
}

func (n *novelAI) Handler(session *discordgo.Session, interact *discordgo.InteractionCreate) (err error) {
	var config = file.ApplicationYaml.Discord
	data := interact.ApplicationCommandData()

	optionMap := discord.ToOptionMap(data)

	description := optionMap["description"].StringValue()
	model := ""
	if bp, ok := optionMap["model"]; ok {
		model = bp.StringValue()
	}

	if interact.ChannelID != fmt.Sprint(config.NsfwChannel) {
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

	if err != nil {
		return
	}

	img, err := aidraw.Draw(aidraw.Payload{
		Prompt: description,
		Model:  model,
	})

	if err != nil {
		_, _ = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
			Content: "生成失败: " + err.Error(),
		})
		return
	}

	var b []byte
	if img.ImgUrl != "" {
		b, err = request.GetBytesByUrl(img.ImgUrl)
	} else if len(img.ImgData) > 0 {
		b = img.ImgData
	} else {
		err = fmt.Errorf("图片数据为空")
	}

	if err != nil {
		_, _ = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
			Content: "图片读取失败: " + err.Error(),
		})
		return
	}

	_, err = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
		Files: []*discordgo.File{
			{
				Name:        fmt.Sprintf("aidraw-%d.png", time.Now().Unix()),
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
