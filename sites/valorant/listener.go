package valorant

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/file"
	bc "github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/eric2788/MiraiValBot/valorant"
)

var (
	listening = &file.DataStorage.Listening
	topic     = func(ch string) string { return fmt.Sprintf("valorant:%s", ch) }
)

func StartListen(name, tag string) (bool, error) {

	if _, err := valorant.GetAccountDetails(name, tag); err != nil {
		return false, err
	}

	id := fmt.Sprintf("%s#%s", name, tag)

	file.UpdateStorage(func() {
		(*listening).Valorant.Add(id)
	})

	info, err := bot.GetModule(bc.Tag)

	if err != nil {
		return false, err
	}

	broadcaster := info.Instance.(*bc.Broadcaster)

	return broadcaster.Subscribe(topic(id), MessageHandler)
}

func StopListen(name, tag string) (bool, error) {

	id := fmt.Sprintf("%s#%s", name, tag)

	if !(*listening).Valorant.Contains(id) {
		return false, nil
	}

	file.UpdateStorage(func() {
		(*listening).Valorant.Delete(id)
	})

	info, _ := bot.GetModule(bc.Tag)

	broadcaster := info.Instance.(*bc.Broadcaster)

	result := broadcaster.UnSubscribe(topic(id))

	return result, nil
}
