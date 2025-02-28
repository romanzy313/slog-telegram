package slogtelegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func checkCredentials(httpClient *http.Client, token, chatId string) error {
	err := getChat(httpClient, token, chatId)

	if err != nil {
		// see if this is a token or a chat issue
		err = getMe(httpClient, token)
		if err != nil {
			return fmt.Errorf("invalid token")
		}

		return fmt.Errorf("invalid chatId")
	}

	return nil
}

// these api methods should really return bool, bool. As networking issues should prevent it from starting up?
// cause i dont want to panic with "invalid token" when http is broken, and not the credentials
// I guess im gonna put this on hold for now...
// docs: https://core.telegram.org/bots/api#getme, https://core.telegram.org/bots/api#authorizing-your-bot
// curl https://api.telegram.org/bot<token>/getMe
func getMe(httpClient *http.Client, token string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getMe", token)

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request to validate token: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token")
	}

	return nil
}

// docs: https://core.telegram.org/bots/api#getchat, https://core.telegram.org/bots/api#authorizing-your-bot
// curl https://api.telegram.org/bot<token>/getChat
//
// Note: this could be used to instead check both the token and chatid
func getChat(httpClient *http.Client, token, chatId string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getChat?chat_id=%s", token, url.QueryEscape(chatId))

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request to validate token: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid token or chatId")
	}

	return nil
}

// docs: https://core.telegram.org/bots/api#sendmessage
//
//	curl -X POST \
//		    -H 'Content-Type: application/json' \
//		    -d '{"chat_id": "<your-chat-id>", "text": "This is a test from curl"}' \
//		    https://api.telegram.org/bot<your-bot-token>/sendMessage
func sendMessage(httpClient *http.Client, token, chatId, msg string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	payload := struct {
		ChatId    string `json:"chat_id"`
		Text      string `json:"text"`
		ParseMode string `json:"parse_mode"`
	}{
		ChatId:    chatId,
		Text:      msg,
		ParseMode: "html", // force parse mode to html as markdown parsing is very picky about formatting
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return fmt.Errorf("failed to marshal JSON: %s", err.Error())
	}

	resp, err := httpClient.Post(url, "application/json", &buf)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to send log: %s", string(errBytes))
	}

	return nil
}

func redactToken(errMsg string, token string) string {
	return strings.ReplaceAll(errMsg, token, "<REDACTED_TOKEN>")
}
