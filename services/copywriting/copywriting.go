package copywriting

import (
	"github.com/eric2788/common-utils/request"
)

const (
	fabingURL        = "https://raw.githubusercontent.com/SAGIRI-kawaii/sagiri-bot/Ariadne-v3/sagiri_bot/handler/handlers/ill/ill_templates.json"
	fadianURL        = "https://raw.githubusercontent.com/FloatTech/zbpdata/e8d06b150b2cf84d9c7dc2f8a9f573da2b2290fd/Fadian/post.json"
	cpDataURL        = "https://raw.githubusercontent.com/SAGIRI-kawaii/sagiri-bot/Ariadne-v3/statics/cp_data.json"
	tiangouURL       = "https://raw.githubusercontent.com/SAGIRI-kawaii/sagiri-bot/Ariadne-v4/modules/self_contained/pero_dog/pero_content.json"
	tiangou2URL      = "https://raw.githubusercontent.com/pcrbot/cappuccilo_plugins/master/generator/diary_data.json"
	ranranURL        = "https://raw.githubusercontent.com/RMYHY/RBot/main/HoshinoBot/hoshino/modules/asill/data.json"
	crazyThursdayURL = "https://raw.githubusercontent.com/MinatoAquaCrews/nonebot_plugin_crazy_thursday/beta/nonebot_plugin_crazy_thursday/post.json"
)

type AsoulWriting struct {
	Person string `json:"person"`
	Text   string `json:"text"`
}

//var logger = utils.GetModuleLogger("services.copywriting")

func GetFabingList() ([]string, string, error) {
	list, err := getJsonList(fabingURL, "data")
	return list, "{target}", err
}

func GetFadianList() ([]string, string, error) {
	list, err := getJsonList(fadianURL, "post")
	return list, "阿咪", err
}

func GetCPList() ([]string, string, string, error) {
	list, err := getJsonList(cpDataURL, "data")
	return list, "<攻>", "<受>", err
}

func GetTianGouList() ([]string, error) {
	return getJsonList(tiangouURL, "data")
}

func GetTiangou2List() ([]string, error) {
	return getList(tiangou2URL)
}

func GetRanranList() ([]AsoulWriting, error) {
	var as []AsoulWriting
	err := request.Get(ranranURL, &as)
	if err != nil {
		return nil, err
	}
	return as, nil
}

func GetCrazyThursdayList() ([]string, error) {
	return getJsonList(crazyThursdayURL, "post")
}

func getList(url string) ([]string, error) {
	var resp []interface{}
	err := request.Get(url, &resp)
	if err != nil {
		return nil, err
	}
	list := make([]string, len(resp))
	for i, v := range resp {
		list[i] = v.(string)
	}
	return list, nil
}

func getJsonList(url, key string) ([]string, error) {
	var resp map[string]interface{}
	err := request.Get(url, &resp)
	if err != nil {
		return nil, err
	}
	list := resp[key].([]interface{})
	results := make([]string, len(list))
	for i, v := range list {
		results[i] = v.(string)
	}
	return results, nil
}
