package logger

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

// Init инициализирует логгер
func Init(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	Logger = log.New(file, "", log.LstdFlags)
	return nil
}

// Info логирует информационное сообщение
func Info(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[INFO] "+format, v...)
	}
}

// Error логирует сообщение об ошибке
func Error(format string, v ...interface{}) {
	if Logger != nil {
		Logger.Printf("[ERROR] "+format, v...)
	}
}

// Close закрывает файл логов
func Close() error {
	if Logger != nil {
		return Logger.Writer().(*os.File).Close()
	}
	return nil
}
