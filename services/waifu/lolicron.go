package waifu

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/eric2788/common-utils/request"
)

const lolicronApi = "https://api.lolicon.app/setu/v2?%s"

type (
	Lolicron struct {
	}

	LolicronResp struct {
		Error string `json:"error"`
		Data  []struct {
			Pid        uint64            `json:"pid"`
			Uid        uint64            `json:"uid"`
			Title      string            `json:"title"`
			Author     string            `json:"author"`
			R18        bool              `json:"r18"`
			Width      int               `json:"width"`
			Height     int               `json:"height"`
			Tags       []string          `json:"tags"`
			Ext        string            `json:"ext"`
			AiType     int               `json:"aiType"`
			Urls       map[string]string `json:"urls"`
			UploadDate int64             `json:"uploadDate"`
		} `json:"data"`
	}
)

func (l *Lolicron) GetImages(option *SearchOptions) ([]*ImageData, error) {
	var resp LolicronResp
	r18 := 0
	if option.R18 {
		r18 = 1
	}
	params := &url.Values{
		"tag":     option.Tags,
		"r18":     []string{fmt.Sprint(r18)},
		"num":     []string{fmt.Sprint(option.Amount)},
		"keyword": []string{option.Keyword},
		"size":    []string{"original"},
	}
	err := request.Get(fmt.Sprintf(lolicronApi, params.Encode()), &resp)
	if err != nil {
		return nil, err
	} else if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}

	var results []*ImageData
	for _, data := range resp.Data {

		img, err := request.GetBytesByUrl(data.Urls["original"])
		if err != nil {
			logger.Errorf("尝试下载图源 %s 时出现错误: %v, 将尝试从pixiv下载", data.Urls["original"], err)
			img, err = GetImageFromIllust(data.Pid)
			if err != nil {
				logger.Errorf("从pixiv下载图源 %d 依然失败: %v, 已略过。", data.Pid, err)
			} else {
				logger.Infof("从pixiv下载图源 %d 成功。", data.Pid)
			}
		}

		results = append(results, &ImageData{
			Pid:    data.Pid,
			Uid:    data.Uid,
			Author: data.Author,
			R18:    data.R18,
			Title:  data.Title,
			Tags:   data.Tags,
			Url:    data.Urls["original"],
			Image:  img,
		})
	}

	return results, nil
}
