package myapp

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/acme/autocert"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	BotToken       = `1041050749:AAHI0gk4ML3WDgtJKt9chEBlERqYg2j5tYI`
	WebHookDNSName = `bot.avdeenko.com`
	WebhookURL     = `https://` + WebHookDNSName + `:443/`
)

//func main() {
	//
	var (
		err error
		bot = &tgbotapi.BotAPI{}
	)

	if bot, err = tgbotapi.NewBotAPI(BotToken); err != nil {
		log.Fatalln(err)
	}
	//bot.Debug = true
	log.Printf("Authorized with name:%v, with ID:%v\n", bot.Self.UserName, bot.Self.ID)

	// Удаляем Webchok если есть и создаем новый
	if whI, err := bot.GetWebhookInfo(); err == nil {
		if whI.IsSet() {
			bot.RemoveWebhook()
			log.Println("Webhook is removed")
		}
		if _, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL + bot.Token)); err != nil {
			log.Fatalln(err)
		}
		log.Println("Webhooks is set")
	} else {
		log.Fatalln(err)
	}
	// chanel с сообщениями
	updates := bot.ListenForWebhook("/" + bot.Token)
	// Запускаем Webhookc LetsEncrypt
	// Надо подусать об обновлении сертификата...
	go func() {
		if err := http.Serve(autocert.NewListener(WebHookDNSName), nil); err != nil {
			log.Fatal(err)
		}
		log.Println("SSL server was starting... ")
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = fmt.Sprintf("Hello! %v", update.Message.From)
			case "help":
				msg.Text = "type /sayhi or /status."
			case "sayhi":
				msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			default:
				msg.Text = "I don't know that command"
			}
			bot.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
//}
