package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type LogEntry struct {
	Level   string `json:"level"`
	Message string `json:"message"`
	Time    string `json:"time"`
}

type Logger struct {
	logger *log.Logger
}


func New() *Logger{
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) log(level, message string){
	entry := LogEntry{
		Level: level,
		Message: message,
		Time: time.Now().Format(time.RFC3339),
	}

	data, _ := json.Marshal(entry)
	l.logger.Println(string(data))
}

func (l *Logger) Info(msg string){
	l.log("INFO", msg)
}

func (l *Logger) Error(msg string){
	l.log("ERROR", msg)
}