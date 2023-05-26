package aidraw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/modules/timer"
	"github.com/eric2788/common-utils/request"
	"github.com/eric2788/common-utils/stream"
	"math/rand"
	"strings"
	"time"
)

type (
	saiResp struct {
		HasError     bool                   `json:"hasError"`
		ErrorMessage string                 `json:"errorMessage"`
		Payload      map[string]interface{} `json:"payload"`
	}

	saiImgInfo struct {
		ImageId string `json:"imageID"`
		Status  string `json:"status"`
	}

	saiGenImgInfo struct {
		saiImgInfo
		IgnoredWords string `json:"ignoredWords"`
	}

	saiImgStatusInfo struct {
		saiImgInfo
		Url string `json:"url"`
	}

	saiOTPResult struct {
		Success bool `json:"success"`
	}

	saiAuthResult struct {
		SessionID       string `json:"sessionID"`
		UserID          string `json:"userID"`
		Email           string `json:"email"`
		UserName        string `json:"username"`
		IsAuthenticated bool   `json:"isAuthenticated"`
		IsDisabled      *bool  `json:"isDisabled"`
	}
)

func (resp *saiResp) Scan(res interface{}) error {
	b, err := json.Marshal(resp.Payload)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, res)
}

var (
	sexyAIRequester = request.New(
		request.WithBaseUrl("https://api.sexy.ai"),
		request.WithHeaders(map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) ",
			"Origin":     "https://sexy.ai",
			"Referer":    "https://sexy.ai/",
		}),
	)

	modelMap = map[string]string{
		"real1": "model2",
		"real2": "model3",
		"real3": "model1",
		"anime": "model4",
	}

	random = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func sexyAIDraw(payload Payload) (*Response, error) {

	var model string
	if m, ok := modelMap[payload.Model]; ok {
		model = m
	} else if payload.Model == "real" { // 随机选择一个real模型
		model = stream.FromMap(modelMap).
			Filter(func(k string, v string) bool {
				return strings.HasPrefix(k, "real")
			}).
			Values().
			Shuffle().
			MustFirst()
	} else {
		return nil, fmt.Errorf("未知模型: %v, 可用模型: %v", payload.Model, stream.FromMap(modelMap).Keys().Join(", "))
	}

	logger.Debugf("Requesting sexyAI Image with sessionID: %v", file.DataStorage.AiDraw.SexyAISession)

	var res saiResp
	_, err := sexyAIRequester.Post("/generateImage", &res,
		request.Data(map[string]interface{}{
			"modelName": model,
			"prompt":    prefixPrompt + payload.Prompt,
			"negprompt": badPrompt,
			"seed":      random.Uint64(),
			"sessionID": file.DataStorage.AiDraw.SexyAISession,
			"steps":     20,
		}),
	)
	if err != nil {
		return nil, err
	} else if res.HasError {
		return nil, errors.New(res.ErrorMessage)
	}

	var genImgInfo saiGenImgInfo
	err = res.Scan(&genImgInfo)
	if err != nil {
		return nil, fmt.Errorf("解析生成图片信息失败: %v", err)
	}

	if genImgInfo.Status != "generating" {
		return nil, fmt.Errorf("生成图片失败: %v", genImgInfo.Status)
	}

	resultPhoto, err := waitForImageGenerated(genImgInfo.ImageId)
	if err != nil {
		return nil, err
	}

	return &Response{
		ImgUrl: resultPhoto,
	}, nil
}

// waitForImageGenerated blocked function
func waitForImageGenerated(imgId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
	for {
		select {
		case <-ctx.Done():
			cancel()
			return "", errors.New("等待图片生成超时")
		default:
			<-time.After(time.Second * 5)
			res, err := sexyPost("/getImageStatus", request.Data(map[string]interface{}{
				"imageID":   imgId,
				"sessionID": file.DataStorage.AiDraw.SexyAISession,
			}))
			if err != nil {
				cancel()
				return "", fmt.Errorf("获取图片状态失败: %v", err)
			}
			var statusInfo saiImgStatusInfo
			err = res.Scan(&statusInfo)
			if err != nil {
				cancel()
				return "", fmt.Errorf("解析图片状态信息失败: %v", err)
			}
			switch statusInfo.Status {
			case "pending":
				logger.Infof("图片 %v 正在生成中...", imgId)
				// continue
			case "complete":
				cancel()
				return statusInfo.Url, nil
			default:
				cancel()
				return "", fmt.Errorf("图片 %v 生成失败: %v", imgId, statusInfo.Status)
			}
		}
	}
}

