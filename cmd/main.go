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
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.AllowedUpdates = []string{"true"}
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)
	// log.Printf("Authorized on account %s", bot.Self.UserName)

	go func() {
		for update := range updates {
			if update.CallbackQuery != nil {
				buttonClicked := update.CallbackQuery.Data
				userClickedID := update.CallbackQuery.Message.Chat.ID
				emoji18 := "🎬"

				fmt.Println("--- BUTTON CLICKED --- ", buttonClicked)
				fmt.Println("--- BUTTON USERCLICKEDID --- ", userClickedID)
				switch buttonClicked {
				case "button_vip":
					msg := tgbotapi.NewMessage(userClickedID, "Temos esses planos para quem quer ser VIP \n\nVIP Silver - 1 Mês de Acesso VIP + Grupo Bônus R$14,99 \n\nVIP Premium - Mês de Acesso VIP + Grupo Bônus R$39,99 \n\nVIP Diamond - Acesso Vitalicio VIP + Grupo Bônus + Chat"+emoji18+" R$60,00")
					buttonVIP := tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("PLANO VIP SILVER", "button_vip"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("PLANO VIP PREMIUM", "button_duvidas"),
						),
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("PLANO VIP DIAMOND", "button_suporte"),
						),
					)

					msg.ReplyMarkup = buttonVIP

					_, err := bot.Send(msg)
					if err != nil {
						fmt.Print("err")
					}
				default:
					fmt.Print("HA HA HA ")
				}
			}
		}
	}()

	for update := range updates {
		if update.Message != nil {
			userID := update.Message.From.ID
			emoji := "✅"
			msg := tgbotapi.NewMessage(userID, "Quer fazer parte do grupo EXTRAS VIP? \n\n"+
				emoji+"Canal VIP com conteúdos completos de mais de 500 Filmes (conteúdo diário e atualizado)\n\n"+
				emoji+"Canal organizado por nome de cada Filme. (fácil de encontra o que procura)\n\n"+
				emoji+"Informações sobre atores (Principais noticias do mundo dos famosos)\n\n"+
				emoji+"Solicitar conteúdo de filmes especificos (sugerir ao Adm algum Filme ou Serie suporte 24h)")
			btn := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("QUERO SER VIP", "button_vip"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("TENHO DÚVIDAS", "button_duvidas"),
				),
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("FALAR COM SUPORTE", "button_suporte"),
				),
			)

			msg.ReplyMarkup = btn

			_, err := bot.Send(msg)
			if err != nil {
				log.Fatal("panic")
			}

			if update.CallbackQuery != nil {
				if update.CallbackQuery.Data == "button_pressed" {
					link := "https://jsonformatter.curiousconcept.com"
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "here's the link: "+link)

					msg.ParseMode = tgbotapi.ModeMarkdown
					msg.DisableWebPagePreview = true
					_, err := bot.Send(msg)
					if err != nil {
						log.Fatal("mensagem não enviada")
					}
				}
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
