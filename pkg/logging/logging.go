package logging

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
)

var loggerInstance *logrus.Logger
var logger *Logger

type Logger struct {
	*logrus.Entry
}

type contextKey string

// Context key для x-request-id
const RequestIDKey = contextKey("x-request-id")

// Context key для хранения логгера
const LoggerKey = contextKey("logger")

// GetLogger возвращает экземпляр логгера
func GetLogger() *Logger {
	if logger == nil {
		initLogger()
	}
	return logger
}

func (l *Logger) GetLoggerWithField(key string, value interface{}) *Logger {
	return &Logger{l.Entry.WithField(key, value)}
}

// GetLoggerWithContext добавляет x-request-id из контекста в логгер
func (l *Logger) GetLoggerWithContext(ctx context.Context) *Logger {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}
	return l.GetLoggerWithField("x-request-id", requestID)
}

func AddLoggerToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, logger)
}

// WithRequestID добавляет x-request-id в контекст
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// initLogger инициализирует глобальный логгер
func initLogger() {
	l := logrus.New()

	// Создание каталога logs, если он не существует
	if err := os.MkdirAll("logs", 0755); err != nil {
		fmt.Printf("Ошибка создания каталога logs: %v\n", err)
	}

	// Открытие/создание файла для записи логов
	allFile, err := os.OpenFile("logs/all.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Ошибка открытия файла логов: %v\n", err)
		return
	}

	// Настройка вывода логов в файл и консоль
	multiWriter := io.MultiWriter(allFile, os.Stdout)
	l.SetOutput(multiWriter)

	// Установка уровня логирования
	l.SetLevel(logrus.TraceLevel)

	loggerInstance = l
	logger = &Logger{logrus.NewEntry(loggerInstance)}
	fmt.Println("Логгер инициализирован")
}

func (l *Logger) InfoCtx(ctx context.Context, args ...interface{}) {
	loggerWithContext := l.GetLoggerWithContext(ctx)
	entry := loggerWithContext.Entry.WithFields(logrus.Fields{
		"caller": getCaller(2), // Смещение для корректного отображения места вызова
	})
	entry.Info(args...)
}

func (l *Logger) ErrorCtx(ctx context.Context, args ...interface{}) {
	loggerWithContext := l.GetLoggerWithContext(ctx)
	entry := loggerWithContext.Entry.WithFields(logrus.Fields{
		"caller": getCaller(2),
	})
	entry.Error(args...)
}

func (l *Logger) DebugCtx(ctx context.Context, args ...interface{}) {
	loggerWithContext := l.GetLoggerWithContext(ctx)
	entry := loggerWithContext.Entry.WithFields(logrus.Fields{
		"caller": getCaller(2),
	})
	entry.Debug(args...)
}

func getCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown:0"
	}

	// Вычитаем базовый путь проекта (если известен)
	basePath, err := os.Getwd() // Путь до корня проекта
	if err != nil {
		return fmt.Sprintf("%s:%d", path.Base(file), line)
	}

	relativePath := file
	if strings.HasPrefix(file, basePath) {
		relativePath = strings.TrimPrefix(file, basePath+"/")
	}

	funcName := runtime.FuncForPC(pc).Name()
	return fmt.Sprintf("%s:%d %s", relativePath, line, path.Base(funcName))
}
