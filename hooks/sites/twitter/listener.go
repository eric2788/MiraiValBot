package twitter

import (
	"fmt"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/file"
	bc "github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/eric2788/common-utils/request"
)

var (
	listening = &file.DataStorage.Listening
	userCache = make(map[string]*ExistUserResp)
	logger    = utils.GetModuleLogger("sites.twitter")
	topic     = func(user string) string { return fmt.Sprintf("twitter:%s", user) }
)

type ExistUserResp struct {
	Exist bool        `json:"exist"`
	Data  ProfileData `json:"data"`
}

type ProfileData struct {
	Id         string `json:"id"`
	ScreenName string `json:"screen_name"`
}

type ErrorResp struct {
	Error []string `json:"error"`
}

func StartListen(user string) (bool, error) {

	if info, err := GetUserInfo(user); err != nil {
		return false, err
	} else if !info.Exist {
		return false, fmt.Errorf("用戶不存在")
	}

	file.UpdateStorage(func() {
		(*listening).Twitter.Add(user)
	})

	info, _ := bot.GetModule(bc.Tag)

	broadcaster := info.Instance.(*bc.Broadcaster)

	return broadcaster.Subscribe(topic(user), MessageHandler)
}

func StopListen(user string) (bool, error) {

	if !(*listening).Twitter.Contains(user) {
		return false, nil
	}

	file.UpdateStorage(func() {
		(*listening).Twitter.Remove(user)
	})

	info, _ := bot.GetModule(bc.Tag)

	broadcaster := info.Instance.(*bc.Broadcaster)

	result := broadcaster.UnSubscribe(topic(user))

	return result, nil
}

func GetUserInfo(user string) (*ExistUserResp, error) {
	if resp, ok := userCache[user]; ok {
		return resp, nil
	}

	var existUserResp = &ExistUserResp{}

	if err := request.Get(fmt.Sprintf("%s/%s", file.ApplicationYaml.Val.TwitterLookUpUrl, user), existUserResp); err != nil {
		if httpError, ok := err.(*request.HttpError); ok {

			var errorResp = &ErrorResp{}

			if err = request.Read(httpError.Response, errorResp); err != nil {
				logger.Warnf("讀取 TwitterLookup 請求時出現錯誤: %v\n", err) // log here
				return nil, httpError                               // return httpError instead of marshall error
			}

			return nil, fmt.Errorf("檢查用戶存在時出現錯誤: %s", strings.Join(errorResp.Error, ", "))

		} else {
			return nil, err
		}
	}

	userCache[user] = existUserResp

	return existUserResp, nil
}
