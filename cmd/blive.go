package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
	"strconv"
)

func bCare(args []string, source *command.MessageSource) error {
	uid, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := qq2.CreateReply(source.Message)

	if bilibili.AddHighlightUser(uid) {
		reply.Append(qq2.NewTextf("新增高亮用户 %d 成功。", uid))
	} else {
		reply.Append(qq2.NewTextf("高亮用户 %d 已存在", uid))
	}

	return qq2.SendGroupMessage(reply)
}

func bUnCare(args []string, source *command.MessageSource) error {
	uid, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := qq2.CreateReply(source.Message)

	if bilibili.RemoveHighlightUser(uid) {
		reply.Append(qq2.NewTextf("刪除高亮用户 %d 成功。", uid))
	} else {
		reply.Append(qq2.NewTextf("高亮用户 %d 不存在", uid))
	}

	return qq2.SendGroupMessage(reply)
}

func bCaring(args []string, source *command.MessageSource) error {

	reply := qq2.CreateReply(source.Message)
	users := file.DataStorage.Bilibili.HighLightedUsers
	if users.Size() > 0 {
		reply.Append(qq2.NewTextf("目前的高亮用户列表: %v", users.ToArr()))
	} else {
		reply.Append(message.NewText("暂无高亮用户"))
	}

	return qq2.SendWithRandomRiskyStrategy(reply)
}

func bClearInfo(args []string, source *command.MessageSource) error {
	room := int64(-1)

	if len(args) > 0 {
		r, err := strconv.ParseInt(args[0], 10, 64)

		if err != nil {
			return err
		}

		room = r
	}

	reply := qq2.CreateReply(source.Message)

	if bilibili.ClearRoomInfo(room) {
		if room != -1 {
			reply.Append(qq2.NewTextf("已成功清除房间 %d 的资讯快取", room))
		} else {
			reply.Append(message.NewText("已成功清除所有房间的快取"))
		}
	} else {
		reply.Append(qq2.NewTextf("房间 %d 没有资讯快取。", room))
	}

	return qq2.SendGroupMessage(reply)
}

func bListen(args []string, source *command.MessageSource) error {
	room, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := qq2.CreateReply(source.Message)

	result, err := bilibili.StartListen(room)

	if err != nil {
		reply.Append(qq2.NewTextf("监听房间时出现错误: %v", err))
	} else if result {
		reply.Append(qq2.NewTextf("开始监听直播房间(%d)。", room))
	} else {
		reply.Append(qq2.NewTextf("该直播间(%d)已经启动监听。", room))
	}

	return qq2.SendGroupMessage(reply)
}

func bTerminate(args []string, source *command.MessageSource) error {

	room, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return err
	}

	reply := qq2.CreateReply(source.Message)

	result, err := bilibili.StopListen(room)

	if err != nil {
		reply.Append(qq2.NewTextf("中止监听房间时出现错误: %v", err))
	} else if result {
		reply.Append(qq2.NewTextf("已中止监听直播房间(%d)。", room))
	} else {
		reply.Append(qq2.NewTextf("你尚未开始监听此直播房间。"))
	}

	return qq2.SendGroupMessage(reply)
}

func bListening(args []string, source *command.MessageSource) error {
	reply := qq2.CreateReply(source.Message)
	listening := file.DataStorage.Listening.Bilibili
	if listening.Size() > 0 {
		reply.Append(qq2.NewTextf("正在监听的房间号: %v", listening.ToArr()))
	} else {
		reply.Append(message.NewText("没有正在监听的房间号"))
	}

	return qq2.SendWithRandomRiskyStrategy(reply)
}

var (
	careCommand       = command.NewNode([]string{"care", "高亮", "关注"}, "新增高亮用户", true, bCare, "<用户ID>")
	unCareCommand     = command.NewNode([]string{"uncare", "删除", "不高亮"}, "删除高亮用户", true, bUnCare, "<用户ID>")
	caringCommand     = command.NewNode([]string{"caring", "正在关注", "关注中", "关注列表"}, "获取高亮用户列表", false, bCaring)
	clearInfoCommand  = command.NewNode([]string{"clearinfo", "清除快取"}, "清除房间资讯快取", true, bClearInfo, "[roomId]")
	bListenCommand    = command.NewNode([]string{"listen", "添加监听"}, "监听", true, bListen, "<房间号>")
	bTerminateCommand = command.NewNode([]string{"terminate", "中止监听", "取消监听"}, "中止监听", true, bTerminate, "<房间号>")
	bListeningCommand = command.NewNode([]string{"listening", "正在监听", "监听列表"}, "获取正在监听的房间号", false, bListening)
)

var bliveCommand = command.NewParent([]string{"blive", "b站", "b站直播"}, "blive 直播间监听指令",
	careCommand,
	unCareCommand,
	caringCommand,
	clearInfoCommand,
	bListenCommand,
	bTerminateCommand,
	bListeningCommand,
)

func init() {
	command.AddCommand(bliveCommand)
}
