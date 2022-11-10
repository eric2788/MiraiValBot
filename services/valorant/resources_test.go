package valorant

import (
	"testing"

	"github.com/eric2788/MiraiValBot/internal/redis"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	redis.Init()
}

func TestGetAgents(t *testing.T) {
	agents, err := GetAgents(AllAgents, TC)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}

	assert.True(t, len(agents) >= 18, "agents should have at least 18")

	for _, agent := range agents {
		t.Log(agent.DisplayName)
	}
}

func TestGetWeapons(t *testing.T) {

	weapons, err := GetWeapons(AllWeapons, TC)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}

	assert.True(t, len(weapons) >= 18, "weapons should have at least 18")

	for _, weapon := range weapons {
		t.Log(weapon.DisplayName)
	}
}

func TestGetBundles(t *testing.T) {

	bundles, err := GetBundles(TC)
	if err != nil {
		if isAllowedStatus(err) {
			return
		}
		t.Fatal(err)
	}

	for _, bundle := range bundles {
		t.Log(bundle.DisplayName)
	}
}

func init() {
	redis.Init()
}
