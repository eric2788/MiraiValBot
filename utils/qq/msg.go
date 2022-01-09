package qq

import (
	"bytes"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"io/ioutil"
	"net/http"
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

func NewImageByUrl(url string) (*message.GroupImageElement, error) {
	return NewImageByUrlWithGroup(ValGroupInfo.Uin, url)
}

// NewImageByUrlWithGroup TODO make caching with redis ?
func NewImageByUrlWithGroup(gp int64, url string) (*message.GroupImageElement, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
	}()
	img, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(img)
	return bot.Instance.UploadGroupImage(gp, reader)
}
