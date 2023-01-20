package cmd

import (
	"fmt"
	"strconv"

	"github.com/Mrs4s/MiraiGo/message"
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/aichat"
	"github.com/eric2788/MiraiValBot/utils/cache"
)

func migrateCache(args []string, source *command.MessageSource) error {
	from, to, path := args[0], args[1], args[2]
	fromCache, err := cache.New(cache.WithPath(path), cache.WithType(from))
	if err != nil {
		return err
	}
	toCache, err := cache.New(cache.WithPath(path), cache.WithType(to))
	if err != nil {
		return err
	}
	remove := false
	if len(args) > 3 {
		remove = args[3] == "true"
	}

	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextf("正在从 %s 缓存 搬迁到 %s 缓存...", from, to))

	_ = qq.SendGroupMessage(msg)

	cache.Migrate(fromCache, toCache, remove).Wait()
	msg = qq.CreateReply(source.Message)
	msg.Append(qq.NewTextf("搬迁完成。请到控制台查看日志输出。"))

	return qq.SendGroupMessage(msg)
}

func fetchEssence(args []string, source *command.MessageSource) error {
	i, err := qq.FetchEssenceListToCache()
	if err != nil {
		return err
	}

	return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(qq.NewTextf("已成功添加 %d 则群精华消息到缓存。", i)))
}

func resetConversation(args []string, source *command.MessageSource) error {
	aichat.ResetGPTConversation()
	return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("已重置 GPT3 对话记录。")))
}

func atUser(args []string, source *command.MessageSource) error {
	times, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) == 0 {
		return fmt.Errorf("没有找到@的用户")
	}
	msg := message.NewSendingMessage()
	at := ats[0]
	for i := 0; i < times; i++ {
		msg.Append(at)
	}
	return qq.SendGroupMessage(msg)
}

var (
	fetchEssenceCommand = command.NewNode([]string{"fetchess"}, "刷新群精华消息到快取", true, fetchEssence)
	migrateCacheCommand = command.NewNode([]string{"migrate", "搬迁"}, "搬迁缓存", true, migrateCache, "<从缓存类型>", "<到缓存类型>", "<缓存路径>", "[是否移除旧资料]")
	resetConverCommand  = command.NewNode([]string{"resetchat", "重置对话"}, "重置gpt3对话记录", true, resetConversation)
	atUserCommand       = command.NewNode([]string{"at", "艾特"}, "艾特用户特定次数", true, atUser, "<次数> [@用户]")
)

var adminCommand = command.NewParent([]string{"admin", "管理员", "管理"}, "管理员指令",
	fetchEssenceCommand,
	migrateCacheCommand,
	resetConverCommand,
	atUserCommand,
)

func init() {
	command.AddCommand(adminCommand)
}
