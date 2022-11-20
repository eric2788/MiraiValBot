package cmd

import (
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
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

var (
	fetchEssenceCommand = command.NewNode([]string{"fetchess"}, "刷新群精华消息到快取", true, fetchEssence)
	migrateCacheCommand = command.NewNode([]string{"migrate", "搬迁"}, "搬迁缓存", true, migrateCache, "<从缓存类型>", "<到缓存类型>", "<缓存路径>", "[是否移除旧资料]")
)

var adminCommand = command.NewParent([]string{"admin", "管理员", "管理"}, "管理员指令",
	fetchEssenceCommand,
	migrateCacheCommand,
)

func init() {
	command.AddCommand(adminCommand)
}
