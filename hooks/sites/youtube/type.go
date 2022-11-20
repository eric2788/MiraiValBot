package youtube

const (
	Live     = "live"
	Idle     = "idle"
	UpComing = "upcoming"
)

type LiveInfo struct {
	Duplicate   bool   `json:"duplicate"` // whether the video id is same as latest checked
	ChannelId   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	Status      string `json:"status"`
	Info        *struct {
		Cover       *string `json:"cover"`
		Title       string  `json:"title"`
		Id          string  `json:"id"`
		PublishTime string  `json:"publishTime"`
		Description string  `json:"description"`
	} `json:"info"`
}
