package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	// fake good data returned from Tlg API
	dataGood := []byte(`{"ok":true,"result":[{"update_id":123123,"message":{"message_id":830,"from":{"id":789456,"is_bot":false,"first_name":"мαяgαяιηє","username":"Testname","language_code":"en"},"chat":{"id":123456,"first_name":"мαяgαяιηє","username":"Testname","type":"private"},"date":1572619731,"text":"boom"}}]}`)

	tlg := Telegram{}

	// decode to ReceiveTlg struct
	recv := ReceiveTlg{}
	err := tlg.decodeJson(dataGood, &recv)

	assert.NoErrorf(t, err, "shouldn't return error")
	assert.Equalf(t, 1, len(recv.Result), "results slice must have 1 elem")

	user, err := NewUser(&recv.Result[0].Message)
	assert.NoError(t, err, "NewUser() shoudn't return error")

	userOrig := recv.Result[0].Message.From

	assert.Equal(t, userOrig.ID, user.Id)
	assert.Equal(t, userOrig.FirstName, user.FirstName)
	assert.Equal(t, userOrig.LastName, user.LastName)
	assert.Equal(t, userOrig.UserName, user.UserName)
	assert.Equal(t, userOrig.IsBot, user.isBot)

}
