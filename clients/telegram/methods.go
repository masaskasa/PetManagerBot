package telegram

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"strconv"
)

const (
	getUpdates          = "getUpdates"
	sendMessage         = "sendMessage"
	answerCallbackQuery = "answerCallbackQuery"
)

func (client *Client) GetUpdates(offset int, limit int) ([]Update, error) {
	query := url.Values{}
	query.Add("offset", strconv.Itoa(offset))
	query.Add("limit", strconv.Itoa(limit))

	data, err := client.getRequest(getUpdates, query)
	if err != nil {
		return nil, err
	}

	var updates ReceivedUpdates

	if err := json.Unmarshal(data, &updates); err != nil {
		slog.Error("GetUpdates: error of parse response data:", err.Error())
		return nil, err
	}

	return updates.Updates, nil
}

func (client *Client) SendMessage(chatID int, text string, replyMarkup InlineKeyboardMarkup) (Message, error) {

	textMessage := createTextMessage(chatID, text, replyMarkup)

	jsonMessage, err := json.Marshal(textMessage)
	if err != nil {
		slog.Error("SendMessage: error of serialize request Message to json:", err.Error())
		return Message{}, err
	}

	data, err := client.postRequest(sendMessage, jsonMessage)
	if err != nil {
		return Message{}, err
	}

	var message Message

	if err := json.Unmarshal(data, &message); err != nil {
		slog.Error("SendMessage: error of parse response data:", err.Error())
		return Message{}, err
	}

	return message, nil
}

func createTextMessage(chatID int, text string, replyMarkup InlineKeyboardMarkup) interface{} {

	if replyMarkup.InlineKeyboard != nil {
		return TextMessageReplyMarkup{ChatID: chatID, Text: text, ReplyMarkup: replyMarkup}
	} else {
		return TextMessage{ChatID: chatID, Text: text}
	}
}

func (client *Client) AnswerCallbackQuery(callbackQueryID string, text string, showAlert bool) (Message, error) {

	jsonMessage, err := json.Marshal(CallbackQueryForAnswer{
		ID:        callbackQueryID,
		Text:      text,
		ShowAlert: showAlert,
	})
	if err != nil {
		slog.Error("AnswerCallbackQuery: error of serialize request CallbackQueryForAnswer to json:", err.Error())
		return Message{}, err
	}

	data, err := client.postRequest(answerCallbackQuery, jsonMessage)
	if err != nil {
		return Message{}, err
	}

	var message Message

	if err := json.Unmarshal(data, &message); err != nil {
		slog.Error("AnswerCallbackQuery: error of parse response data:", err.Error())
		return Message{}, err
	}

	return message, nil
}
