package cmd

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
	"strconv"
)

func care(args []string, source *command.MessageSource) error {
	uid, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := message.NewSendingMessage().Append(message.NewReply(source.Message))

	if bilibili.AddHighlightUser(uid) {
		reply.Append(message.NewText(fmt.Sprintf("新增高亮用户 %d 成功。", uid)))
	} else {
		reply.Append(message.NewText(fmt.Sprintf("高亮用户 %d 已存在", uid)))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)

	return nil
}

func unCare(args []string, source *command.MessageSource) error {
	uid, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := message.NewSendingMessage().Append(message.NewReply(source.Message))

	if bilibili.RemoveHighlightUser(uid) {
		reply.Append(message.NewText(fmt.Sprintf("新增高亮用户 %d 成功。", uid)))
	} else {
		reply.Append(message.NewText(fmt.Sprintf("高亮用户 %d 不存在", uid)))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)

	return nil
}

func caring(args []string, source *command.MessageSource) error {

	reply := message.NewSendingMessage().Append(message.NewReply(source.Message))
	users := file.DataStorage.Bilibili.HighLightedUsers
	if len(users) > 0 {
		reply.Append(message.NewText(fmt.Sprintf("目前的高亮用户列表: %v", users)))
	} else {
		reply.Append(message.NewText("暂无高亮用户"))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

func clearInfo(args []string, source *command.MessageSource) error {
	room := int64(-1)

	if len(args) > 0 {
		r, err := strconv.ParseInt(args[0], 10, 64)

		if err != nil {
			return err
		}

		room = r
	}

	reply := message.NewSendingMessage().Append(message.NewReply(source.Message))

	if bilibili.ClearRoomInfo(room) {
		if room != -1 {
			reply.Append(message.NewText(fmt.Sprintf("已成功清除房间 %d 的资讯快取", room)))
		} else {
			reply.Append(message.NewText("已成功清除所有房间的快取"))
		}
	} else {
		reply.Append(message.NewText(fmt.Sprintf("房间 %d 没有资讯快取。", room)))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

func listen(args []string, source *command.MessageSource) error {
	room, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := message.NewSendingMessage().Append(message.NewReply(source.Message))

	result, err := bilibili.StartListen(room)

	if err != nil {
		reply.Append(message.NewText(fmt.Sprintf("监听房间时出现错误: %v", err)))
	} else if result {
		reply.Append(message.NewText(fmt.Sprintf("开始监听直播房间(%d)。", room)))
	} else {
		reply.Append(message.NewText(fmt.Sprintf("该直播间(%d)已经启动监听。", room)))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

func terminate(args []string, source *command.MessageSource) error {

	room, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := message.NewSendingMessage().Append(message.NewReply(source.Message))

	result, err := bilibili.StopListen(room)

	if err != nil {
		reply.Append(message.NewText(fmt.Sprintf("中止监听房间时出现错误: %v", err)))
	} else if result {
		reply.Append(message.NewText(fmt.Sprintf("已中止监听直播房间(%d)。", room)))
	} else {
		reply.Append(message.NewText(fmt.Sprintf("你尚未开始监听此直播房间。")))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

func listening(args []string, source *command.MessageSource) error {
	reply := message.NewSendingMessage().Append(message.NewReply(source.Message))
	listening := file.DataStorage.Listening.Bilibili
	if len(listening) > 0 {
		reply.Append(message.NewText(fmt.Sprintf("正在监听的房间号: %v", listening)))
	} else {
		reply.Append(message.NewText("没有正在监听的房间号"))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

var (
	careCommand      = command.NewNode([]string{"care", "高亮", "关注"}, "新增高亮用户", true, care, "<用户ID>")
	unCareCommand    = command.NewNode([]string{"uncare", "删除", "不高亮"}, "删除高亮用户", true, unCare, "<用户ID>")
	caringCommand    = command.NewNode([]string{"caring", "正在关注", "关注中", "关注列表"}, "获取高亮用户列表", false, caring)
	clearInfoCommand = command.NewNode([]string{"clearinfo", "清除快取"}, "清除房间资讯快取", true, clearInfo, "[roomId]")
	listenCommand    = command.NewNode([]string{"listen", "添加监听"}, "监听", true, listen, "<房间号>")
	terminateCommand = command.NewNode([]string{"terminate", "中止监听", "取消监听"}, "中止监听", true, terminate, "<房间号>")
	listeningCommand = command.NewNode([]string{"listening", "正在监听", "监听列表"}, "获取正在监听的房间号", false, listening)
)

var bliveCommand = command.NewParent([]string{"blive", "b站", "b站直播"}, "blive 直播间监听指令", false,
	careCommand, unCareCommand, caringCommand, clearInfoCommand, listenCommand, terminateCommand, listeningCommand,
)

func init() {
	command.AddCommand(bliveCommand)
}
