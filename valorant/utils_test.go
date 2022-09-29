package valorant

import (
	"fmt"
	"testing"
)

func TestGetDeathMatchRanking(t *testing.T) {
	match, err := GetMatchDetails("a4e99fec-647d-4a15-9015-967c8e29355a")
	if err != nil {
		t.Fatal(err)
	}
	players := GetDeathMatchRanking(match)
	for i, player := range players {
		t.Logf("%d. %s (%d kills)(score: %d)", 
			i+1, 
			fmt.Sprintf("%s#%s",player.Name, player.Tag), 
			player.Stats.Kills, 
			player.Stats.Score,
		)
	}
}