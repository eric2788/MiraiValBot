package qq

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/redis"
	"github.com/eric2788/common-utils/request"
)

type AppendableMessage struct {
	Texts []*message.TextElement
}

func NewTextf(msg string, arg ...interface{}) *message.TextElement {
	return message.NewText(fmt.Sprintf(msg, arg...))
}

func NewTextfLn(msg string, arg ...interface{}) *message.TextElement {
	return message.NewText(fmt.Sprintf(msg+"\n", arg...))
}

func NewTextLn(msg string) *message.TextElement {
	return message.NewText(msg + "\n")
}

func NextLn() *message.TextElement {
	return message.NewText("\n")
}

func CreateReply(source *message.GroupMessage) *message.SendingMessage {
	return message.NewSendingMessage().Append(message.NewReply(source))
}

func NewTts(text string) (*message.GroupVoiceElement, error) {
	return NewTtsWithGroup(ValGroupInfo.Uin, text)
}

func NewVoiceByUrl(url string) (*message.GroupVoiceElement, error){
	return NewVoiceByUrlWithGroup(ValGroupInfo.Uin, url)
}

func NewVoiceByBytes(b []byte) (*message.GroupVoiceElement, error){
	return NewVoiceByBytesWithGroup(ValGroupInfo.Uin, b)
}

func NewVoiceByBytesWithGroup(gp int64, b []byte) (*message.GroupVoiceElement, error){
	return bot.Instance.UploadGroupPtt(gp, bytes.NewReader(b))
}

func NewVoiceByUrlWithGroup(gp int64, url string) (*message.GroupVoiceElement, error){
	b, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, err
	}
	return NewVoiceByBytesWithGroup(gp, b)
}

func NewTtsWithGroup(gp int64, text string) (voice *message.GroupVoiceElement, err error) {

	key := GroupKey(gp, fmt.Sprintf("tts:%x", md5.Sum([]byte(text))))

	var groupVoiceElement = &message.GroupVoiceElement{}

	if ok, err := redis.Get(key, groupVoiceElement); err != nil {
		return nil, err
	} else if ok {
		logger.Infof("在redis 發現 「%v」 的 voiceElement 緩存， 將使用緩存", text)
		return groupVoiceElement, nil
	}

	logger.Infof("從 redis 找不到 voiceElement (%s), 將使用QQ上傳", key)

	data, err := getTts(text)

	if err != nil {
		return nil, err
	}

	voice, err = bot.Instance.UploadGroupPtt(gp, bytes.NewReader(data))
	if err == nil {
		redisError := redis.Store(key, voice)
		if redisError != nil {
			logger.Warnf("Redis 儲存 群組語音消息 時出現錯誤: %v", redisError)
		} else {
			logger.Infof("Redis 儲存 群組語音消息 成功。")
		}
	}
	return
}

func getTts(text string) (data []byte, err error) {
	key := fmt.Sprintf("qq:tts:%x", md5.Sum([]byte(text)))

	data, notExist, err := redis.GetBytes(key)

	// 非不存在的情況下出現錯誤
	if err != nil && !notExist {
		logger.Warnf("嘗試從 Redis 獲取 TTS 時出現錯誤: %v", err)
		return nil, err
	} else if err == nil { // 找到記錄
		logger.Infof("在redis 發現 「%v」 的 bytes 語音緩存， 將使用緩存", text)
		return data, nil
	}

	logger.Infof("redis 中找不到 TTS (%s), 將使用QQ上傳", key)

	data, err = bot.Instance.GetTts(text)

	if err == nil {
		redisError := redis.StoreBytes(key, data, redis.Permanent)
		if redisError != nil {
			logger.Warnf("Redis 儲存 TTS 時出現錯誤: %v", redisError)
		} else {
			logger.Infof("Redis 儲存 TTS 成功。")
		}
	} else {
		logger.Warnf("QQ 獲取 TTS 時出現錯誤: %v", err)
	}

	return
}

func NewTtsWithPrivate(uid int64, text string) (voice *message.PrivateVoiceElement, err error) {

	key := toPrivateKey(uid, fmt.Sprintf("tts:%x", md5.Sum([]byte(text))))

	var privateVoiceElement = &message.PrivateVoiceElement{}

	if ok, err := redis.Get(key, privateVoiceElement); err != nil {
		return nil, err
	} else if ok {
		return privateVoiceElement, nil
	}

	data, err := getTts(text)

	if err != nil {
		return nil, err
	}

	voice, err = bot.Instance.UploadPrivatePtt(uid, bytes.NewReader(data))
	if err == nil {
		redisError := redis.Store(key, voice)
		if redisError != nil {
			logger.Warnf("Redis 儲存 TTS 時出現錯誤: %v", redisError)
		}
	}
	return
}

func NewImageByUrl(url string) (*message.GroupImageElement, error) {
	return NewImageByUrlWithGroup(ValGroupInfo.Uin, url)
}

func NewImageByByte(img []byte) (*message.GroupImageElement, error) {
	return NewImagesByByteWithGroup(ValGroupInfo.Uin, img)
}

func NewImageByUrlWithPrivate(uid int64, url string) (*message.FriendImageElement, error) {
	img, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(img)
	return bot.Instance.UploadPrivateImage(uid, reader)
}

func NewImageByUrlWithGroup(gp int64, url string) (*message.GroupImageElement, error) {
	img, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(img)
	return bot.Instance.UploadGroupImage(gp, reader)
}

func NewImagesByByteWithGroup(gp int64, img []byte) (*message.GroupImageElement, error) {
	reader := bytes.NewReader(img)
	return bot.Instance.UploadGroupImage(gp, reader)
}

func NewImagesByByteWithPrivate(uid int64, img []byte) (*message.FriendImageElement, error) {
	reader := bytes.NewReader(img)
	return bot.Instance.UploadPrivateImage(uid, reader)
}

func NewVideoByUrl(url, thumbUrl string) (*message.ShortVideoElement, error) {
	return NewVideoByUrlWithGroup(ValGroupInfo.Uin, url, thumbUrl)
}

func NewVideoByUrlWithGroup(gp int64, url, thumbUrl string) (*message.ShortVideoElement, error) {
	video, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, fmt.Errorf("視頻解析失敗(%v)", err)
	}
	thumb, err := request.GetBytesByUrl(thumbUrl)
	if err != nil {
		return nil, fmt.Errorf("封面解析失敗(%v)", err)
	}
	return bot.Instance.UploadGroupShortVideo(gp, bytes.NewReader(video), bytes.NewReader(thumb))
}
