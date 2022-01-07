package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/sites/twitter"
)

func HandleReTweet(bot *bot.Bot, data *twitter.TweetStreamData) error {
	return nil
}

func init() {
	twitter.RegisterDataHandler(twitter.ReTweet, HandleReTweet)
}
