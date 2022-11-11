package cache

import "sync"

func Migrate(from, to *Service, removal ...bool) *sync.WaitGroup {

	wg := &sync.WaitGroup{}

	remove := false
	if len(removal) > 0 {
		remove = removal[0]
	}
	for _, file := range from.List() {
		wg.Add(1)
		go func(file CachedData, wg *sync.WaitGroup) {
			defer wg.Done()
			b, err := from.Get(file.Name)
			if err != nil {
				logger.Errorf("读取缓存 %s 时出现错误: %v", file.Path, err)
			} else {
				if err := to.Set(file.Name, b); err == nil && remove {
					defer func() {
						if err = from.Remove(file.Name); err != nil {
							logger.Warnf("删除档案 %s 失败: %v", file.Path, err)
						}
					}()
				} else if err != nil {
					logger.Warnf("无法复制档案 %s: %v", file.Path, err)
				}
			}
		}(file, wg)
	}

	return wg
}
