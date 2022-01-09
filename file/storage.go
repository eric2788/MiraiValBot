package file

import (
	"encoding/json"
	"github.com/Logiase/MiraiGo-Template/utils"
	"io/fs"
	"os"
	"sync"
)

var logger = utils.GetModuleLogger("file.storage")

const StoragePath = "data/valData.json"

var DataStorage *StorageData
var locker sync.Mutex

var defaultStorageData = storageData{
	Answers:   make(map[string]bool),
	Responses: make(map[string]string),
	Bilibili: &bilibiliSettings{
		HighLightedUsers: []int64{},
	},
	Setting: &setting{
		VerboseDelete: false,
		Verbose:       false,
		YearlyCheck:   true,
		LastChecked:   0,
	},
	Listening: &listening{
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

	err = json.Unmarshal(content, &defaultStorageData)

	if err != nil {
		logger.Fatalf("加載 %s 失敗: %v\n", StoragePath, err)
		os.Exit(1)
	}

	DataStorage = makeWrapper(&defaultStorageData)
}

var edited = false

// UpdateStorage use this function when mutating data for thread safe
func UpdateStorage(updateFunc func()) {
	locker.Lock()
	defer locker.Unlock()
	updateFunc()
	edited = true
}

// SaveStorage should use timer
func SaveStorage() {
	if !edited {
		return
	}
	locker.Lock()
	defer locker.Lock()
	content, err := json.Marshal(DataStorage.toRealStorageData())
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
	edited = false
}
