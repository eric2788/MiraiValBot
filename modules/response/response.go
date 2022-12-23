package response

import (
	"math/rand"
	"reflect"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/game"
	"github.com/eric2788/MiraiValBot/services/copywriting"
	"github.com/eric2788/common-utils/array"
)

const Tag = "valbot.response"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &response{}

	longWongTalks = []string{
		"恭迎龙王 %s (跪拜)",
		"恭喜话痨 %s 成龙王咯",
		"口水多还得是你, %s",
		"%s, YOU 👆 ARE 👆 KING 👑",
		"你就是龙王 %s 吗, 不错",
	}

	pokeTalks = []string{
		"戳你妹戳戳戳, %s!",
		"我记住你了, %s!",
		"你是不是找打, %s?",
		"你戳我干嘛, %s?",
		"滚滚滚, %s!",
		"戳锤子戳, %s!",
		"泻药，刚醒, %s 找我何事",
	}

	responses []Handle
)

type (
	response struct {
		lastInGame bool
	}

	Handle interface {
		ShouldHandle(msg *message.GroupMessage) bool
		Handle(c *client.QQClient, msg *message.GroupMessage) error
	}
)

func (r *response) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(r.handleGroupMessage)
	bot.GroupNotifyEvent.Subscribe(r.handleGroupNotify)
}

func (r *response) handleGroupMessage(c *client.QQClient, msg *message.GroupMessage) {

	if game.IsInGame() {
		r.lastInGame = true
		return
	} else if r.lastInGame { // to make sure it wont reply immediately after game stopped
		r.lastInGame = false
		return
	} else if array.IndexOfInt64(qq.ParseMsgContent(msg.Elements).At, c.Uin) != -1 { // 防止跟 chat_reply 重复
		return
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(responses), func(i, j int) {
		responses[i], responses[j] = responses[j], responses[i]
	})
	for _, response := range responses {
		if response.ShouldHandle(msg) {
			if err := response.Handle(c, msg); err == nil {
				break
			} else {
				logger.Errorf("处理回应时出现错误: %v", err)
			}
		}
	}
}

func (r *response) handleGroupNotify(c *client.QQClient, event client.INotifyEvent) {

	// 非瓦群无视
	if event.From() != qq.ValGroupInfo.Uin {
		return
	}

	rand.Seed(time.Now().UnixNano())

	switch notify := event.(type) {
	case *client.GroupPokeNotifyEvent:

		// 機器人反戳無視
		if notify.Sender == c.Uin {
			return
		}

		msg := message.NewSendingMessage()
		sender := qq.FindGroupMember(notify.Sender)

		// 非机器人
		if notify.Receiver != c.Uin {

			receiver := qq.FindGroupMember(notify.Receiver)

			// 50% 触发CP
			if rand.Intn(100)+1 > 50 {

				list, atk, def, err := copywriting.GetCPList()
				if err != nil {
					logger.Errorf("获取CP列表失败: %v", err)
				} else {
					random := list[rand.Intn(len(list))]
					replacer := strings.NewReplacer(atk, sender.DisplayName(), def, receiver.DisplayName())
					msg.Append(message.NewText(replacer.Replace(random)))
					_ = qq.SendGroupMessage(msg)
				}

			}

			return
		}

		if rand.Intn(100)+1 > 10 {
			random := pokeTalks[rand.Intn(len(pokeTalks))]
			msg.Append(qq.NewTextfLn(random, sender.DisplayName()))
			// 戳回去咯
			c.SendGroupPoke(qq.ValGroupInfo.Code, notify.Sender)
		} else { // 10% 机率触发发病
			if success := sendWriting(msg, sender); !success {
				return
			}
		}

		_ = qq.SendGroupMessage(msg)

	case *client.MemberHonorChangedNotifyEvent:

		msg := message.NewSendingMessage()

		if notify.Uin == c.Uin {

			msg.Append(qq.NewTextf("机器人也能成 %s, 你群是不是该好好反思一下", qq.GetHonorString(notify.Honor)))
			msg.Append(message.NewFace(15))

		} else {

			user := qq.FindGroupMember(notify.Uin)

			// 80% 随机祝贺, 20% 发病
			if rand.Intn(100)+1 > 20 {
				if notify.Honor == client.Talkative {
					random := longWongTalks[rand.Intn(len(longWongTalks))]
					msg.Append(qq.NewTextf(random, user.DisplayName()))
				}
			} else {
				if success := sendWriting(msg, user); !success {
					return
				}
			}
		}

		_ = qq.SendGroupMessage(msg)

	}
}

func sendWriting(msg *message.SendingMessage, sender *client.GroupMemberInfo) bool {
	if rand.Intn(100) > 49 {
		return sendAsWriting(msg, sender)
	} else {
		return sendFabing(msg, sender)
	}
}

func sendAsWriting(msg *message.SendingMessage, sender *client.GroupMemberInfo) bool {
	list, err := copywriting.GetRanranList()
	if err != nil {
		logger.Errorf("获取小作文模板失败: %v", err)
		return false
	}
	random := list[rand.Intn(len(list))]
	msg.Append(message.NewText(strings.ReplaceAll(random.Text, random.Person, sender.DisplayName())))
	return true
}

func sendFabing(msg *message.SendingMessage, sender *client.GroupMemberInfo) bool {
	var getter func() ([]string, string, error)
	if rand.Intn(2) == 1 {
		getter = copywriting.GetFabingList
	} else {
		getter = copywriting.GetFadianList
	}
	if list, replace, err := getter(); err != nil {
		logger.Errorf("获取发病模板失败: %v", err)
		return false
	} else {
		random := list[rand.Intn(len(list))]
		msg.Append(message.NewText(strings.ReplaceAll(random, replace, sender.DisplayName())))
		return true
	}
}

func AddHandle(handle Handle) {
	responses = append(responses, handle)
	logger.Infof("新增 %s 回应成功", reflect.TypeOf(handle).Name())
}

func init() {
	eventhook.RegisterAsModule(instance, "自定義回應", Tag, logger)
}
