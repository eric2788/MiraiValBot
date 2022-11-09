package qq

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/common-utils/request"
)

const cacheDirPath = "cache/"

func saveImages(elements []message.IMessageElement) {
	err := os.MkdirAll(cacheDirPath+"images", os.ModePerm)
	if err != nil {
		logger.Errorf("創建緩存資料夾時出現錯誤: %v", err)
		return
	}

	for _, element := range elements {

		var url string
		var imageId string
		var hash []byte

		switch e := element.(type) {
		case *message.FriendImageElement:
			imageId, hash, url = e.ImageId, e.Md5, e.Url
		case *message.GroupImageElement:
			imageId, hash, url = e.ImageId, e.Md5, e.Url
		case *message.GuildImageElement:
			imageId, hash, url = fmt.Sprint(e.FileId), e.Md5, e.Url
		default:
			continue
		}

		name := hex.EncodeToString(hash)

		b, err := request.GetBytesByUrl(url)
		if err != nil {
			logger.Errorf("下載圖片 %s(%s) 時出現錯誤: %v", imageId, name, err)
			return
		}
		err = os.WriteFile(fmt.Sprintf("%s%s/%s", cacheDirPath, "images", name), b, os.ModePerm)
		if err != nil {
			logger.Errorf("緩存圖片 %s(%s) 時出現錯誤: %v", imageId, name, err)
		} else {
			logger.Infof("緩存圖片 %s(%s) 成功。", imageId, name)
		}
	}
}

func fixGroupImages(gp int64, sending *message.GroupMessage) {
	fixed := make([]message.IMessageElement, len(sending.Elements))
	for _, element := range sending.Elements {
		if groupImage, ok := element.(*message.GroupImageElement); ok {
			name := hex.EncodeToString(groupImage.Md5)
			b, err := os.ReadFile(cacheDirPath + "images/" + name)

			var img *message.GroupImageElement

			if err == nil {
				img, err = NewImagesByByteWithGroup(gp, b)
				if err != nil {
					logger.Errorf("群圖片上傳失敗: %v, 將使用QQ查詢", err)
				}
			} else {
				logger.Errorf("讀取緩存文件 %s 時出現錯誤: %v, 將使用QQ查詢", name, err)
			}

			if img == nil {
				img, err = bot.Instance.QueryGroupImage(gp, groupImage.Md5, groupImage.Size)
				if err != nil {
					logger.Errorf("QQ查詢群圖片失敗: %v, 將繼續使用舊元素發送。", err)
					img = groupImage
				}
			}

			// ensure not null
			if img == nil {
				logger.Warn("檢測到圖片為 null, 將繼續使用舊元素發送。")
				img = groupImage
			}

			fixed = append(fixed, img)
		} else {
			fixed = append(fixed, element)
		}
	}

	sending.Elements = fixed
}
