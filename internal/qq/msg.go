package qq

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/Mrs4s/MiraiGo/utils"
	"github.com/eric2788/MiraiValBot/internal/redis"
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

func CreateAtReply(source *message.GroupMessage) *message.SendingMessage {
	return CreateReply(source).Append(message.NewAt(source.Sender.Uin, source.Sender.DisplayName()))
}

func CreatePrivateReply(source *message.PrivateMessage) *message.SendingMessage {
	return message.NewSendingMessage().Append(message.NewPrivateReply(source))
}

func NewTts(text string) (*message.GroupVoiceElement, error) {
	return NewTtsWithGroup(ValGroupInfo.Uin, text)
}

func NewVoiceByUrl(url string) (*message.GroupVoiceElement, error) {
	data, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, err
	}
	return NewVoiceByBytes(data)
}

func NewVoiceByBytes(b []byte) (*message.GroupVoiceElement, error) {
	return NewVoiceByBytesWithGroup(ValGroupInfo.Uin, b)
}

func NewVoiceByBytesWithGroup(gp int64, b []byte) (*message.GroupVoiceElement, error) {
	return bot.Instance.UploadVoice(NewGroupSource(gp), bytes.NewReader(b))
}

func NewVoiceByUrlWithGroup(gp int64, url string) (*message.GroupVoiceElement, error) {
	b, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, err
	}
	return NewVoiceByBytesWithGroup(gp, b)
}

func NewTtsWithGroup(gp int64, text string) (voice *message.GroupVoiceElement, err error) {
	data, err := getTts(text)

	if err != nil {
		return nil, err
	}
	voice, err = bot.Instance.UploadVoice(NewGroupSource(gp), bytes.NewReader(data))
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

	voice, err = bot.Instance.UploadVoice(NewPrivateSource(uid), bytes.NewReader(data))
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
	b, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, err
	}
	return NewImagesByByteWithPrivate(uid, b)
}

func NewImageByUrlWithGroup(gp int64, url string) (*message.GroupImageElement, error) {
	b, err := request.GetBytesByUrl(url)
	if err != nil {
		return nil, err
	}
	return NewImagesByByteWithGroup(gp, b)
}

func NewImagesByByteWithGroup(gp int64, img []byte) (*message.GroupImageElement, error) {
	reader := bytes.NewReader(img)
	imgElement, err := bot.Instance.UploadImage(NewGroupSource(gp), reader)
	if err != nil {
		return nil, err
	}
	return imgElement.(*message.GroupImageElement), nil
}

func NewImagesByByteWithPrivate(uid int64, img []byte) (*message.FriendImageElement, error) {
	reader := bytes.NewReader(img)
	imgElement, err := bot.Instance.UploadImage(NewPrivateSource(uid), reader)
	if err != nil {
		return nil, err
	}
	return imgElement.(*message.FriendImageElement), nil
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
	return bot.Instance.UploadShortVideo(NewGroupSource(gp), bytes.NewReader(video), bytes.NewReader(thumb), 5)
}

func NewForwardNodeByGroup(msg *message.GroupMessage) *message.ForwardNode {

	// not sure why causing nil pointer reference
	filtered := make([]message.IMessageElement, 0)
	for _, ele := range msg.Elements {
		if ele != nil {
			filtered = append(filtered, ele)
		}
	}

	return &message.ForwardNode{
		GroupId:    msg.GroupCode,
		SenderId:   msg.Sender.Uin,
		Time:       msg.Time,
		SenderName: msg.Sender.DisplayName(),
		Message:    filtered,
	}
}

func NewForwardNodeByPrivate(msg *message.PrivateMessage) *message.ForwardNode {

	// not sure why causing nil pointer reference
	filtered := make([]message.IMessageElement, 0)
	for _, ele := range msg.Elements {
		if ele != nil {
			filtered = append(filtered, ele)
		}
	}

	return &message.ForwardNode{
		SenderId:   msg.Sender.Uin,
		Time:       msg.Time,
		SenderName: msg.Sender.DisplayName(),
		Message:    filtered,
	}
}

// NewForwardNode 以發送信息生成轉發信息節點
// 信息發送身份將為機器人自身
func NewForwardNode(msg *message.SendingMessage) *message.ForwardNode {
	return &message.ForwardNode{
		SenderId:   bot.Instance.Uin,
		Time:       int32(time.Now().Unix()),
		SenderName: bot.Instance.Nickname,
		Message:    msg.Elements,
	}
}

// NewMusicShare image and content can be empty string
func NewMusicShare(title, url, audio, image, content string) *message.ServiceElement {
	xml := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?><msg serviceID="2" templateID="1" action="web" brief="[分享] %s" sourceMsgId="0" url="%s" flag="0" adverSign="0" multiMsgFlag="0"><item layout="2"><audio cover="%s" src="%s"/><title>%s</title><summary>%s</summary></item><source name="音乐" icon="https://i.gtimg.cn/open/app_icon/01/07/98/56/1101079856_100_m.png" url="http://web.p.qq.com/qqmpmobile/aio/app.html?id=1101079856" action="app" a_actionData="com.tencent.qqmusic" i_actionData="tencent1101079856://" appid="1101079856" /></msg>`,
		utils.XmlEscape(title), url, image, audio, utils.XmlEscape(title), utils.XmlEscape(content))
	return &message.ServiceElement{
		Id:      60,
		Content: xml,
		SubType: "music",
	}
}

func NewGroupSource(gp int64) message.Source {
	return message.Source{
		PrimaryID:  gp,
		SourceType: message.SourceGroup,
	}
}

func NewPrivateSource(uid int64) message.Source {
	return message.Source{
		PrimaryID:  uid,
		SourceType: message.SourcePrivate,
	}
}
