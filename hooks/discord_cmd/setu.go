package discord_cmd

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/discord"
	"github.com/eric2788/MiraiValBot/services/waifu"
)

type setu struct {
}

func (s *setu) ApplicationCommand() *discordgo.ApplicationCommand {
	one := float64(1)
	return &discordgo.ApplicationCommand{
		Name: "setu",
		NameLocalizations: &map[discordgo.Locale]string{
			discordgo.ChineseTW: "色圖",
			discordgo.ChineseCN: "色图",
		},
		Description: "来点色图",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type: discordgo.ApplicationCommandOptionInteger,
				Name: "amount",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.ChineseTW: "多少張",
					discordgo.ChineseCN: "多少张",
				},
				Description: "色图数量, 每次最多40张",
				Required:    true,
				MinValue:    &one,
				MaxValue:    40,
			},
			{
				Type: discordgo.ApplicationCommandOptionBoolean,
				Name: "nsfw",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.ChineseTW: "涩涩",
					discordgo.ChineseCN: "涩涩",
				},
				Description: "是否发送R18色图",
				Required:    true,
			},
			{
				Type: discordgo.ApplicationCommandOptionString,
				Name: "tags",
				NameLocalizations: map[discordgo.Locale]string{
					discordgo.ChineseTW: "標籤",
					discordgo.ChineseCN: "标签",
				},
				Description: "选填, 用,分隔",
				Required:    false,
			},
		},
	}
}

func (s *setu) Handler(session *discordgo.Session, interact *discordgo.InteractionCreate) (err error) {
	var config = file.ApplicationYaml.Discord
	data := interact.ApplicationCommandData()
	if config.NsfwChannel == 0 {
		_, err = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
			Content: "未设置 NSFW 频道ID, 因此无法发送色图!",
		})
		return
	}
	err = session.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "正在索取色图....",
		},
	})

	if err != nil {
		return
	}

	optionMap := discord.ToOptionMap(data)

	amount := optionMap["amount"].IntValue()
	r18 := optionMap["nsfw"].BoolValue()
	tags := []string{""}
	if opt, ok := optionMap["tags"]; ok {
		tags = strings.Split(opt.StringValue(), ",")
	}

	var search waifu.Searcher
	if len(tags) == 1 {
		search = waifu.WithKeyword(tags[0])
	} else {
		search = waifu.WithTags(tags...)
	}

	imgs, err := waifu.GetRandomImages(
		waifu.NewOptions(
			search,
			waifu.WithAmount(int(amount)),
			waifu.WithR18(r18),
		),
	)

	if err != nil {
		_, _ = session.FollowupMessageCreate(interact.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("获取色图失败: %v", err),
		})
		return
	}

	content := fmt.Sprintf("正在发送色图到频道 <#%d> ....", config.NsfwChannel)
	_, _ = session.InteractionResponseEdit(interact.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	for _, img := range imgs {
		sendPixivImage(img)
	}
	content = fmt.Sprintf("已成功发送到频道 <#%d>。", config.NsfwChannel)
	_, err = session.InteractionResponseEdit(interact.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	return
}

func sendPixivImage(data *waifu.ImageData) {
	discord.SendEmbed(file.ApplicationYaml.Discord.NsfwChannel, &discordgo.MessageEmbed{
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

func init() {
	discord.RegisterCommand(&setu{})
}
