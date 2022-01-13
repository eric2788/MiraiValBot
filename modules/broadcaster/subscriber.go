package broadcaster

import (
	"context"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/file"
	rdb "github.com/eric2788/MiraiValBot/redis"
	"github.com/eric2788/MiraiValBot/utils/set"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"runtime/debug"
	"time"
)

// waitForPubSubClose this set is to avoid subscribe while that channel is closing PubSub
var waitForPubSubClose = set.NewString()

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
	ifError := make(chan error)

	// pubsub 關閉的 context
	pubsubCtx, pubsubClose := context.WithCancel(ctx)

	go handleMessage(topic, pubsub, handleCtx, pubsubClose, ifError, func(msg *redis.Message) {
		handler.HandleMessage(bot.Instance, msg)
	})

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
		waitForPubSubClose.Delete(topic)
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

func handleMessage(topic string, ps *redis.PubSub, ctx context.Context, close context.CancelFunc, ifError chan<- error, handle func(*redis.Message)) {
	defer func() {
		if err := ps.Close(); err != nil {
			logger.Warnf("停止訂閱 %s 時出現錯誤: %v", topic, err)
		}
		logger.Infof("%s 的訂閱已停止。", topic)
		close()
	}()
	defer func() {
		if err := recover(); err != nil {
			go func() {
				logger.Errorf("處理 %s 訂閱訊息時出現致命錯誤: %v, from %v", topic, err, debug.Stack())
				ifError <- fmt.Errorf("%v", err)
			}()
		}
	}()
	size := file.ApplicationYaml.Redis.Buffer // 每次最大接收数量 (buffer)
	channel := make(chan *redis.Message, int(size))
	go receivePubsub(ps, ctx, channel, ifError)
	for {
		select {
		case <-ctx.Done():
			logger.Debugf("收到中止指令，正在停止訂閱 %s", topic)
			return
		case msg := <-channel:
			handle(msg)
		default:
			break
		}
	}
}

func receivePubsub(ps *redis.PubSub, ctx context.Context, channel chan<- *redis.Message, ifError chan<- error) {
	for {
		// msg, err := fakeReceiveWithError(ps, ctx)
		msg, err := ps.ReceiveMessage(ctx)
		if err != nil {
			go func() {
				ifError <- err
			}()
			return
		}
		channel <- msg
	}
}

func fakeReceiveWithError(ps *redis.PubSub, ctx context.Context) (msg *redis.Message, err error) {
	msg, err = ps.ReceiveMessage(ctx)
	if rand.Intn(1000)%10 == 0 {
		err = fmt.Errorf("test Error")
	}
	return
}
