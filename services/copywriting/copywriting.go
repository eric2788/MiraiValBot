package copywriting

import (
	"github.com/eric2788/common-utils/request"
)

const (
	fabingURL  = "https://raw.githubusercontent.com/SAGIRI-kawaii/sagiri-bot/Ariadne-v3/sagiri_bot/handler/handlers/ill/ill_templates.json"
	fadianURL  = "https://raw.githubusercontent.com/FloatTech/zbpdata/e8d06b150b2cf84d9c7dc2f8a9f573da2b2290fd/Fadian/post.json"
	cpDataURL  = "https://raw.githubusercontent.com/SAGIRI-kawaii/sagiri-bot/Ariadne-v3/statics/cp_data.json"
	tiangouURL = "https://raw.githubusercontent.com/SAGIRI-kawaii/sagiri-bot/Ariadne-v4/modules/self_contained/pero_dog/pero_content.json"
)

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
