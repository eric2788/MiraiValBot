package broadcaster

import (
	"runtime/debug"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/common-utils/set"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type BroadcastHandler[Data any] interface {
	PubSubPrefix() string
	GetOfflineListening() []string
	ToLiveData(message *redis.Message) (*Data, error)
	HandleError(bot *bot.Bot, error error)
	GetCommand(data *Data) string
}

type BroadCastHandle[Data any] struct {
	logger     logrus.FieldLogger
	exception  set.StringSet
	handlerMap map[string]func(bot *bot.Bot, data *Data) error
	handler    BroadcastHandler[Data]
}

func (b *BroadCastHandle[Data]) GetOfflineListening() []string {
	return b.handler.GetOfflineListening()
}

func (b *BroadCastHandle[Data]) HandleMessage(bot *bot.Bot, message *redis.Message) {

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

func (b *BroadCastHandle[Data]) HandleError(bot *bot.Bot, error error) {
	b.handler.HandleError(bot, error)
}

func (b *BroadCastHandle[Data]) AddHandler(cmd string, handle func(*bot.Bot, *Data) error) {
	b.handlerMap[cmd] = handle
	b.logger.Infof("已成功註冊 %s 指令的處理方法。", cmd)
}

func (b *BroadCastHandle[Data]) handleLiveData(bot *bot.Bot, data *Data) {

	command := b.handler.GetCommand(data)

	if b.exception.Contains(command) {
		return
	}

	handle, ok := b.handlerMap[command]

	if !ok {
		b.logger.Debugf("找不到 %s 指令的處理方法，已略過。", command)
		b.exception.Add(command)
		return
	}

	// avoid handle panic
	defer func() {
		if err := recover(); err != nil {
			b.logger.Errorf("處理 %s 指令時出現嚴重錯誤: %v", command, err)
			debug.PrintStack()
		}
	}()

	err := handle(bot, data)

	if err != nil {
		b.logger.Warnf("處理 %s 指令時出現錯誤: %v", command, err)
	}
}

func BuildHandle[Data any](logger logrus.FieldLogger, handler BroadcastHandler[Data]) *BroadCastHandle[Data] {
	return &BroadCastHandle[Data]{
		logger:     logger,
		exception:  *set.NewString(),
		handlerMap: make(map[string]func(bot *bot.Bot, data *Data) error),
		handler:    handler,
	}
}
