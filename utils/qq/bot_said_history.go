package qq

import (
	"github.com/eric2788/MiraiValBot/redis"
	"github.com/eric2788/MiraiValBot/utils/set"
)

var botSaid = &BotSaidHistory{
	said:   set.NewInt64(),
	edited: false,
	key:    "bot_self_msg",
}

type BotSaidHistory struct {
	said   *set.Int64Set
	edited bool
	key    string
}

func (h *BotSaidHistory) Add(id int32) {
	h.said.Add(int64(id))
	h.edited = true
}

func (h *BotSaidHistory) Contains(id int64) bool {
	return h.said.Contains(id)
}

func (h *BotSaidHistory) TakeFromRedis() error {
	botSaidArr := make([]int64, 0)

	if exist, err := redis.Get(GroupKey(ValGroupInfo.Uin, h.key), &botSaidArr); err != nil {
		return err
	} else if exist {
		h.said = set.FromInt64Arr(botSaidArr)
		logger.Infof("已成功獲取機器人的訊息記錄共 %d 條。", h.said.Size())
	}

	return nil
}

func (h *BotSaidHistory) SaveToRedis() {
	if !h.edited {
		return
	}

	if err := redis.Store(GroupKey(ValGroupInfo.Uin, h.key), h.said.ToArr()); err != nil {
		logger.Warnf("Redis 儲存機器人聊天記錄時出現錯誤: %v", err)
	} else {
		logger.Infof("Redis 儲存機器人聊天記錄成功。")
	}

	h.edited = false
}
