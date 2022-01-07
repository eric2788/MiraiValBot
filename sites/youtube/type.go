package youtube

const (
	Live     = "live"
	Idle     = "idle"
	UpComing = "upcoming"
)

type LiveInfo struct {
	ChannelId   string `json:"channelId"`
	ChannelName string `json:"channelName"`
	Status      string `json:"status"`

	Info *struct {
		Cover       *string `json:"cover"`
		Title       string  `json:"title"`
		Id          string  `json:"id"`
		PublishTime int64   `json:"publishTime"`
		Description string  `json:"description"`
	} `json:"info"`
}
