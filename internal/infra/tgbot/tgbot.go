package tgbot

import (
	"time"

	tele "gopkg.in/telebot.v4"
)

type tgBot struct {
	bot *tele.Bot
}

func NewTgBot(token string) tgBot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	return tgBot{
		bot: b,
	}
}

func (b tgBot) handleStart() {

}

func (b tgBot) Start() {
	go b.Start()
}
