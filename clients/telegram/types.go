package telegram

type receivedUpdates struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type Update struct {
	ID      int      `json:"update_id"`
	Message *Message `json:"Message"`
}

type Message struct {
	ID   int    `json:"message_id"`
	From user   `json:"from"`
	Chat chat   `json:"chat"`
	Text string `json:"text"`
}

type user struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
}

type chat struct {
	ID int `json:"id"`
}

type textMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func newTextMessage(chatID int, text string) *textMessage {
	return &textMessage{
		ChatID: chatID,
		Text:   text,
	}
}
