package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommandName(t *testing.T) {

	fakeMsgOther := []byte(`{"message_id":830,"from":{"id":789456,"is_bot":false,"first_name":"мαяgαяιηє","username":"Testname","language_code":"en"},"chat":{"id":123456,"first_name":"мαяgαяιηє","username":"Testname","type":"private"},"date":1572619731,"text":"boom"}`)
	fakeMsgCmd := []byte(`{"message_id":1053,"from":{"id":669450106,"is_bot":false,"first_name":"\u0443\u03c3\u03c5\u044f\u03b9","username":"wolnosc06","language_code":"en"},"chat":{"id":669450106,"first_name":"\u0443\u03c3\u03c5\u044f\u03b9 ","username":"wolnosc06","type":"private"},"date":1572979909,"text":"/mycmd arg1 arg2 arg3","entities":[{"offset":0,"length":7,"type":"bot_command"}]}`)
	tlg := Telegram{}

	var err error
	var msg MessageTlg

	err = tlg.decodeJson(fakeMsgCmd, &msg)
	name := GetCommandName(&msg)
	assert.NoErrorf(t, err, "decodeJson shoudn't return error")
	assert.Equal(t, "mycmd", name, "command doesn't have the same name")

	err = tlg.decodeJson(fakeMsgOther, &msg)
	name = GetCommandName(&msg)
	assert.NoErrorf(t, err, "decodeJson shoudn't return error")
	assert.Empty(t, name, "name shoudn't be empty string")

}

func TestGetInfo(t *testing.T) {
	fakeMsgCmd := []byte(`{"message_id":1053,"from":{"id":669450106,"is_bot":false,"first_name":"\u0443\u03c3\u03c5\u044f\u03b9","username":"wolnosc06","language_code":"en"},"chat":{"id":669450106,"first_name":"\u0443\u03c3\u03c5\u044f\u03b9 ","username":"wolnosc06","type":"private"},"date":1572979909,"text":"/mycmd arg1 arg2   arg3","entities":[{"offset":0,"length":7,"type":"bot_command"}]}`)
	fakeMsgCmdNoTxt := []byte(`{"message_id":1053,"from":{"id":669450106,"is_bot":false,"first_name":"\u0443\u03c3\u03c5\u044f\u03b9","username":"wolnosc06","language_code":"en"},"chat":{"id":669450106,"first_name":"\u0443\u03c3\u03c5\u044f\u03b9 ","username":"wolnosc06","type":"private"},"date":1572979909,"text":"","entities":[{"offset":0,"length":7,"type":"bot_command"}]}`)
	tlg := Telegram{}

	var err error
	var msg MessageTlg

	err = tlg.decodeJson(fakeMsgCmd, &msg)
	name := GetCommandName(&msg)
	assert.NoErrorf(t, err, "decodeJson shoudn't return error")
	assert.Equal(t, "mycmd", name, "command doesn't have the same name")

	expectedCmdInfo := CommandInfo{
		CmdName:      name,
		Msg:          &msg,
		ChatID:       msg.Chat.ID,
		MessageID:    msg.MessageID,
		Restricted:   false,
		isAuthorized: false,
		isAdmin:      false,
		Args:         []string{"arg1", "arg2", "arg3"},
	}

	cmd := Command{
		Enable:            true,
		CmdName:           name,
		NeedAuthorization: false,
	}

	// test with good command data
	actualCmdInfo, err := cmd.GetInfo(&msg)
	assert.NoError(t, err)
	assert.EqualValues(t, expectedCmdInfo, actualCmdInfo)

	// test with empty text command
	err = tlg.decodeJson(fakeMsgCmdNoTxt, &msg)
	name = GetCommandName(&msg)
	assert.NoErrorf(t, err, "decodeJson shoudn't return error")
	expectedCmdInfo.Args = []string{}
	actualCmdInfo, err = cmd.GetInfo(&msg)
	assert.Error(t, err)

}
