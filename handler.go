package slogtelegram

import (
	"context"
	"net/http"

	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

type Option struct {
	// log level (default: debug)
	Level slog.Leveler

	// Telegram bot token
	Token string
	// ChatId is the id of the chat
	ChatId string

	Converter Converter

	// optional: see slog.HandlerOptions
	AddSource   bool
	ReplaceAttr func(groups []string, a slog.Attr) slog.Attr

	// optional: customize HTTP client
	HttpClient *http.Client
}

func (o Option) NewTelegramHandler() slog.Handler {
	if o.Level == nil {
		o.Level = slog.LevelDebug
	}

	if o.Token == "" {
		panic("missing Telegram token")
	}

	if o.ChatId == "" {
		panic("missing Telegram chatId")
	}

	if o.Converter == nil {
		o.Converter = DefaultConverter
	}

	if o.HttpClient == nil {
		o.HttpClient = http.DefaultClient
	}

	// err := checkToken(o.HttpClient, o.Token)
	err := checkCredentials(o.HttpClient, o.Token, o.ChatId)
	if err != nil {
		panic(redactToken(err.Error(), o.Token))
	}

	return &TelegramHandler{
		option:     o,
		attrs:      []slog.Attr{},
		groups:     []string{},
		httpClient: o.HttpClient,
	}
}

var _ slog.Handler = (*TelegramHandler)(nil)

type TelegramHandler struct {
	option     Option
	attrs      []slog.Attr
	groups     []string
	httpClient *http.Client
}

func (h *TelegramHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.option.Level.Level()
}

func (h *TelegramHandler) Handle(ctx context.Context, record slog.Record) error {
	msg := h.option.Converter(h.option.AddSource, h.option.ReplaceAttr, h.attrs, h.groups, &record)

	// non-blocking
	go func() {
		o := h.option
		_ = sendMessage(o.HttpClient, o.Token, o.ChatId, msg)
	}()

	return nil
}

func (h *TelegramHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TelegramHandler{
		option: h.option,
		attrs:  slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...),
		groups: h.groups,
	}
}

func (h *TelegramHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return &TelegramHandler{
		option: h.option,
		attrs:  h.attrs,
		groups: append(h.groups, name),
	}
}

func (h *TelegramHandler) WithHttpClient(client *http.Client) slog.Handler {
	return &TelegramHandler{
		option:     h.option,
		attrs:      h.attrs,
		groups:     h.groups,
		httpClient: client,
	}
}
