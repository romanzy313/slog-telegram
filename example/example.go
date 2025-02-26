package main

import (
	"fmt"
	"time"

	"log/slog"

	slogtelegram "github.com/samber/slog-telegram/v2"
)

func main() {
	token := "..."
	chatId := "..."

	logger := slog.New(slogtelegram.Option{Level: slog.LevelDebug, Token: token, ChatId: chatId}.NewTelegramHandler())
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
}
