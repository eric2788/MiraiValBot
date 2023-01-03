package waifu

import (
	"fmt"
	"io"
	"net/http"

	"github.com/corpix/uarand"
)

func tryGetImage(images ...string) string {
	for _, img := range images {
		if img != "" {
			return img
		}
	}
	return ""
}

func getImageByte(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Referer", "https://pixiv.net")
	req.Header.Set("User-Agent", uarand.GetRandom())

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	} else if res.StatusCode != 200 {
		return nil, fmt.Errorf(res.Status)
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func GetImageFromIllust(id uint64) ([]byte, error) {
	data, err := getIllust(id)
	if err != nil {
		return nil, err
	}
	imgUrl := tryGetImage(
		data.Images.Original,
		data.Images.Large,
		data.Images.Medium,
		data.Images.SquareMedium,
	)
	if imgUrl == "" {
		return nil, fmt.Errorf("插画 %d 的图源为空", id)
	}
	return getImageByte(imgUrl)
}
