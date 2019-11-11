package telegram

import (
	"fmt"
	"strings"
)

type CommandFunc func(cmd CommandInfo)

type Command struct {
	Enable             bool
	CmdName            string
	Callback           CommandFunc
	NeedAuthorization  bool
	AuthorizedUserID   []int
	UnauthorizedUserID []int
}

func NewCommand(cmdName string, callback CommandFunc) *Command {
	return &Command{
		CmdName:  cmdName,
		Enable:   true,
		Callback: callback,
	}
}

func (c *Command) AddAuthorizedId(id ...int) {

}

func (c *Command) DelAuthorizedId(id ...int) {

}

func (c *Command) AddUnauthorizedId(id ...int) {

}

func (c *Command) DelUnauthorizedId(id ...int) {

}

type CommandInfo struct {
	CmdName      string
	Args         []string
	Msg          *MessageTlg
	MessageID    int
	ChatID       int
	Restricted   bool
	isAuthorized bool
	isAdmin      bool
}

// GetInfo parse command and arguments, check authorization
func (c *Command) GetInfo(msg *MessageTlg) (CommandInfo, error) {
	var cmdInfo CommandInfo
	if c.NeedAuthorization {
		userId := msg.From.ID
		cmdInfo.isAuthorized = false // reinitialization to secure value
		for _, id := range c.AuthorizedUserID {
			if id == userId {
				cmdInfo.isAuthorized = true
				break
			}
		}
	}
	if len(msg.Text) < 1 {
		return CommandInfo{}, fmt.Errorf("GetInfo: empty msg.Text")
	}
	cmdInfo.ChatID = msg.Chat.ID
	cmdInfo.MessageID = msg.MessageID
	cmdInfo.Args = strings.Fields(msg.Text)[1:]
	cmdInfo.Msg = msg
	cmdInfo.CmdName = c.CmdName
	return cmdInfo, nil
}

func GetCommandName(msg *MessageTlg) string {
	if len(msg.Entities) > 0 && len(msg.Text) > 0 {
		if msg.Entities[0].Type == "bot_command" && msg.Text[0] == '/' {
			return strings.Fields(msg.Text)[0][1:]
		}
	}
	return ""
}
