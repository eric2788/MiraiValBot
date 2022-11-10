package bilibili

import "encoding/json"

const (
	DanmuMsg         = "DANMU_MSG"
	SendGift         = "SEND_GIFT"
	GuardBuy         = "GUARD_BUY"
	SuperChatMessage = "SUPER_CHAT_MESSAGE"
	Live             = "LIVE"
	InteractWord     = "INTERACT_WORD"
	EntryEffect      = "ENTRY_EFFECT"
)

type MapParser interface {
	Parse(m map[string]interface{}) error
}

type LiveData struct {
	Command  string `json:"command"`
	LiveInfo struct {
		UID    int64   `json:"uid"`
		Title  string  `json:"title"`
		Name   string  `json:"name"`
		Cover  *string `json:"cover"`
		RoomId int64   `json:"room_id"`
	} `json:"live_info"`
	Content map[string]interface{} `json:"content"`
}

type SuperChatMessageData struct {
	UID       int64  `json:"uid"`
	Price     int    `json:"price"`
	Message   string `json:"message"`
	StartTime int64  `json:"start_time"`

	BackgroundColorStart string `json:"background_color_start"`
	BackgroundImage      string `json:"background_image"`
	BackgroundColor      string `json:"background_color"`

	UserInfo struct {
		Face      string `json:"face"`
		NameColor string `json:"name_color"`
		UName     string `json:"uname"`
	} `json:"user_info"`
}

func (d *SuperChatMessageData) Parse(m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, d)
}
