package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"task1crypto/internal/cache"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	update tgbotapi.UpdateConfig
}

func NewBot(token string) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram API: %v", err)
	}

	update := tgbotapi.NewUpdate(0)
	update.Timeout = 1
	return &Bot{api: api, update: update}, nil
}

func (b *Bot) Start(pricesCache *cache.Cache) {
	updates, err := b.api.GetUpdatesChan(b.update)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil && update.CallbackQuery != nil {
			if update.CallbackQuery.Data == "show_courses" {
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, formatPrices(pricesCache.GetData()))
				b.api.Send(msg)
			}
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Добро пожаловать!")
				msg.ReplyMarkup = b.showCoursesButton()
				b.api.Send(msg)
			}
		}
	}

}

func (b *Bot) showCoursesButton() tgbotapi.InlineKeyboardMarkup {
	var keyboard tgbotapi.InlineKeyboardMarkup

	row := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Показать курсы", "show_courses"),
	)

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, row)

	return keyboard
}

func formatPrices(allPrices map[string]map[string][]float64) string {
	result := "📈 Курсы криптовалют:\n\n"
	for cryptoName, valuesMap := range allPrices {
		for fiatName, prices := range valuesMap {
			result += fmt.Sprintf("%s к %s:\n", cryptoName, fiatName)
			result += fmt.Sprintf("🔹 Текущий: %d\n", int(prices[1]))
			result += fmt.Sprintf("🔹 Минимальный: %d\n", int(prices[0]))
			result += fmt.Sprintf("🔹 Максимальный: %d\n\n", int(prices[2]))
		}
	}
	return result
}
