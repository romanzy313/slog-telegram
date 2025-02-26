package slogtelegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// does a request to see if token is correct
// curl https://api.telegram.org/bot<token>/getMe
func (o *Option) checkInit() error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", o.Token)

	return o.makeRequest("GET", url, "")
}

// docs: https://core.telegram.org/bots/api#sendmessage
//
// manually tested with
//
//	curl -X POST \
//	     -H 'Content-Type: application/json' \
//	     -d '{"chat_id": "<your-chat-id>", "text": "This is a test from curl"}' \
//	     https://api.telegram.org/bot<your-bot-token>/sendMessage
func (o *Option) sendMessage(msg string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", o.Token)

	value := struct {
		ChatId    string `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode,omitempty"`
	}{
		ChatId:    o.ChatId,
		Text:      msg,
		ParseMode: o.ParseMode,
	}

	bytes, err := json.Marshal(value)

	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err.Error())
	}

	return o.makeRequest("POST", url, string(bytes))
}

func (o *Option) makeRequest(method, url, body string) (err error) {
	var resp *http.Response
	if method == "GET" {
		resp, err = o.HttpClient.Get(url)
		if err != nil {
			return fmt.Errorf("failed to make http request: %s", err.Error())
		}
	} else if method == "POST" {
		resp, err = o.HttpClient.Post(url, "application/json", strings.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to make http request: %s", err.Error())
		}
	} else {
		panic("unreacheable")
	}

	if resp.StatusCode != http.StatusOK {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("failed to send log: [%s] %s", resp.Status, string(errBytes))
	}

	return nil
}
