package telegram

import "fmt"

type UserInfo struct {
	Id             int
	FirstName      string
	LastName       string
	UserName       string
	isBot          bool
	Authorizations Authorization
}

type Authorization struct {
	authorized bool
	isAdmin    bool
}

func NewUser(msg *MessageTlg) (*UserInfo, error) {
	user := msg.From
	if user.ID == 0 {
		return nil, fmt.Errorf("unknow user, id == 0")
	}
	return &UserInfo{
		Id:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		UserName:  user.UserName,
		isBot:     user.IsBot,
	}, nil
}
