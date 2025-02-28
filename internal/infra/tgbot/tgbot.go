package tgbot

import (
	"b8boost/backend/internal/entities"
	"fmt"
	"log"
	"time"

	tele "gopkg.in/telebot.v4"
)

type TgBot struct {
	bot      *tele.Bot
	userRepo entities.UserRepo
}

func NewTgBot(token string, userRepo entities.UserRepo) TgBot {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		panic(err)
	}

	return TgBot{
		bot:      b,
		userRepo: userRepo,
	}
}

func (b TgBot) SendMessage(userID int64, message string) {
	user, err := b.userRepo.GetByID(uint(userID))
	if err != nil {
		return
	}
	if user.TelegramID == nil || *user.TelegramID == 0 {
		return
	}
	recipient := &tele.User{ID: int64(*user.TelegramID)}

	b.bot.Send(recipient, message)
}

func (b TgBot) handleStart(c tele.Context) error {

	menuButton := &tele.MenuButton{
		Type:   "web_app",
		Text:   "Открыть приложение",
		WebApp: &tele.WebApp{URL: "https://app.b8st.ru"},
	}
	if err := b.bot.SetMenuButton(c.Sender(), menuButton); err != nil {
		log.Printf("Ошибка установки кнопки меню: %v", err)
	}
	return c.Reply(fmt.Sprintf("Привет! 👋 \nСпасибо, что выбрали нашего бота. Для завершения активации вам необходимо отправить ваш telegram_id администратору.\nВаш telegram_id: %d\nПожалуйста, скопируйте этот номер и отправьте его администратору для подтверждения доступа к функциям бота. После проверки вашего ID, администратор активирует ваш аккаунт, и вы сможете пользоваться всеми возможностями нашего сервиса.\nЕсли у вас возникнут вопросы, обращайтесь к администратору.\nЖелаем приятного использования!", c.Sender().ID))
}

func (b TgBot) Start() {
	b.bot.Handle("/start", b.handleStart)
	go b.bot.Start()
}
