package qq

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/common-utils/request"
)

const cacheDirPath = "cache/"

func saveGroupImages(msg *message.GroupMessage) {
	err := os.MkdirAll(cacheDirPath+"images", os.ModePerm)
	if err != nil {
		logger.Errorf("創建緩存資料夾時出現錯誤: %v", err)
		return
	}

	for _, element := range msg.Elements {

		var url string
		var imageId string
		var hash []byte

		switch e := element.(type) {
		case *message.FriendImageElement:
			imageId, hash, url = e.ImageId, e.Md5, e.Url
		case *message.GroupImageElement:

			if e.Flash || e.Url == "" {
				if url, err := bot.Instance.GetGroupImageDownloadUrl(e.FileId, msg.GroupCode, e.Md5); err == nil {
					e.Url = url
				} else {
					logger.Errorf("圖片URL為空或是閃照, 但嘗試獲取圖片 %s 的下載URL時出現錯誤: %v", e.FileId, err)
				}
			}

			imageId, hash, url = e.ImageId, e.Md5, e.Url
		case *message.GuildImageElement:
			imageId, hash, url = fmt.Sprint(e.FileId), e.Md5, e.Url
		default:
			continue
		}

		name := hex.EncodeToString(hash)

		b, err := request.GetBytesByUrl(url)
		if err != nil {
			logger.Errorf("下載圖片 %s 時出現錯誤: %v", strings.ToLower(imageId), name, err)
			return
		}
		err = os.WriteFile(fmt.Sprintf("%s%s/%s", cacheDirPath, "images", name), b, os.ModePerm)
		if err != nil {
			logger.Errorf("緩存圖片 %s 時出現錯誤: %v", strings.ToLower(imageId), err)
		} else {
			logger.Infof("緩存圖片 %s 成功。", strings.ToLower(imageId))
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
				} else {
					logger.Infof("恢复缓存图片 %s 成功。", strings.ToLower(groupImage.ImageId))
				}
			} else {

				logger.Errorf("讀取緩存文件 %s 時出現錯誤: %v, 將使用QQ查詢", name, err)

				if url, err := bot.Instance.GetGroupImageDownloadUrl(groupImage.FileId, gp, groupImage.Md5); err == nil {
					logger.Infof("获取群图片下载链接成功，将尝试使用上传通道")
					img, err = NewImageByUrlWithGroup(gp, url)
					if err == nil {
						logger.Infof("群图片上传成功")
					} else {
						logger.Warnf("群图片上传失败: %v", err)
					}
				}
			}

			if img == nil {
				img, err = bot.Instance.QueryGroupImage(gp, groupImage.Md5, groupImage.Size)
				if err != nil {
					logger.Errorf("QQ查詢群圖片失敗: %v, 將繼續使用舊元素發送。", err)
					img = groupImage
				} else {
					logger.Infof("查询图片 %s 成功。", strings.ToLower(groupImage.ImageId))

					//查詢成功后下載
					url := img.Url
					b, err := request.GetBytesByUrl(url)
					if err != nil {
						logger.Errorf("下載查詢圖片 %s 時出現錯誤: %v", strings.ToLower(groupImage.ImageId), name, err)
					} else {
						err = os.WriteFile(fmt.Sprintf("%s%s/%s", cacheDirPath, "images", name), b, os.ModePerm)
						if err != nil {
							logger.Errorf("緩存查詢圖片 %s 時出現錯誤: %v", strings.ToLower(groupImage.ImageId), err)
						} else {
							logger.Infof("緩存查詢圖片 %s 成功。", strings.ToLower(groupImage.ImageId))
						}
					}
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
