package qq

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/redis"
	"github.com/eric2788/MiraiValBot/utils/request"
)

func NewTextf(msg string, arg ...interface{}) *message.TextElement {
	return message.NewText(fmt.Sprintf(msg, arg...))
}

func NewTextfLn(msg string, arg ...interface{}) *message.TextElement {
	return message.NewText(fmt.Sprintf(msg+"\n", arg...))
}

func NewTextLn(msg string) *message.TextElement {
	return message.NewText(msg + "\n")
}

func CreateReply(source *message.GroupMessage) *message.SendingMessage {
	return message.NewSendingMessage().Append(message.NewReply(source))
}

func NewTts(text string) (*message.GroupVoiceElement, error) {
	return NewTtsWithGroup(ValGroupInfo.Uin, text)
}

func NewTtsWithGroup(gp int64, text string) (voice *message.GroupVoiceElement, err error) {

	key := toGroupKey(gp, fmt.Sprintf("tts:%x", md5.Sum([]byte(text))))

	var groupVoiceElement = &message.GroupVoiceElement{}

	if ok, err := redis.Get(key, groupVoiceElement); err != nil {
		return nil, err
	} else if ok {
		return groupVoiceElement, nil
	}

	data, err := getTTS(text)

	if err != nil {
		return nil, err
	}

	voice, err = bot.Instance.UploadGroupPtt(gp, bytes.NewReader(data))
	if err == nil {
		redisError := redis.Store(key, voice)
		if redisError != nil {
			logger.Warnf("Redis 儲存 群組語音消息 時出現錯誤: %v", redisError)
		}
	}
	return
}

func getTTS(text string) (data []byte, err error) {
	key := fmt.Sprintf("qq:tts:%x", md5.Sum([]byte(text)))

	data, notExist, err := redis.GetBytes(key)
	if err == nil || !notExist {
		return
	}
	data, err = bot.Instance.GetTts(text)
	if err == nil {
		redisError := redis.StoreBytes(key, data)
		if redisError != nil {
			logger.Warnf("Redis 儲存 TTS 時出現錯誤: %v", redisError)
		}
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

	data, err := getTTS(text)

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
