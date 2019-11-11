package telegram

type UserTlg struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsBot     bool   `json:"is_bot"`
	UserName  string `json:"username"`
}

type ChatTlg struct {
	ID int `json:"id"`
}

type MessageEntityTlg struct {
	Type string `json:"type"`
}

type MessageTlg struct {
	MessageID int                `json:"message_id"`
	From      UserTlg            `json:"from"`
	Chat      ChatTlg            `json:"chat"`
	Entities  []MessageEntityTlg `json:"entities"`
	Text      string             `json:"text"`
}

type UpdateTlg struct {
	UpdateId      int        `json:"update_id"`
	Message       MessageTlg `json:"message"`
	EditedMessage MessageTlg `json:"edited_message"`
}

type ReceiveTlg struct {
	Ok          bool        `json:"ok"`
	Result      []UpdateTlg `json:"result"`
	Description string      `json:"description"`
}

type ReceiveReturnMSgTlg struct {
	Ok          bool       `json:"ok"`
	Result      MessageTlg `json:"result"`
	Description string     `json:"description"`
}
