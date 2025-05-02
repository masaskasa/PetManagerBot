package telegram

type ReceivedUpdates struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"Message"`
}

type Message struct {
	ID   int    `json:"message_id"`
	From User   `json:"from"`
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
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

func newTextMessage(chatID int, text string) *TextMessage {
	return &TextMessage{
		ChatID: chatID,
		Text:   text,
	}
}
