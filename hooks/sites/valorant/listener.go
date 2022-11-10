package valorant

import (
	"fmt"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/internal/file"
	bc "github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/eric2788/MiraiValBot/services/valorant"
)

var (
	listening = &file.DataStorage.Listening
	topic     = func(ch string) string { return fmt.Sprintf("valorant:%s", ch) }
)

func StartListen(name, tag string) (bool, error) {

	ac, err := valorant.GetAccountDetails(name, tag)
	if err != nil {
		return false, err
	}

	id := fmt.Sprintf("%s//%s#%s", ac.PUuid, name, tag)

	file.UpdateStorage(func() {
		(*listening).Valorant.Add(id)
	})

	info, err := bot.GetModule(bc.Tag)

	if err != nil {
		return false, err
	}

	broadcaster := info.Instance.(*bc.Broadcaster)

	return broadcaster.Subscribe(topic(ac.PUuid), MessageHandler)
}

func StopListen(name, tag string) (bool, error) {

	nameTag := fmt.Sprintf("%s#%s", name, tag)

	idToDelete, uuidToUnSub := "", ""
	for line := range (*listening).Valorant.Iterator() {
		parts := strings.Split(line, "//")
		if len(parts) != 2 {
			logger.Warnf("Invalid line in listening: %s", line)
			continue
		}
		if parts[1] == nameTag {
			idToDelete = line
			uuidToUnSub = parts[0]
			break
		}
	}

	if idToDelete == "" || uuidToUnSub == "" {
		return false, nil
	}

	file.UpdateStorage(func() {
		(*listening).Valorant.Delete(idToDelete)
	})

	info, _ := bot.GetModule(bc.Tag)

	broadcaster := info.Instance.(*bc.Broadcaster)

	result := broadcaster.UnSubscribe(topic(uuidToUnSub))

	return result, nil
}

func GetListening() []string {
	var displayNames []string
	for line := range (*listening).Valorant.Iterator() {
		parts := strings.Split(line, "//")
		if len(parts) != 2 {
			logger.Warnf("Invalid line in listening: %s", line)
			continue
		}
		displayNames = append(displayNames, parts[1])
	}
	return displayNames
}
