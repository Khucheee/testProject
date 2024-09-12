package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"time"
)

func withLogging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//считываем тело запроса
		requestBody, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			slog.Error("failed to read request body", "error", err)
		}
		// Восстанавливаем тело запроса после чтения (чтобы оно было доступно для последующих обработчиков)
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

		//создаем кастомный writer, чтобы перехватить тело ответа
		cWriter := &customWriter{
			originalWriter: ctx.Writer,
			responseBody:   bytes.NewBufferString(""),
		}
		ctx.Writer = cWriter

		//обрабатываем оригинальный запрос
		startHandlingTime := time.Now()
		ctx.Next()
		handlingTime := time.Since(startHandlingTime)

		//логируем
		slog.Info(
			"incoming request",
			"Method", ctx.Request.Method,
			"host", ctx.Request.Host,
			"path", ctx.Request.URL.Path,
			"client_IP", ctx.ClientIP(),
			"request_handling_time", handlingTime,
			"request_body", string(requestBody),
			"response_body", cWriter.responseBody.String(),
			"status_code", cWriter.statusCode,
		)
	}
}
