package aichat

import (
	"encoding/json"
	"fmt"
)

const qingyunkeURL = "http://api.qingyunke.com/api.php?key=free&appid=0&msg={msg}"

type (
	qingYunKeReply struct {
		Result  int    `json:"result"`
		Content string `json:"content"`
	}

	QingYunKe struct {
	}
)

func (q *QingYunKe) Reply(msg string) (string, error) {
	data, err := getAiReply(qingyunkeURL, msg)
	if err != nil {
		return "", err
	}
	var reply qingYunKeReply
	if err := json.Unmarshal(data, &reply); err != nil {
		return "", err
	}
	if reply.Result != 0 {
		return "", fmt.Errorf("%d: %s", reply.Result, reply.Content)
	} else {
		return replaces(reply.Content, map[string]string{
			"小美人菲菲": "你爹",
			"小美人":   "你爹",
			"菲菲":    "我",
		}), nil
	}
}

func (q *QingYunKe) Name() string {
	return "qingyunke"
}
