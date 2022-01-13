package youtube

import "fmt"

func GetChannelLink(id string) string {
	return fmt.Sprintf("https://youtube.com/channel/%s", id)
}

func GetYTLink(info *LiveInfo) string {
	if info.Info != nil {
		return fmt.Sprintf("https://youtu.be/%s", info.Info.Id)
	} else {
		return fmt.Sprintf("https://youtube.com/channel/%s/live", info.ChannelId)
	}
}
