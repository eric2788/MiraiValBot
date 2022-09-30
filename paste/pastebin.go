package paste

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	PastebinApi = "https://pastebin.com/api/api_post.php"
)

func CreatePasteBin(name, content, lang string) (string, error) {
	devKey := os.Getenv("PASTEBIN_API_KEY")

	if devKey == "" {
		return "", errors.New("env var 'PASTEBIN_API_KEY' is not set")
	}

	postForm := url.Values{}
	postForm.Add("api_option", "paste")
	postForm.Add("api_user_key", "")
	postForm.Add("api_paste_private", "1")
	postForm.Add("api_paste_name", name)
	postForm.Add("api_paste_expire_date", "1D")
	postForm.Add("api_paste_format", "yaml")
	postForm.Add("api_paste_code", content)
	postForm.Add("api_dev_key", devKey)

	res, err := http.Post(PastebinApi, "application/x-www-form-urlencoded", strings.NewReader(postForm.Encode()))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
