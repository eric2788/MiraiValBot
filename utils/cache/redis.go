package cache

import (
	"fmt"

	"github.com/eric2788/MiraiValBot/internal/redis"
)

type RedisCache struct {
}

func (r *RedisCache) Init(path string) error {
	//redis.Init() // init already
	return nil
}

func (r *RedisCache) Save(path string, name string, data []byte) error {
	return redis.StoreBytes(fmt.Sprintf("miralbot_cache:%s:%s", path, name), data, redis.Permanent)
}

func (r *RedisCache) Get(path string, name string) (b []byte, err error) {
	b, exist, err := redis.GetBytes(fmt.Sprintf("miralbot_cache:%s:%s", path, name))
	if !exist {
		err = fmt.Errorf("cache not exist: %s:%s", path, name)
	}
	return
}

func (r *RedisCache) Remove(path string, name string) error {
	return redis.Delete(fmt.Sprintf("miralbot_cache:%s:%s", path, name))
}

func (r *RedisCache) List(path string) []CachedData {
	var results []CachedData
	keys, err := redis.ListKeys(fmt.Sprintf("miralbot_cache:%s", path))
	if err != nil {
		logger.Errorf("获取缓存列表时出现错误: %v", err)
		return results
	}
	for _, key := range keys {
		data := CachedData{
			Name: key,
			Path: fmt.Sprintf("miralbot_cache:%s", path),
		}
		results = append(results, data)
	}
	return results
}
