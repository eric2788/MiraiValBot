package bilibili

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/internal/file"
	bc "github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/eric2788/common-utils/request"
)

type RoomInfo struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`

	Data interface{} `json:"data"`
}

type RoomInfoData struct {
	RoomId    int64  `json:"room_id"`
	Uid       int64  `json:"uid"`
	ShortId   int32  `json:"short_id"`
	Title     string `json:"title"`
	UserCover string `json:"user_cover"`

	// for external use

	LiveStatus int       `json:"live_status"`
	AreaName   string    `json:"area_name"`
	LiveTime   time.Time `json:"live_time"`
	Online     int64     `json:"online"`
	KeyFrame   string    `json:"keyframe"`
}

func (d *RoomInfoData) Parse(data map[string]interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}

const Host = "https://api.live.bilibili.com/room/v1/Room/get_info"

var (
	roomInfoCache = make(map[int64]*RoomInfo)
	listening     = &file.DataStorage.Listening
	topic         = func(room int64) string { return fmt.Sprintf("blive:%d", room) }
)

func StartListen(room int64) (bool, error) {

	if info, err := GetRoomInfo(room); err != nil {
		return false, err
	} else if info.Code != 0 {
		return false, fmt.Errorf(info.Msg)
	} else {
		// 轉換短號為房間號
		data := &RoomInfoData{}
		if m, ok := info.Data.(map[string]interface{}); ok {
			if err := data.Parse(m); err != nil {
				return false, err
			}
			room = data.RoomId
		}
	}

	file.UpdateStorage(func() {
		(*listening).Bilibili.Add(room)
	})

	info, _ := bot.GetModule(bc.Tag)

	broadcaster := info.Instance.(*bc.Broadcaster)

	return broadcaster.Subscribe(topic(room), MessageHandler)
}

func StopListen(room int64) (bool, error) {

	if !(*listening).Bilibili.Contains(room) {
		return false, nil
	}

	file.UpdateStorage(func() {
		(*listening).Bilibili.Delete(room)
	})

	info, _ := bot.GetModule(bc.Tag)

	broadcaster := info.Instance.(*bc.Broadcaster)

	result := broadcaster.UnSubscribe(topic(room))

	return result, nil
}

func GetRoomInfo(room int64) (*RoomInfo, error) {
	if info, ok := roomInfoCache[room]; ok {
		return info, nil
	}

	var info = &RoomInfo{}
	if err := request.Get(fmt.Sprintf("%s?room_id=%d", Host, room), info); err != nil {
		return nil, err
	} else if info.Code == -401 {
		logger.Warnf("請求處於非法訪問，正在更換 User-Agent 並重新請求...")
		return GetRoomInfo(room)
	}

	roomInfoCache[room] = info
	return info, nil
}

func ClearRoomInfo(room int64) bool {
	if room != -1 {
		if _, ok := roomInfoCache[room]; !ok {
			return false
		}
		delete(roomInfoCache, room)
		return true
	} else {
		roomInfoCache = make(map[int64]*RoomInfo)
		return true
	}
}
