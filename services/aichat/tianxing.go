package aichat

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	tianApi = "https://apis.tianapi.com/robot/index?&key=%v&question={msg}"
)

type (
	TianXing struct {
	}

	tianXingReply struct {
		Code   int    `json:"code"`
		Msg    string `json:"msg"`
		Result struct {
			Reply    string `json:"reply"`
			DataType string `json:"datatype"`
		} `json:"result,omitempty"`
	}
)

func (ai *TianXing) Reply(msg string) (string, error) {
	key := os.Getenv("TIAN_API_KEY")
	url := fmt.Sprintf(tianApi, key)
	data, err := getAiReply(url, msg)
	if err != nil {
		return "", err
	}
	reply := new(tianXingReply)
	err = json.Unmarshal(data, &reply)
	if err != nil {
		return "", err
	}
	if reply.Code != 200 {
		return "", fmt.Errorf("%d: %s", reply.Code, reply.Msg)
	}
	return replaces(reply.Result.Reply, map[string]string{"天行": "你爹的"}), nil
}

func (ai *TianXing) Name() string {
	return "tianxing"
}
