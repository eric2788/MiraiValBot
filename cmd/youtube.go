package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func yConvert(args []string, source *command.MessageSource) error {
	url := args[0]
	id, err := youtube.GetChannelId(url)

	if err != nil {
		return err
	}

	reply := qq.CreateReply(source.Message)
	reply.Append(qq.NewTextf("该链接的频道URL为: %s", id))

	return qq.SendGroupMessage(reply)
}

func yListen(args []string, source *command.MessageSource) error {
	channelId := args[0]

	reply := qq.CreateReply(source.Message)
	result, err := youtube.StartListen(channelId)
	if err != nil {
		reply.Append(qq.NewTextf("启动监听时出现错误: %v", err))
	} else if result {
		reply.Append(qq.NewTextf("开始监听频道 %s", channelId))
	} else {
		reply.Append(qq.NewTextf("频道 %s 已启动监听", channelId))
	}

	return qq.SendGroupMessage(reply)
}

func yTerminate(args []string, source *command.MessageSource) error {
	channelId := args[0]
	reply := qq.CreateReply(source.Message)

	result, err := youtube.StopListen(channelId)

	if err != nil {
		reply.Append(qq.NewTextf("中止监听时出现错误: %v", err))
	} else if result {
		reply.Append(qq.NewTextf("已中止监听频道(%s)。", channelId))
	} else {
		reply.Append(message.NewText("你尚未开始监听此频道。"))
	}

	return qq.SendGroupMessage(reply)

}

func yListening(args []string, source *command.MessageSource) error {
	listening := file.DataStorage.Listening.Youtube

	reply := qq.CreateReply(source.Message)
	if listening.Size() > 0 {
		reply.Append(qq.NewTextf("正在监听的房间号: %v", listening))
	} else {
		reply.Append(qq.NewTextf("没有监听的房间号"))
	}

	return qq.SendGroupMessage(reply)
}

var (
	convertCommand    = command.NewNode([]string{"convert", "转换"}, "从用户名转换成频道ID", true, yConvert, "<用户名>")
	yListenCommand    = command.NewNode([]string{"listen", "监听", "启动监听", "启动"}, "启动频道监听", true, yListen, "<频道ID>")
	yTerminateCommand = command.NewNode([]string{"terminate", "中止", "取消"}, "中止监听频道", true, yTerminate, "<频道ID>")
	yListeningCommand = command.NewNode([]string{"listening", "监听中", "正在监听"}, "获取正在监听的频道 id", false, yListening)
)

var youtubeCommand = command.NewParent([]string{"youtube"}, "油管指令",
	convertCommand,
	yListenCommand,
	yListeningCommand,
	yTerminateCommand,
)

func init() {
	command.AddCommand(youtubeCommand)
}
