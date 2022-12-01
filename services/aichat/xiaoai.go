package aichat

import (
	"errors"
	"regexp"
	"strings"
)

const (
	xiaoAiURL = "http://81.70.100.130/api/xiaoai.php?n=text&msg={msg}"
)

var xiaoAiWarningMsg = regexp.MustCompile(`(?s)<.+>`)

type XiaoAi struct {
}

func (ai *XiaoAi) Reply(msg string) (string, error) {
	data, err := getAiReply(xiaoAiURL, msg)
	if err != nil {
		return "", err
	}

	plain := strings.ReplaceAll(string(data), "\r\n", "")
	replaced := xiaoAiWarningMsg.ReplaceAll([]byte(plain), []byte(""))
	reply := replaces(string(replaced), map[string]string{
		"\n":     "",
		"小爱":     "我",
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
