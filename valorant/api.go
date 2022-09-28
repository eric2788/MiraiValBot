package valorant

import (
	"errors"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/corpix/uarand"
	"github.com/eric2788/common-utils/request"
	"github.com/eric2788/common-utils/set"
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
	AllowedSeasons = set.FromStrArr([]string{
		"e5a3",
		"e5a2",
		"e5a1",
		"e4a3",
		"e4a2",
		"e4a1",
		"e3a3",
		"e3a2",
		"e3a1",
		"e2a3",
		"e2a2",
		"e2a1",
		"e1a3",
		"e1a2",
		"e1a1",
	})
	client = &http.Client{Timeout: time.Minute}
	logger = utils.GetModuleLogger("valorant.api")
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
	var data Resp
	err := getRequestCustom(path, &data)
	if err != nil {
		return nil, err
	} else if len(data.Errors) > 0 {
		return nil, &ApiError{data.Errors}
	} else if data.Status != 200 {
		return nil, errors.New(fmt.Sprintf("status: %v", data.Status))
	}
	return &data, nil
}

func getRequestCustom(path string, response interface{}) error {
	logger.Debugf("preparing to do get request: %v", BaseUrl+path)
	req, err := http.NewRequest(http.MethodGet, BaseUrl+path, nil)
	if err != nil {
		return err
	}
	res, err := doRequest(req)
	if err != nil {
		return err
	}
	logger.Debugf("response status for %v: %v", BaseUrl+path, res.Status)
	return request.Read(res, response)
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

func GetLocalizedContentByLocale(locale string) (*Localization, error) {
	resp, err := getRequest(fmt.Sprintf("%v/content?locale=%s", V1, locale))
	if err != nil {
		return nil, err
	}
	localization := &Localization{}
	err = resp.ParseData(localization)
	return localization, err
}

func GetGameStatus(region Region) (*StatusResp, error) {
	var data = &StatusResp{}
	err := getRequestCustom(fmt.Sprintf("%v/status/%s", V1, region), data)
	if err != nil {
		return nil, err
	} else if len(data.Errors) > 0 {
		return nil, &ApiError{data.Errors}
	} else if data.Status != 200 {
		return nil, errors.New(fmt.Sprintf("status code %v", data.Status))
	}
	return data, err
}

func GetMMRHistories(name, tag string, region Region) (*PlayerInfoResp, error) {
	var data = &PlayerInfoResp{}
	err := getRequestCustom(fmt.Sprintf("%v/mmr-history/%s/%s/%s", V1, region, name, tag), data)
	if err != nil {
		return nil, err
	} else if len(data.Errors) > 0 {
		return nil, &ApiError{data.Errors}
	} else if data.Status != 200 {
		return nil, errors.New(fmt.Sprintf("status code %v", data.Status))
	}
	return data, err
}

func GetMMRDetailsV1(name, tag string, region Region) (*MMRV1Details, error) {
	resp, err := getRequest(fmt.Sprintf("%v/mmr/%s/%s/%s", V1, region, name, tag))
	if err != nil {
		return nil, err
	}
	mmrDetails := &MMRV1Details{}
	err = resp.ParseData(mmrDetails)
	return mmrDetails, err
}

func GetMMRDetailsV2(name, tag string, region Region) (*MMRV2Details, error) {
	resp, err := getRequest(fmt.Sprintf("%v/mmr/%s/%s/%s", V2, region, name, tag))
	if err != nil {
		return nil, err
	}
	mmrDetails := &MMRV2Details{}
	err = resp.ParseData(mmrDetails)
	return mmrDetails, err
}

func GetMMRDetailsBySeason(name, tag, filter string, region Region) (*MMRV2SeasonDetails, error) {
	resp, err := getRequest(fmt.Sprintf("%v/mmr/%s/%s/%s?filter=%s", V2, region, name, tag, filter))
	if err != nil {
		return nil, err
	}
	mmrDetails := &MMRV2SeasonDetails{}
	err = resp.ParseData(mmrDetails)
	return mmrDetails, err
}
