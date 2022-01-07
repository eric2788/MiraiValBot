package file

import (
	"encoding/json"
	"github.com/Logiase/MiraiGo-Template/utils"
	"io/fs"
	"os"
	"sync"
)

var logger = utils.GetModuleLogger("file.storage")

type StorageData struct {
	Answers   map[string]bool   `json:"answers"`
	Responses map[string]string `json:"responses"`
	Bilibili  *BilibiliSettings `json:"bilibili"`
	Setting   *Setting          `json:"setting"`
	Listening *Listening        `json:"listening"`
}

type BilibiliSettings struct {
	HighLightedUsers []int64 `json:"highLightedUsers"`
}

type Listening struct {
	Bilibili []int64  `json:"bilibili"`
	Youtube  []string `json:"youtube"`
	Twitter  []string `json:"twitter"`
}

type Setting struct {
	VerboseDelete bool  `json:"verboseDelete"`
	Verbose       bool  `json:"verbose"`
	YearlyCheck   bool  `json:"yearlyCheck"`
	LastChecked   int64 `json:"lastChecked"`
}

const StoragePath = "data/valData.json"

var DataStorage *StorageData = &defaultStorageData
var locker sync.Mutex

var defaultStorageData = StorageData{
	Answers:   make(map[string]bool),
	Responses: make(map[string]string),
	Bilibili: &BilibiliSettings{
		HighLightedUsers: []int64{},
	},
	Setting: &Setting{
		VerboseDelete: false,
		Verbose:       false,
		YearlyCheck:   true,
		LastChecked:   0,
	},
	Listening: &Listening{
		Bilibili: []int64{},
		Youtube:  []string{},
		Twitter:  []string{},
	},
}

func LoadStorage() {
	err := os.MkdirAll("data", fs.ModePerm)

	if err != nil {
		logger.Fatalf("生成 data 文件夾時出現錯誤: %v", err)
		os.Exit(1)
	}

	generate(StoragePath, func() error {

		content, err := json.Marshal(defaultStorageData)

		if err != nil {
			return err
		}

		return os.WriteFile(StoragePath, content, 0775)
	})

	content, err := os.ReadFile(StoragePath)

	if err != nil {
		logger.Fatalf("讀取 %s 失敗: %v\n", StoragePath, err)
		os.Exit(1)
	}

	err = json.Unmarshal(content, DataStorage)

	if err != nil {
		logger.Fatalf("加載 %s 失敗: %v\n", StoragePath, err)
		os.Exit(1)
	}

}

// UpdateStorage use this function when mutating data for thread safe
func UpdateStorage(updateFunc func()) {
	locker.Lock()
	defer locker.Unlock()
	updateFunc()
	content, err := json.Marshal(DataStorage)
	if err != nil {
		logger.Warnf("讀取最新數據內容時出現錯誤: %v\n", err)
		return
	}
	err = os.WriteFile(StoragePath, content, 0755)
	if err != nil {
		logger.Warnf("更新數據內容時出現錯誤: %v\n", err)
		return
	}
	logger.Infof("數據內容已成功更新。")
}
