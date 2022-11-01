package valorant

import (
	"testing"

	"github.com/eric2788/MiraiValBot/redis"
	"github.com/stretchr/testify/assert"
)

func TestGetAgents(t *testing.T) {
	req := NewResourceRequest("/agents")
	req.SetLanguage(TC)
	req.AddQuery("isPlayableCharacter", "true")
	var agents []AgentData

	if err := req.DoRequest(&agents); err != nil {
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
	req := NewResourceRequest("/weapons")
	req.SetLanguage(TC)
	var weapons []WeaponData

	if err := req.DoRequest(&weapons); err != nil {
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
	req := NewResourceRequest("/bundles")
	req.SetLanguage(TC)
	var bundles []BundleData

	if err := req.DoRequest(&bundles); err != nil {
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
