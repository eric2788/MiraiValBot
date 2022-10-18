package valorant

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/eric2788/common-utils/request"
	"github.com/eric2788/common-utils/set"
)

type (
	AgentType   string
	AblitiySlot string
	WeaponType  string
	Language    string
)

var LangAvailable = set.FromStrArr([]string{
	string(EN), 
	string(TC), 
	string(SC), 
	string(JP),
})

const (
	ResourceBaseUrl = "https://valorant-api.com/v1/"

	Inititator AgentType = "1b47567f-8f7b-444b-aae3-b0c634622d10"
	Guard      AgentType = "5fc02f99-4091-4486-a531-98459a3e95e9"
	Duelist    AgentType = "dbe8757e-9e92-4ed4-b39f-9dfc589691d4"
	Controller AgentType = "4ee40330-ecdd-4f2f-98a8-eb1243428373"
	AllAgents  AgentType = "ALL"

	SlotQ   AblitiySlot = "Ability1"
	SlotE   AblitiySlot = "Ability2"
	SlotC   AblitiySlot = "Generade"
	SlotX   AblitiySlot = "Ultimate"
	Passive AblitiySlot = "Passive"

	Heavy      WeaponType = "EEquippableCategory::Heavy"
	Rifle      WeaponType = "EEquippableCategory::Rifle"
	Shotgun    WeaponType = "EEquippableCategory::Shotgun"
	Pistol     WeaponType = "EEquippableCategory::Sidearm"
	Sniper     WeaponType = "EEquippableCategory::Sniper"
	SMG        WeaponType = "EEquippableCategory::SMG"
	Melee      WeaponType = "EEquippableCategory::Melee"
	AllWeapons WeaponType = "ALL"

	EN Language = "en-US"
	SC Language = "zh-CN"
	TC Language = "zh-TW"
	JP Language = "ja-JP"
)

type ResourceSchema struct {
	path     string
	language Language
	query    map[string]string
}

func NewResourceRequest(path string) *ResourceSchema {
	return &ResourceSchema{
		path:     path,
		language: SC,
		query:    make(map[string]string),
	}
}

func (r *ResourceSchema) AddQuery(key, value string) {
	r.query[key] = value
}

func (r *ResourceSchema) SetLanguage(lang Language) {
	r.language = lang
}

func (r *ResourceSchema) DoRequest(arg interface{}) error {

	url, err := url.Parse(ResourceBaseUrl + r.path)
	if err != nil {
		return err
	}
	query := url.Query()
	for key, value := range r.query {
		query.Add(key, value)
	}
	query.Add("Language", string(r.language))
	url.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return err
	}
	res, err := doRequest(req)
	var resp = &ResourceResp{}
	if err != nil {
		if httpErr, ok := err.(*request.HttpError); ok {
			if err := request.Read(httpErr.Response, resp); err == nil {
				return errors.New(resp.Error)
			} else {
				logger.Warnf("cannot parse http error response to Resp: %v, use back http error as error.", err)
			}
		}
		return err
	}
	logger.Debugf("response status for %v: %v", url.String(), res.Status)
	err = request.Read(res, &resp)
	if err != nil {
		return errors.New(res.Status)
	} else if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return resp.ParseData(arg)
}

func GetAgents(agentType AgentType, lang Language) ([]AgentData, error) {

	req := NewResourceRequest("/agents")
	req.SetLanguage(lang)
	req.AddQuery("isPlayableCharacter", "true")

	var agents []AgentData
	if err := req.DoRequest(&agents); err != nil {
		return nil, err
	}

	if agentType == AllAgents {
		return agents, nil
	}

	filtered := make([]AgentData, 0)
	for _, agent := range agents {
		if agent.Role.Uuid == string(agentType) {
			filtered = append(filtered, agent)
		}
	}
	return filtered, nil
}

func GetWeapons(weaponType WeaponType, lang Language) ([]WeaponData, error) {
	req := NewResourceRequest("/weapons")
	req.SetLanguage(lang)
	var weapons []WeaponData
	if err := req.DoRequest(&weapons); err != nil {
		return nil, err
	}

	if weaponType == AllWeapons {
		return weapons, nil
	}

	filtered := make([]WeaponData, 0)
	for _, weapon := range weapons {
		if weapon.Category == weaponType {
			filtered = append(filtered, weapon)
		}
	}
	return filtered, nil
}
