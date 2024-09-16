package logger

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"os"
)

// LogCtx данная структура будет храниться в контексте, содержит в себе поля логирования
type LogCtx struct {
	error     error
	cachePath string
	values    string
}

// logWriter это реализация интерфейса io.Writer, передается в дефолтный JSON хендлер
// нужен, чтобы перехватывать логи в json формате
type logWriter struct {
	logKafkaWorker LogKafkaWorker
}

type HandlerMiddleware struct {
	next             slog.Handler
	handlerOptions   *slog.HandlerOptions
	writerForHandler *logWriter
}

// Write отправляет в канал для воркеров кафки полученный JSON лог, затем выводит его в консоль
func (writer logWriter) Write(p []byte) (n int, err error) {
	writer.logKafkaWorker.GetLogChannel() <- string(p)
	return os.Stdout.Write(p)
}

// InitLogging устанавливает мой логгер как дефолтный, включает кастомное логирвоание
func InitLogging() {
	ctx := context.Background()
	if err := CreateLogTopic(); err != nil {
		slog.ErrorContext(ctx, "failed to init kafka logging")
	}
	handlerOptions := &slog.HandlerOptions{AddSource: false, Level: slog.LevelInfo}
	slog.SetDefault(slog.New(NewHandlerMiddleware(handlerOptions)))
}

func WithLogError(ctx context.Context, err error) context.Context {

	if logCtx, ok := ctx.Value("log").(LogCtx); ok {
		logCtx.error = err
		return context.WithValue(ctx, "log", logCtx)
	}

	return context.WithValue(ctx, "log", LogCtx{error: err})
}

func WithLogValues(ctx context.Context, values interface{}) context.Context {
	var valuesString string
	if data, ok := values.(uuid.UUID); ok {
		valuesString = data.String()
	}
	dataJSON, _ := json.Marshal(values)
	valuesString = string(dataJSON)

	if logCtx, ok := ctx.Value("log").(LogCtx); ok {
		logCtx.values = valuesString
		return context.WithValue(ctx, "log", logCtx)
	}

	return context.WithValue(ctx, "log", LogCtx{values: valuesString})
}

func WithLogCacheKey(ctx context.Context, key string) context.Context {

	if logCtx, ok := ctx.Value("log").(LogCtx); ok {
		logCtx.cachePath = key
		return context.WithValue(ctx, "log", logCtx)
	}

	return context.WithValue(ctx, "log", LogCtx{cachePath: key})
}

func NewHandlerMiddleware(handlerOptions *slog.HandlerOptions) *HandlerMiddleware {
	logKafkaWorker := GetLogKafkaWorker(make(chan string, 200))

	writerForHandler := &logWriter{logKafkaWorker: logKafkaWorker}
	next := slog.NewJSONHandler(writerForHandler, handlerOptions)
	handler := &HandlerMiddleware{
		next:             next,
		handlerOptions:   handlerOptions,
		writerForHandler: writerForHandler,
	}
	return handler
}

func (handler *HandlerMiddleware) Enabled(ctx context.Context, rec slog.Level) bool {
	return handler.next.Enabled(ctx, rec)
}

func (handler *HandlerMiddleware) Handle(ctx context.Context, rec slog.Record) error {
	if logCtx, ok := ctx.Value("log").(LogCtx); ok {
		if logCtx.error != nil {
			rec.Add("error", logCtx.error.Error())
		}
		if logCtx.cachePath != "" {
			rec.Add("error", logCtx.cachePath)
		}
		if logCtx.values != "" {
			rec.Add("values", logCtx.values)
		}
	}
	return handler.next.Handle(ctx, rec)
}

func (handler *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddleware{next: handler.next.WithAttrs(attrs)}
}

func (handler *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return &HandlerMiddleware{
		next:             handler.next.WithGroup(name),
		handlerOptions:   handler.handlerOptions,
		writerForHandler: handler.writerForHandler,
	}
}
