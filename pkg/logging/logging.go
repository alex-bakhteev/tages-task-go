package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"runtime"
)

var loggerInstance *logrus.Logger
var logger *Logger

type Logger struct {
	*logrus.Entry
}

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

// initLogger инициализирует глобальный логгер
func initLogger() {
	l := logrus.New()
	l.SetReportCaller(true)
	l.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}

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
