package aichat

import (
	"os"
	"strings"
	"testing"

	"github.com/eric2788/MiraiValBot/utils/test"
)

var chats = map[string]AIReply{
	"xiaoai":    &XiaoAi{},
	"qingyunke": &QingYunKe{},
	"tianxing":  &TianXing{},
	"moliyun": &MoliYun{},
}

func TestGetXiaoAi(t *testing.T) {

	aichat := chats["xiaoai"]

	msg, err := aichat.Reply("小米是垃圾")
	if err != nil {
		t.Skip(err)
	}

	t.Logf("Reply: %s", msg)
}

func TestGetMoliyun(t *testing.T){
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

func init() {
	test.InitTesting()
}
