package bilibili

import (
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/utils/array"
)

var bSettings = file.DataStorage.Bilibili

func HighlightUserExist(user int64) bool {
	return array.IndexOfInt64(bSettings.HighLightedUsers, user) != -1
}

func AddHighlightUser(user int64) bool {
	if HighlightUserExist(user) {
		return false
	}

	file.UpdateStorage(func() {
		bSettings.HighLightedUsers = append(bSettings.HighLightedUsers, user)
	})

	return true
}

func RemoveHighlightUser(user int64) bool {
	if !HighlightUserExist(user) {
		return false
	}

	index := array.IndexOfInt64(bSettings.HighLightedUsers, user)

	file.UpdateStorage(func() {
		bSettings.HighLightedUsers = array.RemoveInt64(bSettings.HighLightedUsers, index)
	})

	return true
}
