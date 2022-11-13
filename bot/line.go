package bot

import (
	"log"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/spf13/viper"
)

func LineConfig() *linebot.Client {
	secret := viper.GetString("line.channel_secret")
	accesToken := viper.GetString("line.channel_access_token")

	bot, err := linebot.New(secret, accesToken)
	if err != nil {
		log.Fatal("initial line error: ", err)
	}
	log.Println(bot.GetBotInfo())
	log.Println(bot)
	return bot
}
