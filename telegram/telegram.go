package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuri-ncs/novel-chapter-scraper/database"
	"log"
	"os"
)

func Start() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			commandHandler(update, bot)
		}
	}
}

func commandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "start":
		msg.Text = "Welcome to the Novel Scraper Bot!.\n We have a limited set of sites we can get the chapters from.\nType /help for a list of commands."
		msg.ReplyMarkup = startCmdKeyboard()
	case "help":
		msg.Text = "Type /start to re/start the bot."
	case "sites":
		sites := database.GetSitesList()
		msg.Text = "We have the following sites available:\n" + sites
	default:
		msg.Text = "I don't know that command\nType /help for a list of commands."
		msg.ReplyToMessageID = update.Message.MessageID
		msg.ReplyMarkup = mainCmdKeyboard()
	}
	bot.Send(msg)
}

func mainCmdKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var cmdKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/start"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/help"),
		),
	)
	return cmdKeyboard
}

func startCmdKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var cmdKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/sites"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/novels"),
		),
	)
	return cmdKeyboard
}
