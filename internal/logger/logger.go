package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Логирование в stdout
	Log.Out = os.Stdout

	// Установка формата логов
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Установка уровня логирования
	Log.SetLevel(logrus.InfoLevel)
}
