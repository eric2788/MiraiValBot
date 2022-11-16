package response

import (
	"crypto/md5"
	"encoding/binary"
	"math/rand"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/chat_reply"
	"github.com/eric2788/MiraiValBot/utils/misc"
)

const Tag = "valbot.response"

var (
	logger               = utils.GetModuleLogger(Tag)
	instance             = &response{}
	YesNoPattern         = regexp.MustCompile(`^.+是.+吗[\?？]*$`)
	questionMarkReplacer = strings.NewReplacer("?", "", "？", "")
)

type response struct {
	res *chat_reply.AIChatResponse
}

func (r *response) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (r *response) Init() {
	r.res = new(chat_reply.AIChatResponse)
}

func (r *response) PostInit() {
}

func (r *response) Serve(bot *bot.Bot) {
}

func (r *response) Start(bot *bot.Bot) {
	logger.Info("自定義回應模組已啟動")
}

func (r *response) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Info("自定義回應模組已關閉")
}

func (r *response) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(c *client.QQClient, msg *message.GroupMessage) {
		content := msg.ToString()

		if res, ok := file.DataStorage.Responses[content]; ok {
			m := message.NewSendingMessage().Append(message.NewText(res))
			_ = qq.SendGroupMessageByGroup(msg.GroupCode, m)
		} else if YesNoPattern.MatchString(content) {
			m := message.NewSendingMessage()
			if ans, ok := file.DataStorage.Answers[content]; ok {
				logger.Infof("此问题已被手动设置，因此使用被设置的回答")
				m.Append(message.NewText(getResponse(ans)))
			} else {
				ans = getQuestionAns(content)
				logger.Infof("自动回答问题 %s 为 %t", content, ans)
				m.Append(message.NewText(getResponse(ans)))
			}
			_ = qq.SendGroupMessageByGroup(msg.GroupCode, m)
		} else {

			rand.Seed(time.Now().UnixNano())

			// 1/50 (2%) 机率会回复
			if rand.Intn(50) == 25 {

				// 没有文字信息，随机发送龙图?
				if len(qq.ParseMsgContent(msg.Elements).Texts) == 0 {
					send, err := misc.NewRandomDragon()

					if err != nil {
						logger.Errorf("获取龙图失败: %v, 改为发送随机群消息", err)
						send, err = misc.NewRandomImage()
					}

					// 依然失败
					if err != nil {
						logger.Errorf("获取图片失败: %v, 放弃发送。", err)
						return
					}

					_ = qq.SendGroupMessageByGroup(msg.GroupCode, send)
					return
				}

				// 透过 AI 回复信息
				reply, err := r.res.Response(msg)
				if err != nil {
					logger.Errorf("透过 AI 回复对话时出现错误: %v", err)
				} else {

					// create a message with no reply element
					send := message.NewSendingMessage()

					for _, r := range reply.Elements {

						if _, ok := r.(*message.ReplyElement); ok {
							continue
						}

						send.Append(r)
					}

					_ = qq.SendGroupMessageByGroup(msg.GroupCode, send)
				}
			}

		}
	})
}

func getQuestionAns(content string) bool {
	hasher := md5.New()
	question := questionMarkReplacer.Replace(content)
	hashed := hasher.Sum([]byte(question))
	u64 := binary.BigEndian.Uint64(hashed)
	rand.Seed(int64(u64))
	return rand.Intn(2) == 1
}

func getResponse(is bool) string {
	if is {
		return "确实"
	} else {
		return "并不是"
	}
}

func init() {
	bot.RegisterModule(instance)
	eventhook.HookLifeCycle(instance)
}
