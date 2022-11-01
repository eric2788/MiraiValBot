package aichat

import (
	"errors"
	"regexp"

	"github.com/Logiase/MiraiGo-Template/bot"
)

const (
	xiaoAiURL        = "http://jintia.jintias.cn/api/xatx.php?type=text&msg=%v"
)

var xiaoAiWarningMsg = regexp.MustCompile(`(?s)<br.+/>`)

type XiaoAi struct {
}

func (ai *XiaoAi) Reply(msg string) (string, error) {
	data, err := getAiReply(xiaoAiURL, msg)
	if err != nil {
		return "", err
	}

	nick := "Bot"

	if bot.Instance != nil {
		nick = bot.Instance.Nickname
	}

	replaced := xiaoAiWarningMsg.ReplaceAll(data, []byte(""))
	reply := replaces(string(replaced), map[string]string{
		"\n": "",
		"小爱": nick,
		"小米智能助理": "爹",
	})
	if reply == "" {
		return "", errors.New("无法获取回复讯息")
	}
	return reply, nil
}

func (ai *XiaoAi) Name() string {
	return "xiaoai"
}


