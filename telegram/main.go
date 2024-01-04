package telegram

import (
	s "bot/scrapper"
	"fmt"

	log "github.com/sirupsen/logrus"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMessage(message, telegramApi string, telegramChatId int64) {
	bot, err := tgbotapi.NewBotAPI(telegramApi)
	if err != nil {
		log.Panic(err)
	}

	msg := tgbotapi.NewMessage(telegramChatId, message)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	log.Debug(msg)
	_, err = bot.Send(msg)
	if err != nil {
		log.Panic(err)
	}
}

func SendNewPosition(u, id, telegramApi string, telegramChatId int64, p s.Position) {
	message := fmt.Sprintf(`🚨 <b>Ouverture</b> de position
	
	👑 <b>Trader</b>: <a href="https://www.binance.com/en/futures-activity/leaderboard/user/um?encryptedUid=%s">%s</a>
	🚀 <b>Crypto</b>: %s
	🔴 <b>Type</b>: %s

	📈 <b>Entry</b>: %.4f
	✳️ <b>Levier</b>: %dx
	💰 <b>Montant</b>: %.2f
	`, id, u, p.Ticker, p.Direction, p.EntryPrice, p.Leverage, p.Amount)

	sendMessage(message, telegramApi, telegramChatId)
}

func SendClosedPosition(u, id, telegramApi string, telegramChatId int64, p s.Position) {
	message := fmt.Sprintf(`🚨 <b>Fermeture</b> de position
	
	👑 <b>Trader</b>: <a href="https://www.binance.com/en/futures-activity/leaderboard/user?encryptedUid=%s">%s</a>
	🚀 <b>Crypto</b>: %s
	🔴 <b>Type</b>: %s

	📈 <b>Entry / Close</b>: %.4f -> %.4f
	✳️ <b>Levier</b>: %dx
	💰 <b>Montant</b>: %.2f
	💸 <b>Profit</b>: %.2f$ (%.2f%%)
	`, id, u, p.Ticker, p.Direction, p.EntryPrice, p.MarkPrice, p.Leverage, p.Amount, p.Pnl, p.Roe*100)

	sendMessage(message, telegramApi, telegramChatId)
}

func SendAddedToPosition(u, id, telegramApi string, telegramChatId int64, p s.Position) {
	message := fmt.Sprintf(`🚨 <b>Ouverture</b> de position (DCA)
	
	👑 <b>Trader</b>: <a href="https://www.binance.com/en/futures-activity/leaderboard/user?encryptedUid=%s">%s</a>
	🚀 <b>Crypto</b>: %s
	🔴 <b>Type</b>: %s

	📈 <b>Entry</b>: %.4f
	✳️ <b>Levier</b>: %dx
	💰 <b>Montant</b>: %.2f -> %.2f
	`, id, u, p.Ticker, p.Direction, p.EntryPrice, p.Leverage, p.PrevAmount, p.Amount)

	sendMessage(message, telegramApi, telegramChatId)
}
