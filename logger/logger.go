package logger

import (
	"context"
	"log"
	"log/slog"
	"os"
)

type LogCtx struct {
}

type logWriter struct {
	value []byte
}

type HandlerMiddleware struct {
	next             slog.Handler
	kafkaLogProducer LogProducer
	handlerOptions   *slog.HandlerOptions
}

func InitLogging() {
	if err := CreateLogTopic(); err != nil {
		log.Printf("Failed init kafka logging: %s", err)
	}
	handlerOptions := &slog.HandlerOptions{AddSource: true}
	handler := slog.NewJSONHandler(os.Stdout, handlerOptions)
	slog.SetDefault(slog.New(NewHandlerMiddleware(handler, handlerOptions)))
}

func NewHandlerMiddleware(next slog.Handler, handlerOptions *slog.HandlerOptions) *HandlerMiddleware {
	kafkaLogProducer, err := GetLogProducer()
	if err != nil {
		log.Printf("failed to get kafka log producer in log handler: %s", err)
	}
	return &HandlerMiddleware{next: next, kafkaLogProducer: kafkaLogProducer, handlerOptions: handlerOptions}
}

func (writer *logWriter) Write(p []byte) (n int, err error) {
	writer.value = p
	return os.Stdout.Write(p)
}

func (handler *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return handler.next.Enabled(ctx, rec)
}

func (handler *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	writerForHandler := &logWriter{}
	defaultHandler := slog.NewJSONHandler(writerForHandler, handler.handlerOptions)
	err := defaultHandler.Handle(ctx, rec)
	handler.kafkaLogProducer.ProduceLogToKafka(writerForHandler.value)
	return err
}

func (handler *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: handler.next.WithAttrs(attrs)}
}

func (handler *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{
		next:             handler.next.WithGroup(name),
		kafkaLogProducer: handler.kafkaLogProducer,
		handlerOptions:   handler.handlerOptions}
}
