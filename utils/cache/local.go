package cache

import (
	"fmt"
	"os"
)

const localDirPath = "cache/"

type LocalCache struct {
}

func (c *LocalCache) Init(path string) error {
	return os.MkdirAll(localDirPath+path, os.ModePerm)
}

func (c *LocalCache) Save(path, name string, data []byte) error {
	return os.WriteFile(c.path(path, name), data, os.ModePerm)
}

func (c *LocalCache) Get(path, name string) ([]byte, error) {
	return os.ReadFile(c.path(path, name))
}

func (c *LocalCache) Remove(path, name string) error {
	return os.Remove(c.path(path, name))
}

func (c *LocalCache) List(path string) []CachedData {
	
	results := make([]CachedData, 0)
	files, err := os.ReadDir(localDirPath+path)

	if err != nil {
		logger.Errorf("获取缓存列表时出现错误: %v")
		return results
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		data := CachedData{
			Name: file.Name(),
			Path: c.path(path, file.Name()),
		}

		results = append(results, data)
	}

	return results
}

func (c *LocalCache) path(path, name string) string {
	return fmt.Sprintf("%s%s/%s", localDirPath, path, name)
}
