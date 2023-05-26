package broadcaster

import (
	"context"
	"fmt"
	mapset "github.com/deckarep/golang-set/v2"
	"runtime/debug"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/internal/file"
	rdb "github.com/eric2788/MiraiValBot/internal/redis"
	"github.com/go-redis/redis/v8"
)

// waitForPubSubClose this set is to avoid subscribe while that channel is closing PubSub
var waitForPubSubClose = mapset.NewSet[string]()

type MessageHandler interface {
	GetOfflineListening() []string
	HandleMessage(bot *bot.Bot, message *redis.Message)
	HandleError(bot *bot.Bot, error error)
}

func (b *Broadcaster) SubscribeWithSite(topic string, site string) (bool, error) {
	if handler, ok := siteMap[site]; ok {
		return b.Subscribe(topic, handler)
	} else {
		return false, fmt.Errorf("未知的網站類型: %s", site)
	}
}

func (b *Broadcaster) Subscribe(topic string, handler MessageHandler) (bool, error) {

	if waitForPubSubClose.Contains(topic) {
		return false, fmt.Errorf("上一次的解除訂閱尚未完成，請稍候再嘗試。")
	}

	if _, ok := b.subscribeMap[topic]; ok {
		return false, nil
	}

	// 手動關閉的 context
	handleCtx, handleCancel := context.WithCancel(ctx)

	pubsub := rdb.Subscribe(handleCtx, topic)
	ifError := make(chan error, 1)

	// pubsub 關閉的 context
	pubsubCtx, pubsubClose := context.WithCancel(ctx)

	go handleMessage(topic, pubsub, handleCtx, pubsubClose, ifError, handler)

	go handleError(topic, handleCtx, ifError, func(err error) {
		handler.HandleError(bot.Instance, err)
		b.UnSubscribe(topic)
		logger.Warnf("五秒後嘗試重新訂閱...")
		<-time.After(time.Second * 5)
		if _, err = b.Subscribe(topic, handler); err != nil {
			logger.Warnf("重新訂閱時出現錯誤: %v", err)
		} else {
			logger.Infof("重新訂閱 %s 成功。", topic)
		}
	})

	b.subscribeMap[topic] = Subscriber{pubsubCtx, handleCancel}
	logger.Infof("成功訂閱 %s", topic)
	return true, nil
}

func (b *Broadcaster) UnSubscribe(topic string) bool {
	subscriber, ok := b.subscribeMap[topic]
	if !ok || waitForPubSubClose.Contains(topic) {
		return false
	}
	subscriber.CancelFunc()
	waitForPubSubClose.Add(topic)
	logger.Debugf("[Subscribe] 已添加 topic %s 到等待關閉列表", topic)

	go func() {
		<-subscriber.Context.Done()
		delete(b.subscribeMap, topic)
		waitForPubSubClose.Remove(topic)
		logger.Debugf("[Subscribe] pubsub 已成功關閉，已刪除 topic %s 到等待關閉列表", topic)
	}()

	return true
}

func handleError(topic string, ctx context.Context, ifError <-chan error, errorHandle func(error)) {
	select {
	case <-ctx.Done():
		return
	case err := <-ifError:
		logger.Warnf("接收訂閱 %s 時出現錯誤: %v", topic, err)
		go errorHandle(err)
		return
	}
}

func handleMessage(topic string, ps *redis.PubSub, ctx context.Context, cancel context.CancelFunc, ifError chan<- error, handle MessageHandler) {
	defer func() {
		if err := ps.Close(); err != nil {
			logger.Warnf("停止訂閱 %s 時出現錯誤: %v", topic, err)
		}
		logger.Infof("%s 的訂閱已停止。", topic)
		cancel()
	}()
	defer close(ifError)
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("處理 %s 訂閱訊息時出現致命錯誤: %v", topic, err)
			debug.PrintStack()
			ifError <- fmt.Errorf("%v", err)
		}
	}()
	size := file.ApplicationYaml.Redis.Buffer // 每次最大接收数量 (buffer)
	channel := ps.Channel(
		redis.WithChannelSize(int(size)),
		redis.WithChannelHealthCheckInterval(time.Minute),
	)
	for {
		select {
		case <-ctx.Done():
			logger.Debugf("收到中止指令，正在停止訂閱 %s", topic)
			return
		case msg, ok := <-channel:
			if !ok {
				logger.Debugf("訂閱接收閘口關閉，正在停止訂閱 %s", topic)
				return
			}
			go handle.HandleMessage(bot.Instance, msg)
		}
	}
}
