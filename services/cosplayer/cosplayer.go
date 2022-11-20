package cosplayer

import (
	"errors"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/common-utils/request"
	"golang.org/x/exp/maps"
	"math/rand"
	"time"
)

type (
	WebAPI interface {
		GetImages() (*Data, error)
	}

	Data struct {
		Title string
		Urls  []string
	}
)

var (
	logger    = utils.GetModuleLogger("services.cosplayer")
	providers = map[string]WebAPI{
		"ovooa": &OVOOA{},
	}
	singleImageURLs = []string{
		"https://picture.yinux.workers.dev",
		"https://api.jrsgslb.cn/cos/url.php?return=img",
	}
)

func GetImageRandom() ([]byte, error) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(singleImageURLs), func(i, j int) {
		singleImageURLs[i], singleImageURLs[j] = singleImageURLs[j], singleImageURLs[i]
	})
	for _, url := range singleImageURLs {
		if data, err := request.GetBytesByUrl(url); err == nil {
			return data, nil
		} else {
			logger.Errorf("使用 %s 獲取 cosplayer 圖片失敗: %v, 將使用下一個API", url, err)
		}
	}
	return nil, errors.New("沒有可用的 cosplay 圖片 API")
}

func GetImagesRandom() (*Data, error) {
	rand.Seed(time.Now().UnixNano())
	keys := maps.Keys(providers)

	rand.Shuffle(len(keys), func(i, j int) {
		keys[i], keys[j] = keys[j], keys[i]
	})

	for _, name := range keys {
		p := providers[name]
		if data, err := p.GetImages(); err == nil {
			return data, nil
		} else {
			logger.Errorf("使用 %s 獲取 cosplayer 圖片失敗: %v, 將使用下一個API", name, err)
		}
	}
	return nil, errors.New("沒有可用的 cosplay 圖片 API")
}

func GetImages(provider string) (*Data, error) {
	if p, ok := providers[provider]; ok {
		return p.GetImages()
	}
	return nil, errors.New("沒有可用的 cosplay 圖片 API")
}
