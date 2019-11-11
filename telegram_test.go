package telegram

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindCommand(t *testing.T) {

}

func TestSendMessage(t *testing.T) {

	// msg to send
	msg := url.QueryEscape("Test message!")
	chatid := "123456789"
	data200 := []byte(`{"ok":true,"result":{"message_id":828,"from":{"id":789456,"is_bot":true,"first_name":"Tester","username":"tester_bot"},"chat":{"id":` + chatid + `,"first_name":"\u043c\u03b1\u044fg\u03b1\u044f\u03b9\u03b7\u0454","username":"PowerUser","type":"private"},"date":1572617995,"text":"` + msg + `"}}`)
	data400 := []byte(`{"ok":false,"error_code":400, "description":"Bad request"}`)
	data404 := []byte(`{"ok":false,"error_code":404, "description":"Not Found"}`)
	dataBad := []byte(`'"l(kjrlkjtml"thmflkhg'"(`)
	// fake server emulate telegram API
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/sendMessage" {
			w.WriteHeader(http.StatusNotFound)
			w.Write(data404)
			return
		}
		r.ParseForm()
		chatid := r.Form.Get("chat_id")
		text := r.Form.Get("text")
		// require fields (tlg api doc)
		if chatid == "" || text == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(data400)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data200)
	}))
	defer server.Close()

	tlg := Telegram{
		BotURL: server.URL + "/",
	}

	var err error
	// normal
	err = tlg.SendChatMessage(123456, "test")
	assert.NoError(t, err, "shouldn't return error")

	// bad return from server
	data200 = dataBad
	err = tlg.SendChatMessage(123456, "test")
	assert.Error(t, err, "Shoud return error")

	// bad address
	tlg.BotURL = "127.0.0.2:65555/"
	err = tlg.SendChatMessage(123456, "test")
	assert.Error(t, err, "Shoud return error")
}

func TestDecodeJson(t *testing.T) {
	updateid := "123123"
	updateidInt, _ := strconv.Atoi(updateid)
	dataGood := []byte(`{"ok":true,"result":[{"update_id":` + updateid + `,"message":{"message_id":830,"from":{"id":789456,"is_bot":false,"first_name":"мαяgαяιηє","username":"Testname","language_code":"en"},"chat":{"id":123456,"first_name":"мαяgαяιηє","username":"Testname","type":"private"},"date":1572619731,"text":"boom"}}]}`)
	dataBad1 := []byte(`"ok":true,"result":[{"update_id":` + updateid + `,"message":{"message_id":830,"from":{"id":789456,"is_bot":false,"first_name":"мαяgαяιηє","username":"Testname","language_code":"en"},"chat":{"id":123456,"first_name":"мαяgαяιηє","username":"Testname","type":"private"},"date":1572619731,"text":"boom"}}]}`)
	dataBad2 := []byte(`élkjlkjmlkhfsd`)

	tlg := Telegram{}
	d := &ReceiveTlg{}

	var err error
	err = tlg.decodeJson(dataGood, d)
	assert.NoError(t, err, "shouldn't return error")
	assert.Equal(t, updateidInt, d.Result[0].UpdateId, "update_id not equal")

	err = tlg.decodeJson(dataBad1, d)
	assert.Error(t, err, "should return error")

	err = tlg.decodeJson(dataBad2, d)
	assert.Error(t, err, "should return error")
}

func TestPopMessage(t *testing.T) {
	var updateTlg []UpdateTlg
	tlg := Telegram{}
	updateid := 123456

	// create 10 fakes results
	for i := 0; i < 10; i++ {
		result := UpdateTlg{}
		iStr := strconv.Itoa(i)
		currentId := strconv.Itoa(updateid + i)
		d := []byte(`{"update_id":` + currentId + `,"message":{"message_id":` + iStr + `,"from":{"id":789456,"is_bot":false,"first_name":"мαяgαяιηє","username":"Testname","language_code":"en"},"chat":{"id":123456,"first_name":"мαяgαяιηє","username":"Testname","type":"private"},"date":1572619731,"text":"boom"}}`)
		err := tlg.decodeJson(d, &result)
		assert.NoError(t, err, "shouldn't return error")
		updateTlg = append(updateTlg, result)
	}

	// preserve original struct for expected values
	updateTlgOrig := make([]UpdateTlg, len(updateTlg))
	copy(updateTlgOrig, updateTlg)
	recv := ReceiveTlg{
		Ok:     true,
		Result: updateTlg,
	}

	i := 0
	for {
		id, msg, valid := tlg.popMessage(&recv)
		if !valid {
			break
		}
		assert.Equal(t, updateTlgOrig[i].UpdateId, id)
		assert.Equal(t, msg, updateTlgOrig[i].Message)
		assert.Equal(t, msg.MessageID, updateTlgOrig[i].Message.MessageID)
		i++
	}
}
