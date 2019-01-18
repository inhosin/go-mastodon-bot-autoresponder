package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mastodon "github.com/mattn/go-mastodon"
	"github.com/spf13/viper"
)

const (
	LogPrefix = "mastodon-bot-autoresponder-> "
)

type Bot struct {
	mstd   *mastodon.Client
	logger *log.Logger
}

func NewBot() *Bot {
	return &Bot{
		mstd: mastodon.NewClient(&mastodon.Config{
			Server:       viper.GetString("mastodon.server"),
			ClientID:     viper.GetString("mastodon.clientId"),
			ClientSecret: viper.GetString("mastodon.clientSecret"),
			AccessToken:  viper.GetString("mastodon.accessToken"),
		}),
		logger: log.New(
			os.Stderr,
			LogPrefix,
			log.Ldate|log.Ltime,
		),
	}
}

func (bot *Bot) log(v ...interface{}) {
	bot.logger.Println(v...)
}

func (bot *Bot) logError(err error) {
	bot.log("[ERROR]", err)
}

func (bot *Bot) sendToot(toot string) error {
	if _, err := bot.mstd.PostStatus(context.Background(), &mastodon.Toot{
		Status:     toot,
		Visibility: viper.GetString("response.visibility"),
	}); err != nil {
		bot.logError(err)
		return err
	}
	return nil
}

func (bot *Bot) sendWelcomeToot(toot string) error {
	if _, err := bot.mstd.PostStatus(context.Background(), &mastodon.Toot{
		Status:     toot,
		Visibility: "public",
	}); err != nil {
		bot.logError(err)
		return err
	}
	return nil
}

func (bot *Bot) Run() error {
	// start listen user event
	evChan, err := bot.mstd.StreamingUser(context.Background())
	if err != nil {
		return err
	}

	bot.log("start streaming")

	for env := range evChan {
		switch envType := env.(type) {
		case *mastodon.NotificationEvent:
			time.AfterFunc(time.Duration(viper.GetInt32("response.delay"))*time.Second, func() {
				notification := envType.Notification
				botName := fmt.Sprintf("@%s", viper.GetString("bot.name"))
				authorAccName := fmt.Sprintf("@%s", notification.Account.Acct)
				var toMessage string
				// switch notification type
				switch notification.Type {
				// follow notification type
				case "follow":
					bot.log("Welcomed new follower: @" + notification.Account.Acct)
					toMessage = fmt.Sprintf("%s @%s", viper.GetString("bot.welcome_message"), notification.Account.Acct)
					if err := bot.sendWelcomeToot(toMessage); err != nil {
						bot.logError(err)
					}
				// mention notification type
				case "mention":
					if notification.Status.InReplyToID != nil {
						//bot.log(ExtractMessage(notification.Status.Content))
						return
					}
					fromMessage, err := ExtractMessage(notification.Status.Content)
					if err != nil {
						bot.logError(err)
						return
					}
					bot.log(notification.Account.Acct+": ", fromMessage)
					fromMessage = strings.TrimPrefix(fromMessage, botName)
					toMessage = fmt.Sprintf("%s by %s", fromMessage, authorAccName)
					if err := bot.sendToot(toMessage); err != nil {
						bot.logError(err)
					}
				// others notification type
				default:
					bot.log("Recieved: " + notification.Type)
				}
			})
		case *mastodon.ErrorEvent:
			bot.logError(envType)
		}
	}

	bot.log("stop streaming")

	return nil
}
