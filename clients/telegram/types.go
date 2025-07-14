package telegram

type ReceivedUpdates struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type Update struct {
	ID            int            `json:"update_id"`
	CallbackQuery *CallbackQuery `json:"callback_query"`
	Message       *Message       `json:"message"`
}

type Message struct {
	ID   int    `json:"message_id"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type CallbackQuery struct {
	ID   string `json:"id"`
	From User   `json:"from"`
	Data string `json:"data"`
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

type Chat struct {
	ID int `json:"id"`
}

type TextMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type TextMessageReplyMarkup struct {
	ChatID      int                  `json:"chat_id"`
	Text        string               `json:"text"`
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func NewInlineKeyboardMarkup() InlineKeyboardMarkup {

	buttons := make([][]InlineKeyboardButton, 0, 10)
	for i := range buttons {
		buttons[i] = make([]InlineKeyboardButton, 0, 10)
	}

	return InlineKeyboardMarkup{InlineKeyboard: buttons}
}

func (markup *InlineKeyboardMarkup) AddButtonInlineKeyboardMarkup(button *InlineKeyboardButton) {
	markup.InlineKeyboard = append(markup.InlineKeyboard, []InlineKeyboardButton{*button})
}

type CallbackQueryForAnswer struct {
	ID        string `json:"callback_query_id"`
	Text      string `json:"text"`
	ShowAlert bool   `json:"show_alert"`
}
