package aichat

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/modules/timer"
	"github.com/eric2788/chatgpt"
)

var (
	ctx *chatgpt.ChatContext
	mu  sync.Mutex
)

// openAI
type Chatgpt3 struct {
}

func (c *Chatgpt3) Reply(msg string) (string, error) {
	return c.replyWithAPI(msg)
}

func (c *Chatgpt3) replyWithAPI(msg string) (string, error) {
	apiKey := os.Getenv("CHATGPT_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("chatgpt api key not set")
	}

	c.ensureInit()

	cli := chatgpt.New(apiKey, "watta-aidraw", time.Duration(0))
	cli.ChatContext = ctx
	defer cli.Close()

	logger.Debugf("chatgpt3 接收信息: %s", msg)

	resp, err := cli.ChatWithContext(msg)

	if err != nil {
		logger.Debugf("chatgpt3 返回錯誤: %v", err)

		if err = c.handleError(err); err == nil {
			return c.Reply(msg) // retry again
		}

		return "", err
	}

	logger.Debugf("chatgpt3 輸出信息: %s", resp)

	return replaces(resp, map[string]string{
		"答：": "",
	}), nil
}

func (c *Chatgpt3) Name() string {
	return "chatgpt3"
}

func (c *Chatgpt3) ensureInit() {
	if ctx == nil {
		ctx = chatgpt.NewContext(
			chatgpt.WithMaxSeqTimes(3000),
			chatgpt.WithMaintainSeqTimes(true),
			chatgpt.WithOldConversation("data/aidraw.conversation"),
		)
		ctx.SetBackground("")
		ctx.SetPreset("")
	}
}

// handleError handle error that can be handled and return error that can not be handled
func (c *Chatgpt3) handleError(err error) error {
	// 已达到最大对话次数，重置对话上下文
	if err == chatgpt.OverMaxSequenceTimes {
		logger.Warnf("chatgpt3 已达到最大对话次数, 将重置对话上下文")
		ctx = nil
		return nil
	} else if err == chatgpt.OverMaxQuestionLength {
		return fmt.Errorf("chatgpt3 问题超过最大长度")
	} else if err == chatgpt.OverMaxTextLength {
		logger.Warnf("chatgpt3 文本超过最大长度，将尝试删除先前的对话记录: %d => %d", ctx.GetConversationsSize(), ctx.GetConversationsSize()-1)
		ctx.PollConversation()
		return nil
	}
	return err
}

func SaveGPTConversation() error {
	if ctx != nil {
		mu.Lock()
		defer mu.Unlock()
		return ctx.SaveConversation("data/aidraw.conversation")
	}
	return nil
}

func ResetGPTConversation() {
	if ctx == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	ctx.ResetConversation()
}

func isEdited() bool {
	if ctx == nil {
		return false
	}
	return ctx.IsEdited()
}

func init() {
	timer.RegisterTimer("save-gpt-conversation", 10*time.Minute, func(bot *bot.Bot) (err error) {
		if !isEdited() { //anti spam
			return
		}
		err = SaveGPTConversation()
		if err == nil {
			logger.Infof("保存AI对话成功")
		}
		return
	})
}
