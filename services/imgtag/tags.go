package imgtag

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/huggingface"
	"github.com/eric2788/MiraiValBot/utils/misc"
	"github.com/eric2788/common-utils/request"
)

const tagURL = "https://nsfwtag.azurewebsites.net/api/tag?limit=%f&url=%s"

var logger = utils.GetModuleLogger("service.imgtag")

type ImageTagger func(imgUrl string, confidence float64) (map[string]float64, error)

var taggerProviders = map[string]ImageTagger{
	"hfspace2": getTagFromHFSpace2,
	"hfspace":  getTagFromHFSpace,
	"azure":    getTagFromAzure,
}

func GetTagsFromImage(imgUrl string) ([]string, error) {
	var dict map[string]float64
	var err error

	for name, provider := range taggerProviders {

		dict, err = provider(imgUrl, file.DataStorage.Setting.TagClassifyLimit)
		if err != nil {
			logger.Errorf("从 %s 获取图片鉴别标签错误: %v, 将使用下一个API", name, err)
		} else {
			break
		}
	}

	if err != nil {
		return nil, err
	}

	var tags []string
	for key := range dict {
		if strings.HasPrefix(key, "rating:") {
			continue
		} else { //filter rating:xxx
			tags = append(tags, strings.ReplaceAll(key, "_", " "))
		}
	}
	return tags, nil
}

func getTagFromAzure(imgUrl string, confidence float64) (map[string]float64, error) {
	var dict map[string]float64
	err := request.Get(fmt.Sprintf(tagURL, confidence, url.QueryEscape(imgUrl)), &dict)
	return dict, err
}

// getTagFromHFSpace using model: mayhug-rainchan-anime-image-label
func getTagFromHFSpace(imgUrl string, confidence float64) (map[string]float64, error) {

	b64, t, err := misc.ReadURLToSrcData(imgUrl)

	if err != nil {
		return nil, err
	} else if !strings.HasPrefix(t, "image/") {
		return nil, fmt.Errorf("url is not image type")
	}

	api := huggingface.NewSpaceApi("mayhug-rainchan-anime-image-label",
		b64,
		confidence,
		"ResNet50",
	)
	return api.EndPoint("api/predict/").GetClassifiedLabels()
}

// getTagFromHFSpace2 using model: hysts-deepdanbooru
func getTagFromHFSpace2(imgUrl string, confidence float64) (tags map[string]float64, err error) {
	b64, t, err := misc.ReadURLToSrcData(imgUrl)
	if err != nil {
		return nil, err
	} else if !strings.HasPrefix(t, "image/") {
		return nil, fmt.Errorf("url is not image type")
	}

	api := huggingface.NewSpaceApi("hysts-deepdanbooru", b64, confidence)
	tags, err = api.EndPoint("api/predict/").GetClassifiedLabels()
	if err != nil {
		return
	}
	for tag := range tags {
		if strings.HasPrefix(tag, "rating:") {
			delete(tags, tag)
		}
	}
	return
}
