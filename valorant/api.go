package valorant

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/corpix/uarand"
	"github.com/eric2788/MiraiValBot/redis"
	"github.com/eric2788/common-utils/request"
	"github.com/eric2788/common-utils/set"
)

type Region string

const (
	HenrikBaseUrl = "https://api.henrikdev.xyz/valorant"
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
		return nil, &ApiError{
			Status: data.Status,
			Errors: data.Errors,
		}
	}
	return &data, nil
}

func getRequestCustom(path string, response interface{}) error {
	logger.Debugf("preparing to do get request: %v", HenrikBaseUrl+path)
	req, err := http.NewRequest(http.MethodGet, HenrikBaseUrl+path, nil)
	if err != nil {
		return err
	}
	res, err := doRequest(req)
	if err != nil {
		if httpErr, ok := err.(*request.HttpError); ok {
			if err := request.Read(httpErr.Response, response); err == nil {
				return nil
			}else{
				logger.Warnf("cannot parse http error response to Resp: %v, use back http error as error.", err)
				if b, err := io.ReadAll(httpErr.Response.Body); err == nil {
					logger.Debugf("response body: %q", string(b))
				}else{
					logger.Debugf("cannot print response body: %v", err)
				}
				
			}
		}
		return err
	}
	logger.Debugf("response status for %v: %v", HenrikBaseUrl+path, res.Status)
	return request.Read(res, response)
}

func UpdateAccountDetails(name, tag string) error {
	_, err := getRequest(fmt.Sprintf("%v/account/%s/%s?force=true", V1, name, tag))
	return err
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

func GetMatchHistoriesAPI(name, tag string, region Region) ([]MatchData, error) {
	resp, err := getRequest(fmt.Sprintf("%v/matches/%s/%s/%s", V3, region, name, tag))
	if err != nil {
		return nil, err
	}
	var matchHistories []MatchData
	err = resp.ParseData(&matchHistories)
	return matchHistories, err
}

func GetMatchHistories(name, tag string, region Region) ([]MatchData, error) {
	matchHistories, err := GetMatchHistoriesAPI(name, tag, region)
	if err == nil {
		go cacheMatchHistories(matchHistories)
	}
	return matchHistories, err
}

func GetMatchDetailsAPI(matchId string) (*MatchData, error) {
	resp, err := getRequest(fmt.Sprintf("%v/match/%s", V2, matchId))
	if err != nil {
		return nil, err
	}
	matchDetails := &MatchData{}
	err = resp.ParseData(matchDetails)
	return matchDetails, err
}

func GetMatchDetails(matchId string) (*MatchData, error) {
	matchDetails := &MatchData{}

	if exist, err := redis.Get(matchId, matchDetails); exist {
		return matchDetails, nil
	} else if err != nil {
		logger.Warnf("从 redis 提取快取时出现错误: %v, 将使用 HTTP 请求.", err)
	}

	matchDetails, err := GetMatchDetailsAPI(matchId)
	if err == nil {
		go cacheMatchData(matchDetails)
	}
	return matchDetails, err
}

func GetLocalizedContent() (Localization, error) {
	var local Localization
	err := getRequestCustom(fmt.Sprintf("%v/content", V1), &local)
	if err != nil {
		return nil, err
	}
	return local, err
}

func GetLocalizedContentByLocale(locale string) (Localization, error) {
	var localization Localization
	err := getRequestCustom(fmt.Sprintf("%v/content?locale=%s", V1, locale), &localization)
	if err != nil {
		return nil, err
	}
	return localization, err
}

func GetGameStatus(region Region) (*StatusResp, error) {
	var data = &StatusResp{}
	err := getRequestCustom(fmt.Sprintf("%v/status/%s", V1, region), data)
	if err != nil {
		return nil, err
	} else if len(data.Errors) > 0 {
		return nil, &ApiError{data.Status, data.Errors}
	} else if data.Status != 200 {
		return nil, fmt.Errorf("status code %v", data.Status)
	}
	return data, err
}

func GetMMRHistories(name, tag string, region Region) (*PlayerInfoResp, error) {
	var data = &PlayerInfoResp{}
	err := getRequestCustom(fmt.Sprintf("%v/mmr-history/%s/%s/%s", V1, region, name, tag), data)
	if err != nil {
		return nil, err
	} else if len(data.Errors) > 0 {
		return nil, &ApiError{data.Status, data.Errors}
	} else if data.Status != 200 {
		return nil, fmt.Errorf("status code %v", data.Status)
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
	if !AllowedSeasons.Contains(filter) {
		return nil, fmt.Errorf("无效的赛季: %v, 可用的赛季: %s", filter, strings.Join(AllowedSeasons.ToArr(), ", "))
	}
	resp, err := getRequest(fmt.Sprintf("%v/mmr/%s/%s/%s?filter=%s", V2, region, name, tag, filter))
	if err != nil {
		return nil, err
	}
	mmrDetails := &MMRV2SeasonDetails{}
	err = resp.ParseData(mmrDetails)
	return mmrDetails, err
}
