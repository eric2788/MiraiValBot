package valorant

import (
	"fmt"
	"testing"
)

func TestCache(t *testing.T) {

	datas := []MatchData{
		generateMatchData("aaa"),
		generateMatchData("bbb"),
		generateMatchData("ccc"),
		generateMatchData("dddd"),
		generateMatchData("eee"),
	}

	for _, data := range datas {
		printMatchId(&data) // add go will make all become last value
	}
}

func printMatchId(data *MatchData) {
	fmt.Println(data.MetaData.MatchId)
}

func generateMatchData(id string) MatchData {
	return MatchData{
		MetaData: MatchMetaData{
			MatchId: id,
		},
	}
}
