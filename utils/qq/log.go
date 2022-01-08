package qq

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/eventhook"
)

type log struct {
}

func (l *log) HookEvent(bot *bot.Bot) {
	bot.OnLog(func(qqClient *client.QQClient, event *client.LogEvent) {
		logger.WithField("type", event.Type).Debug(event.Message)
	})

	bot.OnDisconnected(func(cli *client.QQClient, e *client.ClientDisconnectedEvent) {
		logger.Warn("bot 已離線: ", e.Message)
	})

	bot.OnGroupMessage(func(cli *client.QQClient, msg *message.GroupMessage) {
		logger.Infof("%s (%d) 在群 %s 發送了消息: %s", msg.Sender.Nickname, msg.Sender.Uin, msg.GroupName, msg.ToString())
	})

	bot.OnPrivateMessage(func(cli *client.QQClient, msg *message.PrivateMessage) {
		logger.Infof("%s (%d) 向機器人 發送了消息: %s", msg.Sender.Nickname, msg.Sender.Uin, msg.ToString())
	})

	bot.OnTempMessage(func(cli *client.QQClient, event *client.TempMessageEvent) {
		msg := event.Message
		logger.Infof("%s (%d) 向機器人 發送了臨時會話消息: %s", msg.Sender.Nickname, msg.Sender.Uin, msg.ToString())
	})

	bot.OnNewFriendAdded(func(cli *client.QQClient, e *client.NewFriendEvent) {
		logger.Infof("已新增好友 %s (%d)", e.Friend.Nickname, e.Friend.Uin)
	})

	bot.OnNewFriendRequest(func(cli *client.QQClient, req *client.NewFriendRequest) {
		logger.Infof("收到好友請求 %s (%d): %s", req.RequesterNick, req.RequesterUin, req.Message)
	})

	bot.OnSelfPrivateMessage(func(cli *client.QQClient, msg *message.PrivateMessage) {
		friend := cli.FindFriend(msg.Target)
		logger.Infof("向 %s (%d) 發送私人消息: %s", friend.Nickname, msg.Target, msg.ToString())
	})

	bot.OnSelfGroupMessage(func(cli *client.QQClient, msg *message.GroupMessage) {
		logger.Infof("向群 %s (%d) 發送消息: %s", msg.GroupName, msg.GroupCode, msg.ToString())
	})
}

func init() {
	eventhook.HookLifeCycle(&log{})
}
