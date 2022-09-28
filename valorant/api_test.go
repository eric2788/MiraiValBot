package valorant

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func TestGetAccountDetails(t *testing.T) {
	detail, err := GetAccountDetails("勝たんしかrinrin", "JP1")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "勝たんしかrinrin", detail.Name)
	assert.Equal(t, "JP1", detail.Tag)
	assert.Equal(t, "2b5e388c-7359-5382-a29f-e5add6e50ed6", detail.PUuid)
	assert.Equal(t, "ap", detail.Region)
}

func TestGetMatchHistories(t *testing.T) {
	histories, err := GetMatchHistories("勝たんしかrinrin", "JP1", AsiaSpecific)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(histories))
}

func TestGetMatchDetails(t *testing.T) {
	data, err := GetMatchDetails("33ae90f4-76b4-4aa0-aa16-331214c7c1dd")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "33ae90f4-76b4-4aa0-aa16-331214c7c1dd", data.MetaData.MatchId)
	assert.Equal(t, "ap", data.MetaData.Region)
	assert.Equal(t, "Unrated", data.MetaData.Mode)
	assert.Equal(t, "7a85de9a-4032-61a9-61d8-f4aa2b4a84b6", data.MetaData.SeasonId)
}
