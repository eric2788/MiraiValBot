package valorant

import (
	"errors"
	"fmt"
	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
	"net/http"
	"os"
	"time"
)

type Region string

const (
	BaseUrl = "https://api.henrikdev.xyz/valorant"
	V1      = "/v1"
	V2      = "/v2"
	V3      = "/v3"

	Europe       Region = "eu"
	NorthAmerica Region = "na"
	AsiaSpecific Region = "ap"
	Korea        Region = "kr"
	LatinAmerica Region = "latam"
	Brazil       Region = "br"
)

var (
	client = &http.Client{Timeout: time.Minute}
)

func doRequest(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", uarand.GetRandom())
	if apiKey := os.Getenv("HENRIK_VALORANT_API_KEY"); apiKey != "" {
		req.Header.Set("Authorization", apiKey)
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	} else if res.StatusCode != 200 {
		return nil, &request.HttpError{
			Code:     res.StatusCode,
			Status:   res.Status,
			Response: res,
		}
	}
	return res, nil
}

func getRequest(path string) (*Resp, error) {
	req, err := http.NewRequest(http.MethodGet, BaseUrl+path, nil)
	if err != nil {
		return nil, err
	}
	res, err := doRequest(req)
	if err != nil {
		return nil, err
	}
	var data Resp
	err = request.Read(res, &data)
	if err != nil {
		return nil, err
	} else if len(data.Errors) > 0 {
		return nil, &ApiError{data.Errors}
	} else if data.Status != 200 {
		return nil, errors.New(fmt.Sprintf("status code %v", data.Status))
	}
	return &data, nil
}

func GetAccountDetails(name, tag string) (*AccountDetails, error) {
	resp, err := getRequest(fmt.Sprintf("%v/account/%s/%s", V1, name, tag))
	if err != nil {
		return nil, err
	}
	accountDetails := &AccountDetails{}
	err = resp.ParseData(accountDetails)
	return accountDetails, err
}

func GetMatchHistories(name, tag string, region Region) ([]MatchData, error) {
	resp, err := getRequest(fmt.Sprintf("%v/matches/%s/%s/%s", V3, region, name, tag))
	if err != nil {
		return nil, err
	}
	var matchHistories []MatchData
	err = resp.ParseData(&matchHistories)
	return matchHistories, err
}

func GetMatchDetails(matchId string) (*MatchData, error) {
	resp, err := getRequest(fmt.Sprintf("%v/match/%s", V2, matchId))
	if err != nil {
		return nil, err
	}
	matchDetails := &MatchData{}
	err = resp.ParseData(matchDetails)
	return matchDetails, err
}

func GetLocalizedContent() (*Localization, error) {
	resp, err := getRequest(fmt.Sprintf("%v/content", V1))
	if err != nil {
		return nil, err
	}
	localization := &Localization{}
	err = resp.ParseData(localization)
	return localization, err
}

func GetGameStatus(region Region) (*GameStatus, error) {
	resp, err := getRequest(fmt.Sprintf("%v/status/%s", V1, region))
	if err != nil {
		return nil, err
	}
	gameStatus := &GameStatus{}
	err = resp.ParseData(gameStatus)
	return gameStatus, err
}
