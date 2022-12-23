package aichat

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/eric2788/MiraiValBot/services/copywriting"
	"github.com/eric2788/MiraiValBot/utils/test"
)

var chats = map[string]AIReply{
	"xiaoai":    &XiaoAi{},
	"qingyunke": &QingYunKe{},
	"tianxing":  &TianXing{},
	"moliyun":   &MoliYun{},
	"chatgpt3":  &Chatgpt3{},
}

func TestGetXiaoAi(t *testing.T) {

	aichat := chats["xiaoai"]

	msg, err := aichat.Reply("小米是垃圾")
	if err != nil {
		t.Skip(err)
	}

	t.Logf("Reply: %s", msg)
}

func TestGetMoliyun(t *testing.T) {
	aichat := chats["moliyun"]
	msg, err := aichat.Reply("你真笨")
	if err != nil {
		t.Skip(err)
	}

	t.Logf("Reply: %s", msg)
}

func TestQingYunKe(t *testing.T) {
	aichat := chats["qingyunke"]

	msg, err := aichat.Reply("你真垃圾")
	if err != nil {
		t.Skip(err)
	}

	t.Logf("Reply: %s", msg)
}

func TestTianXing_Reply(t *testing.T) {
	aichat := chats["tianxing"]

	if os.Getenv("TIAN_API_KEY") == "" {
		t.Skip("TIAN_API_KEY is empty")
	}

	msg, err := aichat.Reply("你好，你叫什么？")
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			t.Skip(err)
		}
		t.Fatal(err)
	}

	t.Logf("Reply: %s", msg)
}

func TestChatgpt3_Reply(t *testing.T) {

	if os.Getenv("CHATGPT_API_KEY") == "" {
		t.Skip("CHATGPT_API_KEY is empty")
	}

	conversations := []string{
		"你好，我叫老陈，你叫什么？",
		"你知道我叫什麼嗎",
	}

	for i, conversation := range conversations {

		ai := &Chatgpt3{}
		msg, err := ai.Reply(conversation)
		if err != nil {
			t.Skip(err)
		}
		<-time.After(time.Second * 5)
		t.Logf("Reply %d: %s", i+1, msg)
	}
}

func TestChatGptMaximumConversation(t *testing.T) {
	if os.Getenv("CHATGPT_API_KEY") == "" && os.Getenv("OPENAI_ACCESS_TOKEN") == "" {
		t.Skip("CHATGPT_API_KEY or OPENAI_ACCESS_TOKEN is empty")
	}

	questions, err := copywriting.GetTianGouList()
	if err != nil {
		t.Skip(err)
	}
	for i := 0; i < len(questions); i++ {
		<-time.After(time.Second * 1)
		ai := &Chatgpt3{}
		msg, err := ai.Reply(questions[i])
		if err != nil {
			t.Log(err)
			continue
		}
		t.Logf("Reply %d: %s", i+1, msg)
	}
}

func init() {
	test.InitTesting()
}
