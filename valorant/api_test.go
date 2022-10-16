package valorant

import (
	"github.com/eric2788/common-utils/datetime"
	"testing"
	"time"

	"github.com/eric2788/common-utils/request"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type Void struct{}

var (
	void              = &Void{}
	allowedStatusCode = map[int]*Void{
		408: void,
		403: void,
		429: void,
		503: void,
	}
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func isAllowedStatus(err error) bool {
	status := 0
	if api, apiOK := err.(*ApiError); apiOK {
		status = api.Status
	} else if http, httpOK := err.(*request.HttpError); httpOK {
		status = http.Code
	}
	_, ok := allowedStatusCode[status]
	if ok {
		logger.Debugf("%d is in allowed status code, skipped", status)
	}
	return ok
}

func formatTime(timeStr string) string {
	if timeStr == "" {
		return "无"
	}
	ti, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		logger.Errorf("无法解析时间: %s, 将返回厡讯息", timeStr)
		return timeStr
	}
	return datetime.FormatISO(ti)
}

func TestGetGameStatus(t *testing.T) {
	status, err := GetGameStatus(AsiaSpecific)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	for _, inc := range status.Data.Incidents {
		t.Log(formatTime(inc.CreatedAt))
	}
	assert.NotEmpty(t, status)
}

func TestGetAccountDetails(t *testing.T) {
	detail, err := GetAccountDetails("勝たんしかrinrin", "JP1")
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, "勝たんしかrinrin", detail.Name)
	assert.Equal(t, "JP1", detail.Tag)
	assert.Equal(t, "2b5e388c-7359-5382-a29f-e5add6e50ed6", detail.PUuid)
	assert.Equal(t, "ap", detail.Region)
}

func TestGetMatchHistories(t *testing.T) {
	histories, err := GetMatchHistoriesAPI("勝たんしかrinrin", "JP1", AsiaSpecific)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(histories))
}

func TestGetMatchDetails(t *testing.T) {
	data, err := GetMatchDetailsAPI("33ae90f4-76b4-4aa0-aa16-331214c7c1dd")
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, "33ae90f4-76b4-4aa0-aa16-331214c7c1dd", data.MetaData.MatchId)
	assert.Equal(t, "ap", data.MetaData.Region)
	assert.Equal(t, "Unrated", data.MetaData.Mode)
	assert.Equal(t, "7a85de9a-4032-61a9-61d8-f4aa2b4a84b6", data.MetaData.SeasonId)
}

func TestGetMMRHistories(t *testing.T) {
	mmrHistories, err := GetMMRHistories("勝たんしかrinrin", "JP1", AsiaSpecific)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, "勝たんしかrinrin", mmrHistories.Name)
	assert.Equal(t, "JP1", mmrHistories.Tag)
}

func TestGetMMRDetailsV1(t *testing.T) {
	mmrDetails, err := GetMMRDetailsV1("勝たんしかrinrin", "JP1", AsiaSpecific)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, "勝たんしかrinrin", mmrDetails.Name)
	assert.Equal(t, "JP1", mmrDetails.Tag)
	assert.NotEmpty(t, mmrDetails.CurrentTierPatched)
}

func TestGetMMRDetailsV2(t *testing.T) {
	mmrDetails, err := GetMMRDetailsV2("勝たんしかrinrin", "JP1", AsiaSpecific)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, "勝たんしかrinrin", mmrDetails.Name)
	assert.Equal(t, "JP1", mmrDetails.Tag)
}

func TestGetMMRDetailsBySeason(t *testing.T) {
	mmrDetails, err := GetMMRDetailsBySeason("勝たんしかrinrin", "JP1", "e3a3", AsiaSpecific)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}
	assert.Equal(t, 7, mmrDetails.Wins)
	assert.Equal(t, 12, mmrDetails.FinalRank)
}
