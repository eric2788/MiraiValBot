package broadcaster

import (
	"context"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/go-redis/redis/v8"
	"time"
)

type MessageHandler interface {
	GetOfflineListening() []string

	HandleMessage(bot *bot.Bot, message *redis.Message)

	HandleError(bot *bot.Bot, error error)
}

func (b *Broadcaster) SubscribeWithSite(topic string, site string) (bool, error) {
	if handler, ok := siteMap[site]; ok {
		return b.Subscribe(topic, handler), nil
	} else {
		return false, fmt.Errorf("未知的網站類型: %s", site)
	}
}

func (b *Broadcaster) Subscribe(topic string, handler MessageHandler) bool {
	if _, ok := b.subscribeMap[topic]; ok {
		return false
	}
	pubsub := b.rdb.Subscribe(ctx, topic)
	handleCtx, handleCancel := context.WithCancel(ctx)
	ifError := make(chan error)
	go handleMessage(topic, pubsub, handleCtx, ifError, func(msg *redis.Message) {
		handler.HandleMessage(b.bot, msg)
	})
	go handleError(topic, handleCtx, ifError, func(err error) {
		handler.HandleError(b.bot, err)
		b.UnSubscribe(topic)
		logger.Warnf("五秒後嘗試重新訂閱...")
		<-time.After(time.Second * 5)
		b.Subscribe(topic, handler)
	})
	b.subscribeMap[topic] = handleCancel
	logger.Infof("成功訂閱 %s\n", topic)
	return true
}

func (b *Broadcaster) UnSubscribe(topic string) bool {
	cancelSubscribe, ok := b.subscribeMap[topic]
	if !ok {
		return false
	}
	cancelSubscribe()
	return true
}

func handleError(topic string, ctx context.Context, ifError <-chan error, errorHandle func(error)) {
	select {
	case <-ctx.Done():
		return
	case err := <-ifError:
		logger.Warnf("接收訂閱 %s 時出現錯誤: %v\n", topic, err)
		go errorHandle(err)
		return
	}
}

func handleMessage(topic string, ps *redis.PubSub, ctx context.Context, ifError chan<- error, handle func(*redis.Message)) {
	defer func() {
		if err := ps.Close(); err != nil {
			logger.Warnf("停止訂閱 %s 時出現錯誤: %v\n", topic, err)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("已停止訂閱 %s", topic)
			return
		default:
			msg, err := ps.ReceiveMessage(ctx)
			if err != nil {
				go func() {
					ifError <- err
				}()
				return
			}
			handle(msg)
		}
	}
}
