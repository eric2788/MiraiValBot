package broadcaster

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/utils/set"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"strings"
)

type BroadcastHandler interface {
	PubSubPrefix() string
	GetOfflineListening() []string
	HandleLiveData(bot *bot.Bot, data interface{}, handle interface{}) error
	ToLiveData(message *redis.Message) (interface{}, error)
	HandleError(bot *bot.Bot, error error)
	GetCommand(data interface{}) string
}

type BroadCastHandle struct {
	logger     *logrus.Entry
	exception  set.StringSet
	handlerMap map[string]interface{}
	handler    BroadcastHandler
}

func (b *BroadCastHandle) GetOfflineListening() []string {
	return b.handler.GetOfflineListening()
}

func (b *BroadCastHandle) HandleMessage(bot *bot.Bot, message *redis.Message) {

	if !strings.HasPrefix(message.Channel, b.handler.PubSubPrefix()) {
		b.logger.Warnf("未知的 topic: %v", message.Channel)
		return
	}

	data, err := b.handler.ToLiveData(message)
	if err != nil {
		b.logger.Warnf("解析 JSON 內容時出現錯誤: %v\n", err)
		return
	}

	b.handleLiveData(bot, data)
}

func (b *BroadCastHandle) HandleError(bot *bot.Bot, error error) {
	b.handler.HandleError(bot, error)
}

func (b *BroadCastHandle) AddHandler(cmd string, handle interface{}) {
	b.handlerMap[cmd] = handle
	b.logger.Infof("已成功註冊 %s 指令的處理方法。", cmd)
}

func (b *BroadCastHandle) handleLiveData(bot *bot.Bot, data interface{}) {

	command := b.handler.GetCommand(data)

	if b.exception.Contains(command) {
		return
	}

	handle, ok := b.handlerMap[command]

	if !ok {
		b.logger.Warnf("找不到 %s 指令的處理方法，已略過。", command)
		b.exception.Add(command)
		return
	}

	// avoid handle panic
	defer func() {
		if err := recover(); err != nil {
			b.logger.Errorf("處理 %s 指令時出現嚴重錯誤: %v", command, err)
		}
	}()

	err := b.handler.HandleLiveData(bot, data, handle)

	if err != nil {
		b.logger.Warnf("處理 %s 指令時出現錯誤: %v", command, err)
	}
}

func BuildHandle(logger *logrus.Entry, handler BroadcastHandler) *BroadCastHandle {
	return &BroadCastHandle{
		logger:     logger,
		exception:  *set.NewString(),
		handlerMap: make(map[string]interface{}),
		handler:    handler,
	}
}
