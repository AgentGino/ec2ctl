package logger

import (
	"fmt"
	"log"
	"os"
)

func Log(message string) {
	log.Println(message)
}
func Error(message string, args ...interface{}) {
	// print with red color and emojis
	log.Println("\033[31m" + "❌ " + fmt.Sprintf(message, args...) + "\033[0m")
	os.Exit(1)
}

func Success(message string, args ...interface{}) {
	log.Println("\033[32m" + "✅ " + fmt.Sprintf(message, args...) + "\033[0m")
}

func Warning(message string, args ...interface{}) {
	log.Println("\033[33m" + "⚠️ " + fmt.Sprintf(message, args...) + "\033[0m")
}

func Info(message string, args ...interface{}) {
	log.Println("\033[34m" + "ℹ️ " + fmt.Sprintf(message, args...) + "\033[0m")
}
