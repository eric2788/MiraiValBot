package file

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

// 文件實例

var ApplicationYaml *Configuration = &defaultConfig

//

func loadYaml(filename string, t interface{}) error {

	b, err := os.ReadFile(filename)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, t)
}

// LoadApplicationYaml 加載 application.yaml 並加載到 文件實例 上
func LoadApplicationYaml() {
	if err := loadYaml("application.yaml", ApplicationYaml); err != nil {
		fmt.Printf("加載 application.yaml 失敗: %v\n", err)
		os.Exit(1)
	}
}
