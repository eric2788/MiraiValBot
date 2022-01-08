package bilibili

import (
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/utils/array"
)

var bSettings = file.DataStorage.Bilibili

func AddHighlightUser(user int64) bool {

	if array.IndexOfInt64(bSettings.HighLightedUsers, user) != -1 {
		return false
	}

	file.UpdateStorage(func() {
		bSettings.HighLightedUsers = array.AddInt64(bSettings.HighLightedUsers, user)
	})

	return true
}

func RemoveHighlightUser(user int64) bool {

	index := array.IndexOfInt64(bSettings.HighLightedUsers, user)

	if index == -1 {
		return false
	}

	file.UpdateStorage(func() {
		bSettings.HighLightedUsers = array.RemoveInt64(bSettings.HighLightedUsers, index)
	})

	return true
}
