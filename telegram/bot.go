package bot

import (
	"TransBot/command"
	"crypto/ecdsa"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

type connect struct {
	*ethclient.Client
}

func Run(token string, url string) {
	var chatID int64

	log.Println(token)
	log.Println(url)
	userMap := make(map[string]*ecdsa.PrivateKey)
	ethConnect, err := New(url)
	if err != nil {
		log.Fatal("dont connect")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))
	log.Println(fmt.Sprintf("ChatID: %v", chatID))

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	_, err = bot.Send(tgbotapi.NewMessage(chatID, "Hello, guys!"))
	if err != nil {
		return
	}

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		reply := ""
		if update.Message == nil {
			reply = command.Help()

			continue
		}

		log.Println(fmt.Sprintf("[%s] %s", update.Message.From.UserName, update.Message.Text))

		switch update.Message.Command() {
		case "balance":
			reply = command.Ethbalance(userMap, update, ethConnect)
		case "transfer":
			reply = command.Transfer(userMap, update, ethConnect)
		default:
			reply = command.Help()
		}
		update.Message.CommandArguments()

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)

		_, err := bot.Send(msg)
		if err != nil {
			return
		}
	}
}

func New(url string) (*connect, error) {
	conn, err := ethclient.Dial(url)

	return &connect{conn}, err
}
