package bilibili

import (
	"github.com/eric2788/MiraiValBot/internal/file"
)

var bSettings = &file.DataStorage.Bilibili

func AddHighlightUser(user int64) bool {

	if (*bSettings).HighLightedUsers.Contains(user) {
		return false
	}

	file.UpdateStorage(func() {
		(*bSettings).HighLightedUsers.Add(user)
	})

	return true
}

func RemoveHighlightUser(user int64) bool {

	if !(*bSettings).HighLightedUsers.Contains(user) {
		return false
	}

	file.UpdateStorage(func() {
		(*bSettings).HighLightedUsers.Delete(user)
	})

	return true
}

func IsHighlighter(user int64) bool {
	return (*bSettings).HighLightedUsers.Contains(user)
}
