package valorant

import (
	"fmt"
	"time"

	"github.com/eric2788/MiraiValBot/internal/redis"
)

func cacheMatchHistories(histories []MatchData) {
	for _, data := range histories {
		cacheMatchData(&data)
	}
}

func cacheMatchData(data *MatchData) {
	if err := redis.StoreTimely(matchKey(data.MetaData.MatchId), data, time.Hour*24*30); err != nil {
		logger.Errorf("储存对战数据 (%s) 到快取时出现错误: %v", data.MetaData.MatchId, err)
	} else {
		logger.Infof("对战数据 (%s) 储存快取成功。", data.MetaData.MatchId)
	}
}

func matchKey(matchId string) string {
	return fmt.Sprintf("valorant:match:%s", matchId)
}
