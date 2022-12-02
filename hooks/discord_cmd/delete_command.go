package discord_cmd

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/services/discord"
)

type deleteCommand struct {
}

func (d *deleteCommand) ApplicationCommand() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "delete",
		Description: "delete command with id (admin use only)",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "id",
				Description: "command id",
				Required:    true,
			},
		},
	}
}

func (d *deleteCommand) Handler(session *discordgo.Session, interact *discordgo.InteractionCreate) (err error) {

	optMap := discord.ToOptionMap(interact.ApplicationCommandData())

	err = session.InteractionRespond(interact.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "正在刪除指令...",
		},
	})

	if err != nil {
		return
	}

	err = session.ApplicationCommandDelete(interact.AppID, interact.GuildID, optMap["id"].StringValue())

	if err != nil {
		content := fmt.Sprintf("刪除指令失敗: %v", err)
		_, err = session.InteractionResponseEdit(interact.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	} else {
		content := "刪除指令成功"
		_, err = session.InteractionResponseEdit(interact.Interaction, &discordgo.WebhookEdit{
			Content: &content,
		})
		return
	}
}

func init() {
	discord.RegisterCommand(&deleteCommand{})
}
