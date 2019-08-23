package tools

import (
	"log"
)

func IsError(err error) (isError bool) {
	if isError = (err != nil); isError {
		log.Print(err)
	}
	return isError
}

type Logger struct {
	name string
}

func NewLogger(name string) Logger {
	return Logger{
		name: name,
	}
}

func (l *Logger) Info(logs interface{}) {
	style := "INFO " + l.name + ": "
	log.Print(style, logs)
}

func (l *Logger) Debug(logs interface{}) {
	style := "DEBUG " + l.name + ": "
	log.Print("********** ********** ********** ********** **********")
	log.Print(style, logs)
	log.Print("********** ********** ********** ********** **********")
}
