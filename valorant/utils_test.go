package valorant

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

	assert.Equal(t, 12, len(players))
}

func TestPercentageDisplay(t *testing.T) {
	total, a, b := 23, 11, 12
	t.Logf("A %.1f%% (%d) B %.1f%% (%d)", float64(a)/float64(total)*100, a, float64(b)/float64(total)*100, b)
}

