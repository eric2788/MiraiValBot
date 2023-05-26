package cache

import (
	"fmt"

	"github.com/eric2788/MiraiValBot/services/github"
)

// GitCache using github repository as caching
// Deprecated: github repostiory has size limit, and doing this will against TOS
type GitCache struct {
}

func (g *GitCache) Init(path string) error {
	github.Init()
	return github.VerifySuccess()
}

func (g *GitCache) Save(path, name string, data []byte) error {
	return github.UpdateFile(fmt.Sprintf("%s/%s", path, name), data)
}

func (g *GitCache) Get(path, name string) ([]byte, error) {
	return github.DownloadFile(fmt.Sprintf("%s/%s", path, name))
}

func (g *GitCache) Remove(path, name string) error {
	return github.RemoveFile(fmt.Sprintf("%s/%s", path, name))
}

func (g *GitCache) List(path string) []CachedData {

	results := make([]CachedData, 0)

	files, err := github.ListDir(path)

	if err != nil {
		logger.Errorf("获取缓存列表时出现错误: %v", err)
		return results
	}
	for _, file := range files {

		if file.GetType() == "dir" {
			continue
		}

		data := CachedData{
			Name: file.GetName(),
			Path: file.GetPath(),
		}

		results = append(results, data)

	}

	return results
}
