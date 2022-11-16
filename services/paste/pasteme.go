package paste

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/corpix/uarand"
)

const (
	PasteMeBaseURL = "https://pasteme.cn"
)

var (
	client = &http.Client{
		Jar: createCookieJar(),
	}
	logger = utils.GetModuleLogger("paste.me")
)

type Resp struct {
	Code int    `json:"code"`
	Key  string `json:"key"`
}

type ErrResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err *ErrResp) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}

func createCookieJar() http.CookieJar {
	jar, _ := cookiejar.New(nil)
	return jar
}

func CreatePasteMe(lang, content string) (string, error) {

	userAgent := uarand.GetRandom()
	if err := browseMainPage(userAgent); err != nil {
		return "", err
	}

	field := map[string]interface{}{
		"content":       content,
		"lang":          lang,
		"password":      "",
		"expire_count":  1,
		"expire_second": 300,
		"self_destruct": true,
	}

	body, err := json.Marshal(field)

	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, PasteMeBaseURL+"/api/v3/paste", bytes.NewReader(body))

	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")

	if url, err := url.Parse(PasteMeBaseURL); err == nil && url != nil {
		logger.Debugf("cookies of api client: %+v", client.Jar.Cookies(url))
	} else {
		logger.Debugf("Error while parsing url %s: %s", PasteMeBaseURL, err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if res.StatusCode == 201 {
		var resp Resp
		err = json.Unmarshal(b, &resp)
		return fmt.Sprintf("%s#%s", PasteMeBaseURL, resp.Key), err
	} else {
		var errResp ErrResp
		if err = json.Unmarshal(b, &errResp); err == nil {
			return "", &errResp
		} else {
			return "", errors.New(fmt.Sprintf("%d: %s", res.StatusCode, res.Status))
		}
	}
}

func getCookiesFromPage(userAgent, url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	logger.Debugf("get cookies from url %s: %+v", req.URL.String(), res.Cookies())
	client.Jar.SetCookies(req.URL, res.Cookies())
	return nil
}

// to get cookie
func browseMainPage(userAgent string) error {
	urls := []string{
		PasteMeBaseURL,
		PasteMeBaseURL + "/api/v3/?method=beat",
		PasteMeBaseURL + "/?encode=text",
	}

	for _, url := range urls {
		if err := getCookiesFromPage(userAgent, url); err != nil {
			return err
		}
	}

	return nil
}
