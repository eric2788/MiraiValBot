package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	ApplicationCommand() *discordgo.ApplicationCommand
	Handler(session *discordgo.Session, interact *discordgo.InteractionCreate) (err error)
}

var commandHandler = make(map[string]Command)

func RegisterCommand(c Command) {
	commandHandler[c.ApplicationCommand().Name] = c
}

func HookCommands() {
	for name, cmd := range commandHandler {
		_, err := client.ApplicationCommandCreate(client.State.User.ID, fmt.Sprint(config.Guild), cmd.ApplicationCommand())
		if err != nil {
			logger.Errorf("注册指令 %s 失败: %v", name, err)
			continue
		}
	}
	client.AddHandler(handleCommands)
}

func handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if cmd, ok := commandHandler[i.ApplicationCommandData().Name]; ok {
		if err := cmd.Handler(s, i); err != nil {
			logger.Errorf("执行指令 %s 失败: %v", cmd.ApplicationCommand().Name, err)
		}
	}
}

func ToOptionMap(data discordgo.ApplicationCommandInteractionData) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	options := data.Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}
	return optionMap
}
