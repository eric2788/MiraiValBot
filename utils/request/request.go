package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpError struct {
	Code     int
	Status   string
	Response *http.Response
}

func (e HttpError) Error() string {
	return fmt.Sprintf("%v: %s", e.Code, e.Status)
}

func Get(url string, response interface{}) error {

	res, err := http.Get(url)

	if err != nil {
		return err
	} else if res.StatusCode != 200 {
		return &HttpError{
			Code:     res.StatusCode,
			Status:   res.Status,
			Response: res,
		}
	}

	return Read(res, response)
}

func GetHtml(url string) (string, error) {

	res, err := http.Get(url)

	if err != nil {
		return "", err
	}

	return ReadString(res)
}

func GetBytesByUrl(url string) (img []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = res.Body.Close()
	}()
	img, err = ioutil.ReadAll(res.Body)
	return
}

func ReadString(res *http.Response) (string, error) {
	var err error

	defer func() {
		err = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Read(res *http.Response, response interface{}) error {

	var err error

	defer func() {
		err = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, response)
	return err
}
