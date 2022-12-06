package test

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/sirupsen/logrus"
	"github.com/subosito/gotenv"
)

var logger = utils.GetModuleLogger("utils.test")

func InitTesting() {
	logrus.SetLevel(logrus.DebugLevel)
	logger.Debugf("Logging Level set to debug")
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		logger.Warnf("unable to get the current filename")
		return
	}
	dirname := filepath.Dir(filename)

	path := "/"

	if runtime.GOOS == "windows" {
		path = "\\"
	}

	if err := gotenv.OverLoad(dirname + path + ".env.local"); err == nil {
		logger.Debugf("successfully loaded local environment variables.")
	}
}

func StringifySendingMessage(msg *message.SendingMessage) (res string) {
	for _, elem := range msg.Elements {
		switch e := elem.(type) {
		case *message.TextElement:
			res += e.Content
		case *message.FaceElement:
			res += "[" + e.Name + "]"
		case *message.MarketFaceElement:
			res += "[" + e.Name + "]"
		case *message.GroupImageElement:
			res += "[Image: " + e.ImageId + "]"
		case *message.AtElement:
			res += e.Display
		case *message.RedBagElement:
			res += "[RedBag:" + e.Title + "]"
		case *message.ReplyElement:
			res += "[Reply:" + strconv.FormatInt(int64(e.ReplySeq), 10) + "]"
		}
	}
	return
}

func FakeUploadImageByte(b []byte) (*message.GroupImageElement, error) {
	return &message.GroupImageElement{
		ImageId: fmt.Sprintf("%d bytes image", len(b)),
		Url:     base64.StdEncoding.EncodeToString(b),
	}, nil
}

func FakeUploadImageUrl(url string) (*message.GroupImageElement, error) {
	return &message.GroupImageElement{
		ImageId: url,
		Url:     url,
	}, nil
}
