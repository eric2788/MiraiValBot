package aichat

import "testing"


var chats = map[string]AIReply{
	"xiaoai": &XiaoAi{},
	"qingyunke": &QingYunKe{},
}

func TestGetXiaoAi(t *testing.T) {
	aichat := chats["xiaoai"]

	msg, err := aichat.Reply("你好，你叫什么？")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Reply: %s", msg)
}

func TestQingYunKe(t *testing.T) {
	aichat := chats["qingyunke"]

	msg, err := aichat.Reply("你好，你叫什么？")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("Reply: %s", msg)
}