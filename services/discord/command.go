package discord

import (
	"fmt"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	ApplicationCommand() *discordgo.ApplicationCommand
	Handler(session *discordgo.Session, interact *discordgo.InteractionCreate) (err error)
}

var (
	commandHandler = make(map[string]Command)
)

func RegisterCommand(c Command) {
	commandHandler[c.ApplicationCommand().Name] = c
}

func HookCommands() {
	// delete all commands before hook
	// UnRegisterCommands()
	for name, cmd := range commandHandler {
		_, err := client.ApplicationCommandCreate(client.State.User.ID, fmt.Sprint(config.Guild), cmd.ApplicationCommand())
		if err != nil {
			logger.Errorf("注册指令 %s 失败: %v", name, err)
			continue
		} else {
			logger.Infof("注册 Discord 指令 %s 成功", name)
		}
	}
	client.AddHandler(handleCommands)
}

func UnRegisterCommands() {
	commands, err := client.ApplicationCommands(client.State.User.ID, fmt.Sprint(config.Guild))
	if err != nil {
		logger.Errorf("获取指令列表失败: %v", err)
		return
	}
	for _, cmd := range commands {
		err = client.ApplicationCommandDelete(client.State.User.ID, fmt.Sprint(config.Guild), cmd.ID)
		if err != nil {
			logger.Errorf("删除指令 %s 失败: %v", cmd.Name, err)
		} else {
			logger.Infof("删除指令 %s 成功", cmd.Name)
		}
	}
}

func handleCommands(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if cmd, ok := commandHandler[i.ApplicationCommandData().Name]; ok {

		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("執行指令 %s 出現嚴重錯誤: %v", i.ApplicationCommandData().Name, err)
				debug.PrintStack()
			}
		}()

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
