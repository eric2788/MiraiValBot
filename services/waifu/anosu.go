package waifu

import (
	"fmt"
	"net/url"

	"github.com/eric2788/common-utils/request"
)

const anosuTopURL = "https://image.anosu.top/pixiv/json?%s"

type (
	AnosuTop struct {
	}

	AnosuTopResp struct {
		Pid    uint64   `json:"pid"`
		Uid    uint64   `json:"uid"`
		Title  string   `json:"title"`
		Author string   `json:"author"`
		R18    int      `json:"r18"`
		Width  int      `json:"width"`
		Height int      `json:"height"`
		Tags   []string `json:"tags"`
		Url    string   `json:"url"`
	}
)

func (a *AnosuTop) GetImages(option *SearchOptions) ([]*ImageData, error) {
	var resp []AnosuTopResp
	r18 := 0
	if option.R18 {
		r18 = 1
	}
	// this api only support one keyword
	keyword := option.Keyword
	if keyword == "" && len(option.Tags) > 0 {
		keyword = option.Tags[0]
	}
	params := &url.Values{
		"r18":     []string{fmt.Sprint(r18)},
		"num":     []string{fmt.Sprint(option.Amount)},
		"keyword": []string{keyword},
	}
	err := request.Get(fmt.Sprintf(anosuTopURL, params.Encode()), &resp)
	if err != nil {
		return nil, err
	}
	var results []*ImageData
	for _, data := range resp {

		img, err := request.GetBytesByUrl(data.Url)
		if err != nil {
			logger.Errorf("尝试下载图源 %s 时出现错误: %v, 将尝试从pixiv下载", data.Url, err)
			img, err = GetImageFromIllust(data.Pid)
			if err != nil {
				logger.Errorf("从pixiv下载图源 %d 依然失败: %v, 已略过。", data.Pid, err)
			} else {
				logger.Infof("从pixiv下载图源 %d 成功。", data.Pid)
			}
		}

		results = append(results, &ImageData{
			Title:  data.Title,
			Pid:    data.Pid,
			Uid:    data.Uid,
			Author: data.Author,
			R18:    data.R18 == 1,
			Tags:   data.Tags,
			Url:    data.Url,
			Image:  img,
		})
	}
	return results, nil
}
