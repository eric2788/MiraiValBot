package broadcaster

import (
	"context"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/go-redis/redis/v8"
	"sync"
)

const Tag = "valbot.broadcaster"

type Subscriber struct {
	Context    context.Context
	CancelFunc context.CancelFunc
}

var (
	instance = &Broadcaster{
		subscribeMap: make(map[string]Subscriber),
	}
	logger  = utils.GetModuleLogger(Tag)
	ctx     = context.Background()
	siteMap = make(map[string]MessageHandler)
)

func init() {
	bot.RegisterModule(instance)
}

type Broadcaster struct {
	rdb          *redis.Client
	subscribeMap map[string]Subscriber
	bot          *bot.Bot
}

func (b *Broadcaster) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (b *Broadcaster) Init() {
	redisConfig := file.ApplicationYaml.Redis
	host := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	b.rdb = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: redisConfig.Password,
		DB:       redisConfig.Database,
	})
}

func (b *Broadcaster) PostInit() {
	// 啟動 discord 機器人以作廣播通知
	discord.Start()
}

func (b *Broadcaster) Serve(bot *bot.Bot) {
	b.bot = bot

	logger.Info("正在重離線獲取訂閱...")
	count := 0
	for _, handler := range siteMap {
		// 獲取所有離線訂閱
		for _, topic := range handler.GetOfflineListening() {
			// 進行訂閱
			_, _ = b.Subscribe(topic, handler)
			count += 1
		}
	}
	logger.Infof("已從離線重新訂閱 %d 個 topic。", count)
}

func (b *Broadcaster) Start(bot *bot.Bot) {
	logger.Info("Redis 訂閱監控已啟用。")
}

func (b *Broadcaster) Stop(bot *bot.Bot, wg *sync.WaitGroup) {

	defer wg.Done()

	subWaitGroup := &sync.WaitGroup{}
	subWaitGroup.Add(len(b.subscribeMap))

	// 解除所有訂閱
	for topic, subscriber := range b.subscribeMap {
		ctx, cancel := subscriber.Context, subscriber.CancelFunc
		go func(topic string) {
			cancel()
			logger.Debugf("[Redis] 正在等待 topic %s 關閉 pubsub...", topic)
			<-ctx.Done() // 等待 pubsub 關閉
			logger.Debugf("[Redis] topic %s pubsub 關閉完成", topic)
			subWaitGroup.Done()
		}(topic)
	}

	subWaitGroup.Wait()

	// 關閉 Redis
	if err := b.rdb.Close(); err != nil {
		logger.Warnf("關閉 Redis 時出現錯誤: %v", err)
	} else {
		logger.Info("Redis 已成功關閉連接")
	}

}

func RegisterHandler(site string, handler MessageHandler) bool {
	if _, ok := siteMap[site]; ok {
		// site already exist
		return false
	} else {
		siteMap[site] = handler
		logger.Infof("已成功註冊網站 %s 的 PubSub 處理方式", site)
		return true
	}
}
