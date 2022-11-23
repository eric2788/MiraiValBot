package waifu

import (
	"errors"
	"fmt"
	"github.com/eric2788/common-utils/request"
	"net/url"
	"os"
	"strings"
	"time"
)

const danbooruURL = "https://danbooru.donmai.us"

const (
	RatingQuestionable = "q"
	RatingExplicit     = "e"
	RatingSafe         = "s"
)

type (
	Rating string

	Danbooru struct {
	}

	DanbooruPostResp struct {
		Id                 uint64    `json:"id"`
		CreatedAt          time.Time `json:"created_at"`
		UploaderId         uint64    `json:"uploader_id"`
		Source             string    `json:"source"`
		Score              int       `json:"score"`
		Md5                string    `json:"md5"`
		Rating             Rating    `json:"rating"`
		ImageWidth         int       `json:"image_width"`
		ImageHeight        int       `json:"image_height"`
		TagString          string    `json:"tag_string"`
		FileExt            string    `json:"file_ext"`
		TagCount           int       `json:"tag_count"`
		PixivId            uint64    `json:"pixiv_id"`
		TagStringGeneral   string    `json:"tag_string_general"`
		TagStringCharacter string    `json:"tag_string_character"`
		TagStringArtist    string    `json:"tag_string_artist"`
		TagStringMeta      string    `json:"tag_string_meta"`
		FileUrl            string    `json:"file_url"`
		LargeFileUrl       string    `json:"large_file_url"`
		PreviewFileUrl     string    `json:"preview_file_url"`
	}

	DanbooruErrResp struct {
		Success   bool     `json:"success"`
		SError    string   `json:"error"`
		Message   string   `json:"message"`
		BackTrace []string `json:"backtrace"`
	}
)

func (d DanbooruErrResp) Error() string {
	return d.Message
}

func (d *Danbooru) GetImages(option *SearchOptions) ([]*ImageData, error) {
	apiKey, login := os.Getenv("DANBOORU_API_KEY"), os.Getenv("DANBOORU_LOGIN")
	if apiKey == "" || login == "" {
		return nil, errors.New("未设置DANBOORU_API_KEY或DANBOORU_LOGIN")
	}

	// convert tags to danbooru format
	tags := make([]string, len(option.Tags))
	for i, tag := range option.Tags {
		tags[i] = strings.ReplaceAll(tag, " ", "_")
	}

	var resp []DanbooruPostResp
	params := &url.Values{
		"login":   []string{login},
		"api_key": []string{apiKey},
		"limit":   []string{fmt.Sprint(option.Amount)},
		"tags":    tags,
		"random":  []string{"true"},
	}
	err := request.Get(fmt.Sprintf("%s/posts.json?%s", danbooruURL, params.Encode()), &resp)
	if err != nil {
		if httpErr, ok := err.(*request.HttpError); ok {
			var errResp DanbooruErrResp
			defer httpErr.Response.Body.Close()
			if derr := request.Read(httpErr.Response, &errResp); derr == nil {
				return nil, errResp
			}
		}
		return nil, err
	}
	var images map[uint64]*ImageData

	for _, post := range resp {

		r18 := d.isR18(post)

		if r18 && !option.R18 {
			continue
		}

		images[post.Id] = &ImageData{
			Title:  post.TagStringCharacter,
			Url:    post.LargeFileUrl,
			Pid:    post.PixivId,
			Uid:    post.UploaderId,
			R18:    d.isR18(post),
			Author: fmt.Sprint(post.UploaderId),
			Tags:   strings.Split(post.TagString, " "),
		}
	}

	var results []*ImageData
	for _, image := range images {
		results = append(results, image)
	}

	if option.Amount > len(results) {
		ids, err := d.GetImages(NewOptions(
			WithTags(tags...),
			WithAmount(option.Amount-len(results)),
			WithR18(option.R18),
		))
		if err != nil {
			return nil, err
		}

		results = append(results, ids...)
	}

	return results, nil
}

func (d *Danbooru) isR18(post DanbooruPostResp) bool {
	return post.Rating != RatingSafe
}
