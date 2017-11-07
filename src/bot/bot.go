package main

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/yaml.v2"

	"flag"
	"io/ioutil"
	"log"
	"parseforum"
	"time"
)

var (
	configFile = flag.String("c", "config.yml", "config file") // читаем переданные параметры.
	confDeug   = flag.Bool("v", false, "debug log")
)

type ConfigStr struct { // структура файла конфига
	BotToken    string `yaml:"BotToken"`
	ChatId      int64  `yaml:"ChatId"`
	ForumUrl	string `yaml:"ForumUrl"`
	UrlLogin    string `yaml:"UrlLogin"`
	UrlFindNew  string `yaml:"UrlFindNew"`
	UrlMarkRead string `yaml:"UrlMarkRead"`
	UserName    string `yaml:"UserName"`
	Password    string `yaml:"Password"`
}

func main() {

	flag.Parse()       // парсим параметры
	Conf := ReadConf() // читаю конфиг из аргумента или дефолтного пути

	bot, err := tgbotapi.NewBotAPI(Conf.BotToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = *confDeug

	if *confDeug {
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	go func() {
		for {
			reply := parseforum.GetNew(Conf.ForumUrl, Conf.UserName, Conf.Password, *confDeug)
			if *confDeug {
				log.Printf("GetNew return: " + reply)
			}
			if reply != "" {
				reply = "Новое сообщение на форуме в теме: " + reply
				msg := tgbotapi.NewMessage(Conf.ChatId, reply)
				bot.Send(msg)
			}
			time.Sleep(60000 * time.Millisecond)
		}
	}()

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			switch update.Message.Command() {
			case "start":
				reply := "hi im bot"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
			case "new":
				reply := "Проверка сообщений. \n Новое сообщение на форуме в теме: "
				reply += parseforum.GetNew(Conf.ForumUrl, Conf.UserName, Conf.Password, *confDeug)
				if *confDeug {
					log.Printf("GetNew return: " + reply)
				}
				msg := tgbotapi.NewMessage(Conf.ChatId, reply)
				bot.Send(msg)
			}
		}
	}
}

func ReadConf() ConfigStr {
	if *confDeug {
		log.Print("Open config file")
	}
	data, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	var Conf ConfigStr
	if *confDeug {
		log.Print("Parse config file")
	}
	err = yaml.Unmarshal([]byte(data), &Conf)
	if err != nil {
		log.Fatal(err)
	}
	if *confDeug {
		log.Println(Conf)
	}
	return Conf
}
