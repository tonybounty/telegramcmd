package telegram

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/tonybounty/longpoll"
)

const (
	urlApiTelegram = "https://api.telegram.org/"
)

type Telegram struct {
	// URL to the Telegram BOT API, eg https://api.telegram.org/botTOKEN/
	BotURL string

	// offset parameter for getUpdates
	LastOffset string

	Timeout int

	// PollService   longpoll.LongPoll
	EndPoll       chan struct{}
	LoopStarted   bool
	cancelRequest context.CancelFunc

	Commands      []Command
	OtherTextFunc func(msg *MessageTlg)
}

func New(token string) *Telegram {
	return &Telegram{
		BotURL:        urlApiTelegram + "bot" + token + "/",
		OtherTextFunc: nil,
		LoopStarted:   false,
		Timeout:       120,
	}
}

func (t *Telegram) decodeJson(jsonData []byte, structData interface{}) error {
	err := json.Unmarshal(jsonData, structData)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

// checkReceiveErrTlg check if there is no error from API
func (t *Telegram) checkReceiveErrTlg(recv *ReceiveTlg) error {
	if recv.Ok != true {
		return fmt.Errorf("telegram api: %s", recv.Description)
	}
	return nil
}

// popMessage pop & return Message struct and update id from ReceiveTlg struct
func (t *Telegram) popMessage(recv *ReceiveTlg) (updateid int, message MessageTlg, valid bool) {
	if len(recv.Result) == 0 {
		return -1, MessageTlg{}, false
	}
	msgid := recv.Result[0].Message.MessageID
	msgEdid := recv.Result[0].EditedMessage.MessageID
	if msgEdid == 0 && msgid == 0 {
		return -1, MessageTlg{}, false
	}

	var msg MessageTlg
	if msgid > 0 {
		msg = recv.Result[0].Message
	} else {
		msg = recv.Result[0].EditedMessage
	}

	id := recv.Result[0].UpdateId
	recv.Result = recv.Result[1:]

	return id, msg, true
}

func (t *Telegram) poll(ctx context.Context, url string,
	outcomingMsg chan *MessageTlg, done chan struct{}) {
	r, err := longpoll.RunWithParm(ctx, url, map[string]string{
		"offset":  t.LastOffset,
		"timeout": strconv.Itoa(t.Timeout),
	})
	if err != nil {
		log.Println("Error polling: ", err)
		return
	}

	recv := &ReceiveTlg{}
	err = t.decodeJson(r, recv)
	if err != nil {
		fmt.Println("decode error: ", err)
		return
	}
	if err := t.checkReceiveErrTlg(recv); err != nil {
		log.Println("Error api: ", err)
		return
	}
	for {

		updateid, message, valid := t.popMessage(recv)
		if !valid {
			break
		}
		// sending message to startPollLoop
		outcomingMsg <- &message
		// updating offset
		t.LastOffset = strconv.Itoa(updateid + 1)
	}
	done <- struct{}{}

}

// Start or restart long polling loop for waiting message from telegram server
// Emit signal on EndPoll chan when server loop is stopped
func (t *Telegram) Start() (EndPoll <-chan struct{}) {
	t.EndPoll = make(chan struct{})
	t.startLoopHandler()
	t.LoopStarted = true
	return t.EndPoll
}

// Stop current polling, stop also the current request.
func (t *Telegram) Stop() {
	t.cancelRequest()
	t.LoopStarted = false
}

func (t *Telegram) startLoopHandler() {
	url := t.BotURL + "getUpdates"

	// context is used to stop http request in long polling loop
	ctx, cancelPoll := context.WithCancel(context.Background())
	t.cancelRequest = cancelPoll
	incomingMsg := make(chan *MessageTlg)
	done := make(chan struct{})

	go func() {
		go t.poll(ctx, url, incomingMsg, done)
		for {
			select {
			case <-ctx.Done():
				log.Println("context.Done() cancel called")
				close(t.EndPoll)
				return
			case msg := <-incomingMsg:
				go t.dispatchMessage(*msg)
			case <-done:
				go t.poll(ctx, url, incomingMsg, done)
			}
		}
	}()
}

func (t *Telegram) BindCommand(cmd *Command) {
	t.Commands = append(t.Commands, *cmd)
}

func (t *Telegram) BindOtherText(callback func(msg *MessageTlg)) {
	t.OtherTextFunc = callback
}

func (t *Telegram) handleCommand(msg *MessageTlg) {
	cmdName := GetCommandName(msg)
	var command *Command

	// search for command
	for _, cmd := range t.Commands {
		if cmd.CmdName == cmdName {
			command = &cmd
			break
		}
	}
	if command == nil {
		log.Printf("Command '%s' not found\n", cmdName)
		return
	}
	cmdInfo, err := command.GetInfo(msg)
	if err != nil {
		log.Println("handleCommand error: ", err)
		return
	}
	command.Callback(cmdInfo)
}

func (t *Telegram) dispatchMessage(msg MessageTlg) {
	if len(msg.Entities) > 0 {
		if msg.Entities[0].Type == "bot_command" {
			t.handleCommand(&msg)
			return
		}
	}

	// if not a command
	if t.OtherTextFunc != nil {
		t.OtherTextFunc(&msg)
	}
}

func (t *Telegram) addAuthorizedId(id int) {

}

func (t *Telegram) delAuthorizedId(id int) {

}

func (t *Telegram) addAdminId(id int) {

}

func (t *Telegram) delAdminId(id int) {

}

func (t *Telegram) getMe() {

}

func (t *Telegram) SendChatMessage(chatID int, text string) error {
	return t.SendMessageParam(map[string]string{
		"chat_id": strconv.Itoa(chatID),
		"text":    text,
	})
}

func (t *Telegram) ReplyChatMessage(chatID int, msgID int, text string) error {
	return t.SendMessageParam(map[string]string{
		"chat_id":             strconv.Itoa(chatID),
		"reply_to_message_id": strconv.Itoa(msgID),
		"text":                text,
	})
}

func (t *Telegram) SendMessageParam(param map[string]string) error {
	botUrl := t.BotURL + "sendMessage?"
	for paramName, value := range param {
		botUrl += "&" + paramName + "=" + url.QueryEscape(value)
	}
	log.Println("sending message with botUrl:", botUrl)
	res, err := http.Get(botUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	d := ReceiveReturnMSgTlg{}
	err = t.decodeJson(b, &d)
	if err != nil {
		return fmt.Errorf("json decode: %v", err)
	}

	if !d.Ok || res.StatusCode != http.StatusOK {
		return fmt.Errorf("api error: %s", d.Description)
	}

	return nil
}
