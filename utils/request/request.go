package request

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Get(url string, response interface{}) error {

	res, err := http.Get(url)

	if err != nil {
		return err
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
