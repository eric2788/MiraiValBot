package waifu

import (
	"fmt"
	"net/url"

	"github.com/eric2788/common-utils/request"
)

const yubanURL = "https://setu.yuban10703.xyz/setu?%s"

type (
	Yuban struct {
	}

	YubanResp struct {
		Detail string   `json:"detail"`
		Count  int      `json:"count"`
		Tags   []string `json:"tags"`
		Data   []struct {
			ArtWork struct {
				Title string `json:"title"`
				Id    uint64 `json:"id"`
			} `json:"artwork"`
			Author struct {
				Name string `json:"name"`
				Id   uint64 `json:"id"`
			} `json:"author"`
			SanityLevel int    `json:"sanity_level"`
			R18         bool   `json:"r18"`
			CreateDate  string `json:"create_date"`
			Size        struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"size"`
			Tags []string `json:"tags"`
			Urls struct {
				Original string `json:"original"`
				Large    string `json:"large"`
				Medium   string `json:"medium"`
			} `json:"urls"`
		}
	}
)

func (y *Yuban) GetImages(option *SearchOptions) ([]*ImageData, error) {
	var resp YubanResp
	r18 := 0
	if option.R18 {
		r18 = 1
	}
	tags := append(option.Tags, option.Keyword)
	params := &url.Values{
		"r18":  []string{fmt.Sprint(r18)},
		"num":  []string{fmt.Sprint(option.Amount)},
		"tags": tags,
	}
	err := request.Get(fmt.Sprintf(yubanURL, params.Encode()), &resp)
	if err != nil {
		return nil, err
	}
	var results []*ImageData
	for _, data := range resp.Data {

		url := tryGetImage(
			data.Urls.Original,
			data.Urls.Large,
			data.Urls.Medium,
		)

		img, err := getImageByte(url)

		if err != nil {
			logger.Errorf("尝试下载图源 %s 时出现错误: %v, 将尝试从pixiv下载", url, err)
			img, err = GetImageFromIllust(data.ArtWork.Id)
			if err != nil {
				logger.Errorf("从pixiv下载图源 %d 依然失败: %v, 已略过。", data.ArtWork.Id, err)
			} else {
				logger.Infof("从pixiv下载图源 %d 成功。", data.ArtWork.Id)
			}
		}

		results = append(results, &ImageData{
			Title:  data.ArtWork.Title,
			Pid:    data.ArtWork.Id,
			Uid:    data.Author.Id,
			Author: data.Author.Name,
			R18:    data.R18,
			Tags:   data.Tags,
			Url:    url,
			Image:  img,
		})
	}

	return results, nil
}