func init() {
	initializeSession()
}

func initializeSession() {
	res, err := sexyPost("/getSelfUser", request.Data(map[string]interface{}{
		"isAtLeast18Confirmed": true,
		"sessionID":            file.DataStorage.AiDraw.SexyAISession,
	}))
	if err != nil {
		logger.Errorf("sexyAI初始化失败: %v", err)
		return
	}
	var auth saiAuthResult
	err = res.Scan(&auth)
	if err != nil {
		logger.Errorf("sexyAI初始化失败: %v", err)
		return
	}

	if !auth.IsAuthenticated && file.DataStorage.AiDraw.SexyAISession != "" {
		logger.Infof("檢測到 SexyAI Session 已失效, 正在刷新新的 Session ID...")
		file.UpdateStorage(func() {
			file.DataStorage.AiDraw.SexyAISession = ""
		})
		initializeSession()
		return
	}

	if auth.SessionID != file.DataStorage.AiDraw.SexyAISession {
		logger.Infof("检测到 SexyAI 的 SessionId 已变更: %v -> %v", file.DataStorage.AiDraw.SexyAISession, auth.SessionID)
		file.UpdateStorage(func() {
			file.DataStorage.AiDraw.SexyAISession = auth.SessionID
		})
	}

	logger.Infof("成功獲取 SexyAI SessionID: %v, Username: %v", file.DataStorage.AiDraw.SexyAISession, auth.UserName)
	timer.RegisterTimer("sexyAI-keepAlive", time.Minute*5, keepAlive)

	drawableSources["sexyai"] = sexyAIDraw
}

func SaiRequestOTP(email string, newUser bool) (bool, error) {
	res, err := sexyPost("/requestOTP", request.Data(map[string]interface{}{
		"email":     email,
		"newUser":   newUser,
		"sessionID": file.DataStorage.AiDraw.SexyAISession,
	}))
	if err != nil {
		return false, err
	}
	var otp saiOTPResult
	err = res.Scan(&otp)
	if err != nil {
		return false, fmt.Errorf("解析OTP信息失败: %v", err)
	}
	return otp.Success, nil
}

func SaiAuth(email, otp string, register bool) (string, error) {
	path := "/authenticateUser"
	formData := map[string]interface{}{
		"email":                email,
		"isAtLeast18Confirmed": true,
		"otp":                  otp,
		"sessionID":            file.DataStorage.AiDraw.SexyAISession,
	}
	if register {
		path = "/createUser"
		formData["refList"] = ""
		formData["rs"] = ""
		// generate a random username
		formData["username"] = fmt.Sprintf("user%v", random.Intn(1000000))
	}

	logger.Debugf("requesting with data: %v", formData)

	res, err := sexyPost(path, request.Data(formData))
	if err != nil {
		return "", err
	}
	var auth saiAuthResult
	err = res.Scan(&auth)
	if err != nil {
		return "", fmt.Errorf("解析Auth信息失败: %v", err)
	}

	if !auth.IsAuthenticated {
		return auth.UserName, fmt.Errorf("登錄失败: %v", auth.Email)
	}

	return auth.UserName, nil
}

func keepAlive(bot *bot.Bot) error {
	res, err := sexyPost("/getSelfUser", request.Data(map[string]interface{}{
		"isAtLeast18Confirmed": true,
		"sessionID":            file.DataStorage.AiDraw.SexyAISession,
	}))
	var auth saiAuthResult
	err = res.Scan(&auth)
	if err != nil {
		return err
	}

	if auth.IsAuthenticated {
		return nil
	}

	logger.Warnf("SexyAI SessionID %v 已失效，需要重新登入", file.DataStorage.AiDraw.SexyAISession)
	return nil
}

func sexyPost(url string, confs ...request.Configurer) (*saiResp, error) {
	var res saiResp
	_, err := sexyAIRequester.Post(url, &res, confs...)
	if err != nil {
		return nil, err
	} else if res.HasError {
		return nil, errors.New(res.ErrorMessage)
	}
	return &res, nil
}
