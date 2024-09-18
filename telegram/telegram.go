package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/yuri-ncs/novel-chapter-scraper/database"
	"github.com/yuri-ncs/novel-chapter-scraper/models"
	"html/template"
	"log"
	"os"
	"strings"
	"sync"
)

var userStates = make(map[int64]string)
var mu sync.Mutex
var Bot *tgbotapi.BotAPI

func Start() {
	var err error
	Bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))

	if err != nil {
		panic(err)
	}
	Bot.Debug = true

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			commandHandler(update, Bot)
		}
	}
}

func UpdateUserState(userID int64, state string) {
	mu.Lock()
	defer mu.Unlock()
	userStates[userID] = state
}

func GetUserState(userID int64) string {
	mu.Lock()
	defer mu.Unlock()
	if state, ok := userStates[userID]; ok {
		return state
	}
	return "main" // Default to the main menu if no state is stored
}

func commandHandler(update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	switch update.Message.Command() {
	case "start":
		msg.Text = "Welcome to the Novel Scraper Bot!.\n We have a limited set of sites we can get the chapters from.\nType /help for a list of commands."
		UpdateUserState(update.Message.Chat.ID, "main_menu")
		msg.ReplyMarkup = startCmdKeyboard()
	case "user":
		UpdateUserState(update.Message.Chat.ID, "user_menu")
		msg.Text = "Click /signup to sign up for notifications.\nCLick /list to show your tracked novels\nClick /return to return to the main menu."
		msg.ReplyMarkup = userCmdKeyboard()
	case "help":
		UpdateUserState(update.Message.Chat.ID, "help_menu")
		msg.Text = "The following commands are available:\n/sites - List the available sites\n/novels - List the available novels\n/return - Return to the main menu"
	case "sites":
		sites := database.GetSitesList()
		msg.Text = "We have the following sites available:\n" + sites
	case "novels":
		msg.Text = "Choose an option."
		UpdateUserState(update.Message.Chat.ID, "novels_menu")
		msg.ReplyMarkup = novelsCmdKeyboard()
	case "add":
		if GetUserState(update.Message.Chat.ID) != "novels_menu" {
			break
		}
		msg.Text = "Use the example below to add a novel to the bot.\n"
		msg.Text += "/add_novel [novel name] [url]\n"
		msg.Text += "Example: /add_novel [The Novel] [https://www.example.com]"
		msg.Text += "Obs: Use the exact name of the novel and the full url of a supported site."
	case "add_novel":

		//split the response to get the name and the url
		response := update.Message.CommandArguments()
		fields := strings.Fields(response)

		if len(fields) < 2 {
			fmt.Println(fields)
			fmt.Println(len(fields))
			msg.Text = "Invalid command. Use the example below to add a novel to the bot.\n"
			msg.Text += "/add_novel [novel name] [url]\n"
			msg.Text += "Example: /add [The Novel] [https://www.example.com]"
			break
		}

		url := fields[len(fields)-1]

		fields = fields[:len(fields)-1] //removed the url from the fields

		var merged string

		for i := 0; i < len(fields); i++ {
			merged += fields[i]
			if i < len(fields)-1 {
				merged += " "
			}
		}
		//remove the brackets
		merged = strings.Replace(merged, "[", "", -1)
		merged = strings.Replace(merged, "]", "", -1)

		url = strings.Replace(url, "[", "", -1)
		url = strings.Replace(url, "]", "", -1)

		fmt.Println(merged)
		fmt.Println(url)

		//verify if the site in the field 1 is supported
		supported, siteID := database.VerifySupportedSite(url)

		if !supported {
			msg.Text = "Oops! The site is not supported yet."
			break
		}

		novel := models.Novel{
			Name:   merged,
			URL:    template.URL(url),
			SiteID: siteID,
		}

		err := database.CreateNovel(&novel)
		if err != nil {
			msg.Text = "Error creating the novel."
			break
		}

		msg.Text = "Novel added successfully!"

	case "return":
		msg.Text = "Type /start to re/start the bot."
		msg.ReplyMarkup = mainCmdKeyboard()
	case "track":
		if GetUserState(update.Message.Chat.ID) != "novels_menu" {
			break
		}
		response := update.Message.CommandArguments()
		fields := strings.Fields(response)

		var novelName string
		if len(fields) >= 2 {
			for i := 0; i < len(fields); i++ {
				novelName += fields[i]
				if i < len(fields)-1 {
					novelName += " "
				}
			}
		}

		novelName = strings.Replace(novelName, "[", "", -1)
		novelName = strings.Replace(novelName, "]", "", -1)

		fmt.Println(novelName)

		novel, err := database.GetNovelByName(novelName)

		if err != nil {
			msg.Text = "Novel not found.\n"
			msg.Text += "Verify for any typos and try again."
			break
		}

		user, err := database.GetUserByChatID(update.Message.Chat.ID)

		if err != nil {
			msg.Text = "Sign-up First!"
			break
		}

		tracking := database.TrackNovel(user.ID, novel.ID)

		if tracking {
			msg.Text = "Novel tracked successfully!"

		} else {
			msg.Text = "Novel already being tracked."
		}

	case "list":

		switch GetUserState(update.Message.Chat.ID) {
		case "novels_menu":
			novels := database.GetActiveNovels()
			msg.Text = "We have the following novels available:\n"
			for _, novel := range novels {
				// Append each novel title as a code block
				msg.Text += "`" + novel.Name + "`\n"

			}
			msg.ParseMode = "MarkdownV2"
			bot.Send(msg)
			msg = tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = "Use the example below to track a novel.\n\n"
			msg.Text += "/track [novel name]\n\n"

		case "user_menu":
			user, _ := database.GetUserByChatID(update.Message.Chat.ID)
			novels, _ := database.GetTrackedNovels(user.ID)
			msg.Text = "You are tracking the following novels:\n\n"
			for _, novel := range novels {
				msg.Text += novel + "\n"
			}
		default:
			msg.Text = "Type /novels to see the available novels."

		}

		if userStates[update.Message.Chat.ID] != "novels_menu" {
			break
		}

	case "signup":
		msg.Text = "You have been signed up for notifications.\nDont forget to select the novels you want to be notified about."
		user := models.User{
			ChatID: update.Message.Chat.ID,
		}
		database.CreateUser(&user)

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
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/user"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/return"),
		),
	)
	return cmdKeyboard
}

func userCmdKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var cmdKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/signup"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/list"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/return"),
		),
	)
	return cmdKeyboard
}

func novelsCmdKeyboard() tgbotapi.ReplyKeyboardMarkup {
	var cmdKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/list"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/add"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("/return"),
		),
	)
	return cmdKeyboard
}

func SendNotification(chapter models.Chapter, novelName string) {

	users, err := database.GetUsersByNovelID(chapter.NovelID)

	if err != nil {
		panic(err)
	}

	for _, user := range users {
		msg := tgbotapi.NewMessage(user.ChatID, "")
		msg.Text = "New chapter for " + novelName + "!\n"
		msg.Text += chapter.Title + "\n"
		msg.Text += "Link: " + chapter.Href

		Bot.Send(msg)
	}
}
