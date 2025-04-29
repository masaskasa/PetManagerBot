package telegram

const (
	getUpdates = "getUpdates"
)

type receivedUpdates struct {
	Ok      bool     `json:"ok"`
	Updates []update `json:"result"`
}

type update struct {
	ID      int      `json:"update_id"`
	Message *message `json:"message"`
}

type message struct {
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
