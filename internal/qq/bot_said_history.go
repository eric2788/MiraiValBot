package qq

import (
	"github.com/eric2788/MiraiValBot/internal/redis"
)

var botSaid = &BotSaidHistory{
	key: "bot_self_msg",
}

type BotSaidHistory struct {
	key string
}

func (h *BotSaidHistory) Remove(id int64) {
	if err := redis.SetRemove(GroupKey(ValGroupInfo.Uin, h.key), id); err != nil {
		logger.Warnf("Redis 移除機器人聊天記錄時出現錯誤: %v", err)
	} else {
		logger.Infof("Redis 移除機器人聊天記錄成功。")
	}
}

func (h *BotSaidHistory) Add(id int32) {
	if err := redis.SetAdd(GroupKey(ValGroupInfo.Uin, h.key), id); err != nil {
		logger.Warnf("Redis 儲存機器人聊天記錄時出現錯誤: %v", err)
	} else {
		logger.Infof("Redis 儲存機器人聊天記錄成功。")
	}
}

func (h *BotSaidHistory) Contains(id int64) bool {
	exist, err := redis.SetContains(GroupKey(ValGroupInfo.Uin, h.key), id)
	if err != nil {
		logger.Errorf("嘗試從 redis 檢查 bot 訊息列表時出現錯誤: %v", err)
		return false
	}
	return exist
}
