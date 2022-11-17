package discord

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/services/waifu"
)

// type commandHandler func(s *discordgo.Session, i *discordgo.InteractionCreate)

func HookCommand() {
	one := float64(1)
	cmd := &discordgo.ApplicationCommand{
		Name:        "色图",
		Description: "来点色图",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "多少张",
				Description: "色图数量, 每次最多40张",
				Required:    true,
				MinValue:    &one,
				MaxValue:    40,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "是否R18",
				Description: "是否发送R18色图",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "标签",
				Description: "选填, 用,分隔",
				Required:    false,
			},
		},
	}

	_, err := client.ApplicationCommandCreate(client.State.User.ID, fmt.Sprint(config.Guild), cmd)
	if err != nil {
		logger.Errorf("注册色图指令失败: %v", err)
	}

	client.AddHandler(handleSetu)
}

func handleSetu(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	if data.Name != "色图" {
		return
	}
	if config.NsfwChannel == 0 {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "未设置 NSFW 频道ID, 因此无法发送色图!",
		})
		return
	}
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "正在索取色图....",
		},
	})

	if err != nil {
		logger.Errorf("发送回应失败: %v", err)
	}

	options := data.Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	amount := optionMap["多少张"].IntValue()
	r18 := optionMap["是否R18"].BoolValue()
	tags := []string{""}
	if opt, ok := optionMap["标签"]; ok {
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
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("获取色图失败: %v", err),
		})
		return
	}

	content := fmt.Sprintf("正在发送色图到频道 <#%d> ....", config.NsfwChannel)
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		logger.Errorf("更改回应失败: %v", err)
	}

	for _, img := range imgs {
		SendNSFWImage(img)
	}

	content = "发送成功。"
	_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
	if err != nil {
		logger.Errorf("更改回应失败: %v", err)
	}
}
