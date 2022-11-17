package waifu

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
	"github.com/everpcpc/pixiv"
	"github.com/lucas-clemente/quic-go/http3"
)

const (
	pixivMoeApi    = "https://api.pixivel.moe/v2/pixiv/illust/search/%s?page=%d&sortpop=true"
	pixivMoeTagApi = "https://api.pixivel.moe/v2/pixiv/tag/search/%s?page=%d&sortpop=true"
)

var http3Client = http.Client{
	Transport: &http3.RoundTripper{},
}

type (
	PixelMoe struct {
	}

	PixivMoeResp struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
		Data    struct {
			HasNext bool `json:"has_next"`
			Illusts []struct {
				ID        uint64        `json:"id"`
				Title     string        `json:"title"`
				AltTitle  string        `json:"altTitle"`
				Sanity    int           `json:"sanity"`
				Width     int           `json:"width"`
				Height    int           `json:"height"`
				Tags      []PixivMoeTag `json:"tags"`
				Statistic struct {
					Views     int `json:"views"`
					Likes     int `json:"likes"`
					Bookmarks int `json:"bookmarks"`
				} `json:"statistic"`
			} `json:"illusts"`
		} `json:"data"`
	}

	PixivMoeTag struct {
		Name        string `json:"name"`
		Translation string `json:"translation"`
	}
)

func (p *PixelMoe) GetImages(option *SearchOptions) ([]ImageData, error) {
	var ids []uint64
	var err error
	if option.Keyword != "" {
		ids, err = p.getPixivIdsByKeyword(option.Keyword, 0, option.Amount, option.R18)
	} else if len(option.Tags) > 0 {
		ids, err = p.getPixivIdsByTags(option.Tags, 0, option.Amount, option.R18)
	} else {
		return nil, fmt.Errorf("unknown search option")
	}
	if err != nil {
		return nil, err
	}
	var results []ImageData
	for _, id := range ids {
		data, err := getIllust(id)
		if err != nil {
			logger.Errorf("獲取 pixiv 圖片 %d 失敗: %v", id, err)
			continue
		} else if data == nil || data.Images == nil {
			logger.Warnf("獲取 pixiv 圖片 %d 失敗: 無法獲取圖片資訊", id)
			continue
		}

		imgUrl := p.tryGetImage(data.Images)

		if imgUrl == "" {
			logger.Warnf("所有圖像ID %d 的網址為空: %v", id, data.Images)
			continue
		}

		results = append(results, ImageData{
			Pid:    data.ID,
			Uid:    data.User.ID,
			R18:    p.checkTagIsR18(data.Tags),
			Author: data.User.Name,
			Title:  data.Title,
			Url:    imgUrl,
			Tags:   p.toArr(data.Tags),
		})

	}

	// 所有pixiv图片均无法获取
	if len(results) == 0 && len(ids) > 0 {
		return nil, fmt.Errorf("无法索取所有 pixiv 图片, 原因请查看服务端输出")
	}

	return results, nil
}

func (p *PixelMoe) tryGetImage(images *pixiv.Images) string {
	if images.Original != "" {
		return images.Original
	} else if images.Large != "" {
		return images.Large
	} else if images.Medium != "" {
		return images.Medium
	} else if images.SquareMedium != "" {
		return images.SquareMedium
	}
	return ""
}

func (p *PixelMoe) toArr(tags []pixiv.Tag) []string {
	var results []string
	for _, tag := range tags {
		results = append(results, tag.Name)
	}
	return results
}

func (p *PixelMoe) checkTagIsR18(tags []pixiv.Tag) bool {
	for _, tag := range tags {
		if tag.Name == "R-18" {
			return true
		}
	}
	return false
}

func (p *PixelMoe) _checkTagIsR18(tags []PixivMoeTag) bool {
	for _, tag := range tags {
		if tag.Name == "R-18" {
			return true
		}
	}
	return false
}

func (p *PixelMoe) httpGet(url string, response interface{}) error {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", uarand.GetRandom())
	req.Header.Set("Referer", "https://pixivel.moe/")
	res, err := http3Client.Do(req)

	if err != nil {
		return err
	} else if res.StatusCode != 200 {
		return &request.HttpError{
			Code:     res.StatusCode,
			Status:   res.Status,
			Response: res,
		}
	}

	return request.Read(res, response)
}

func (p *PixelMoe) getPixivIdsByTags(tags []string, page, amount int, r18 bool) ([]uint64, error) {
	var resp PixivMoeResp
	err := p.httpGet(fmt.Sprintf(pixivMoeTagApi, url.QueryEscape(strings.Join(tags, ",")), page), &resp)
	if err != nil {
		return nil, err
	} else if resp.Error {
		return nil, fmt.Errorf(resp.Message)
	}
	var results []uint64
	for _, illust := range resp.Data.Illusts {
		if p._checkTagIsR18(illust.Tags) && !r18 {
			continue
		}
		results = append(results, illust.ID)
		// 获取特定数量
		if len(results) >= amount {
			break
		}
	}
	if amount > len(results) && resp.Data.HasNext {
		ids, err := p.getPixivIdsByTags(tags, page+1, amount-len(results), r18)
		if err != nil {
			return nil, err
		}
		results = append(results, ids...)
	}
	return results, nil
}

func (p *PixelMoe) getPixivIdsByKeyword(keyword string, page, amount int, r18 bool) ([]uint64, error) {
	var resp PixivMoeResp
	err := p.httpGet(fmt.Sprintf(pixivMoeApi, url.QueryEscape(keyword), page), &resp)
	if err != nil {
		return nil, err
	} else if resp.Error {
		return nil, fmt.Errorf(resp.Message)
	}
	var results []uint64
	for _, illust := range resp.Data.Illusts {
		if p._checkTagIsR18(illust.Tags) && !r18 {
			continue
		}
		results = append(results, illust.ID)
		// 获取特定数量
		if len(results) >= amount {
			break
		}
	}
	if amount > len(results) && resp.Data.HasNext {
		ids, err := p.getPixivIdsByKeyword(keyword, page+1, amount-len(results), r18)
		if err != nil {
			return nil, err
		}
		results = append(results, ids...)
	}
	return results, nil
}
