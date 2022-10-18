package valorant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAgents(t *testing.T) {
	agents, err := GetAgents(AllAgents, EN)
	if err != nil {
		if isAllowedStatus(err){
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
	weapons, err := GetWeapons(AllWeapons, EN)
	if err != nil {
		if isAllowedStatus(err){
			return
		}
		t.Fatal(err)
	}

	assert.True(t, len(weapons) >= 18, "weapons should have at least 18")

	for _, weapon := range weapons {
		t.Log(weapon.DisplayName)
	}
}