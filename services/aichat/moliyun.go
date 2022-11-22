package aichat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/eric2788/MiraiValBot/internal/qq"
)

const moliyunURL = "https://api.mlyai.com/reply"

type (
	MoliYun struct {
	}

	// reference: https://mlyai.com/docs
	moliYunPayload struct {
		Content string `json:"content"`
		Type    int    `json:"type"`
		From    int64  `json:"from"`
		To      int64  `json:"to"`
	}

	moliyunResp struct {
		Code    string  `json:"code"`
		Message string  `json:"message"`
		Plugin  *string `json:"plugin"`
		Data    []struct {
			Content string `json:"content"`
			Typed   int    `json:"typed"`
		} `json:"data"`
	}
)

func (m *MoliYun) Reply(msg string) (string, error) {
	key, secret := os.Getenv("MOLIYUN_API_KEY"), os.Getenv("MOLIYUN_API_SECRET")
	if key == "" || secret == "" {
		return "", fmt.Errorf("API KEY or API SECRET not set")
	}
	b, err := json.Marshal(moliYunPayload{
		Content: msg,
		Type:    2,
		From:    qq.ValGroupInfo.OwnerUin,
		To:      qq.ValGroupInfo.Code,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, moliyunURL, bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Api-Key", key)
	req.Header.Set("Api-Secret", secret)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err = io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	resp := new(moliyunResp)
	err = json.Unmarshal(b, resp)
	if err != nil {
		return "", err
	}
	if resp.Code != "00000" {
		return "", fmt.Errorf(resp.Message)
	}
	for _, d := range resp.Data {
		if d.Content != "" {
			return replaces(d.Content, map[string]string{
				"茉莉云": "你爹的",
			}), nil
		}
	}
	return "", fmt.Errorf("没有回应: %v", resp.Data)
}

func (m *MoliYun) Name() string {
	return "moliyun"
}
