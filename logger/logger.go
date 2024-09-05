package logger

import (
	"context"
	"log"
	"log/slog"
	"os"
)

type HandlerMiddleware struct {
	next             slog.Handler
	kafkaLogProducer LogProducer
}

func InitLogging() {
	if err := CreateLogTopic(); err != nil {
		log.Printf("Failed init kafka logging: %s", err)
	}
	handler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(NewHandlerMiddleware(handler)))
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	kafkaLogProducer, err := GetLogProducer()
	if err != nil {
		log.Printf("failed to get kafka log producer in log handler: %s", err)
	}
	return &HandlerMiddleware{next: next, kafkaLogProducer: kafkaLogProducer}
}

func (handler *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return handler.next.Enabled(ctx, rec)
}

func (handler *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	//сдесь буду вызывать метод отправки сообщения в кафку
	handler.kafkaLogProducer.ProduceLogToKafka(rec.Message)
	return handler.next.Handle(ctx, rec)
}

func (handler *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: handler.next.WithAttrs(attrs)}
}

func (handler *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{handler.next.WithGroup(name), handler.kafkaLogProducer}
}
