package verbose

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/file"
	qq2 "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

var translation = map[string]string{
	"started":        "的监听初始化成功",
	"stopped":        "的监听已关闭",
	"existed":        "的监听已经开始",
	"server-closed":  "已关闭",
	"server-started": "已启动",
	"error":          "的监听在初始化时出现错误",
}

type (
	liveRoomStatus struct {
		Platform string `json:"platform"`
		Id       string `json:"id"`
		Status   string `json:"status"`
	}

	liveRoomStatusHandler struct {
	}
)

func (s *liveRoomStatus) GetRoom() string {
	if s.Id == "-1" || s.Id == "server" {
		return "监控服务器"
	} else {
		return fmt.Sprintf("房间 %s", s.Id)
	}
}

func (s *liveRoomStatus) GetError() string {
	if strings.HasPrefix(s.Status, "error:") {
		return strings.Split(s.Status, ":")[1]
	} else {
		return ""
	}
}

func (l *liveRoomStatusHandler) GetOfflineListening() []string {
	return nil
}

func (l *liveRoomStatusHandler) HandleMessage(bot *bot.Bot, rd *redis.Message) {

	if !file.DataStorage.Setting.Verbose {
		return
	}

	var status liveRoomStatus
	if err := json.Unmarshal([]byte(rd.Payload), &status); err != nil {
		logger.Warnf("解析 JSON 内容时出现错误: %v", err)
		return
	}

	msg := message.NewSendingMessage()
	if err := status.GetError(); err != "" {
		msg.Append(qq2.NewTextf("【%s】%s 初始化监听时出现错误: %s", status.Platform, status.GetRoom(), err))
	} else {
		txt, ok := translation[status.Status]
		if !ok {
			txt = status.Status
		}
		msg.Append(qq2.NewTextf("【%s】%s %s", status.Platform, status.GetRoom(), txt))
	}

	go qq2.SendRiskyMessage(5, 5, func(try int) error {
		return qq2.SendGroupMessage(msg)
	})
}

func (l *liveRoomStatusHandler) HandleError(bot *bot.Bot, error error) {
}

func verboseLiveRoomStatus() {

	m, err := bot.GetModule(broadcaster.Tag)

	if err != nil {
		logger.Warnf("订阅房间状态讯息时出现错误: %v", err)
		return
	}

	bc := m.Instance.(*broadcaster.Broadcaster)

	if _, err = bc.Subscribe("live-room-status", &liveRoomStatusHandler{}); err != nil {
		logger.Warnf("订阅房间状态讯息时出现错误: %v", err)
	}

}
