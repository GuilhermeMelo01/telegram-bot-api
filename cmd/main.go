package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Print("No .env file found")
	}

}
func main() {
	botToken := os.Getenv("BOT_TOKEN")
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 5

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			userID := update.Message.From.ID
			msg := tgbotapi.NewMessage(userID, "Seu plano está espirando")

			msg.DisableWebPagePreview = true

			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal("tua prima")
			}
		}
	}
}

func getChatIdListFromBot(telegramBotToken string) int64 {

	apiUrl := fmt.Sprintf("%s%s%s", "https://api.telegram.org/bot", telegramBotToken, "/getUpdates")

	request, err := http.NewRequest("GET", apiUrl, nil)

	if err != nil {
		fmt.Println(err)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	responseBody, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
	}

	var data map[string]interface{}

	if err := json.Unmarshal(responseBody, &data); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	// Check if "result" key exists in the response
	if updates, ok := data["result"].([]interface{}); ok {
		// Iterate through updates
		for _, update := range updates {
			fmt.Println(update)
			if message, ok := update.(map[string]interface{})["message"].(map[string]interface{}); ok {
				if chat, ok := message["chat"].(map[string]interface{}); ok {
					if chatID, ok := chat["id"].(float64); ok {
						return int64(chatID)
					}
				}
			}
		}
	}

	return 0
}
