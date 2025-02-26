//go:build integration
// +build integration

package slogtelegram

// run with TOKEN=<token> CHAT_ID=<chatId> go test -tags=integration

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic as expected")
		}
	}()
	f()
}

func TestPanicWithInvalidToken(t *testing.T) {
	assertPanic(t, func() {
		Option{Token: "bad-token", ChatId: os.Getenv("CHAT_ID")}.NewTelegramHandler()
	})
}

// currently does not work. need to get the full chat in order to verify that chat id is invalid
// func TestPanicWithInvalidChatId(t *testing.T) {
// 	assertPanic(t, func() {
// 		Option{Token: os.Getenv("TOKEN"), ChatId: "lala"}.NewTelegramHandler()
// 	})
// }

func TestSuccessfulInit(t *testing.T) {
	handler := Option{Token: os.Getenv("TOKEN"), ChatId: os.Getenv("CHAT_ID")}.NewTelegramHandler()
	if handler == nil {
		t.Errorf("Expected handler to be non-nil for valid credentials, got %v", handler)
	}
}

// this test "leaks goroutines"
func TestSendMessage(t *testing.T) {
	token := os.Getenv("TOKEN")
	chatId := os.Getenv("CHAT_ID")

	logger := slog.New(Option{Level: slog.LevelDebug, Token: token, ChatId: chatId}.NewTelegramHandler())
	logger = logger.With("release", "v1.0.0")

	logger.
		With(
			slog.Group("user",
				slog.String("id", "user-123"),
				slog.Time("created_at", time.Now().AddDate(0, 0, -1)),
			),
		).
		With("environment", "dev").
		With("error", fmt.Errorf("an error")).
		Error("Hello slog")

	time.Sleep(1 * time.Second)
}
