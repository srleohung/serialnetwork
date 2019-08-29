package tools

import (
	"log"
)

type Logger struct {
	name     string
	logLevel map[logLevel]bool
}

type logLevel int

const (
	EMERG   logLevel = 0 // The system is unusable.
	ALERT   logLevel = 1 // Actions that must be taken care of immediately.
	CRIT    logLevel = 2 // Critical conditions.
	ERR     logLevel = 3 // Noncritical error conditions.
	WARNING logLevel = 4 // Warning conditions that should be taken care of.
	NOTICE  logLevel = 5 // Normal, but significant events.
	INFO    logLevel = 6 // Informational messages that require no action.
	DEBUG   logLevel = 7
)

var defaultLogLevel map[logLevel]bool = map[logLevel]bool{
	EMERG:   true,
	ALERT:   true,
	CRIT:    true,
	ERR:     false,
	WARNING: true,
	NOTICE:  true,
	INFO:    true,
	DEBUG:   false,
}

func NewLogger(name string) Logger {
	return Logger{
		name:     name,
		logLevel: defaultLogLevel,
	}
}

func (l *Logger) SetLogLevel(logLevel map[logLevel]bool) {
	l.logLevel = logLevel
}

func (l *Logger) Debug(logs ...interface{}) {
	if l.logLevel[DEBUG] {
		style := "DEBUG " + l.name + ": "
		log.Print(style, logs)
	}
}

func (l *Logger) Info(logs ...interface{}) {
	if l.logLevel[INFO] {
		style := "INFO " + l.name + ": "
		log.Print(style, logs)
	}
}

func (l *Logger) Notice(logs ...interface{}) {
	if l.logLevel[NOTICE] {
		style := "NOTICE " + l.name + ": "
		log.Print(style, logs)
	}
}

func (l *Logger) Warning(logs ...interface{}) {
	if l.logLevel[WARNING] {
		style := "WARNING " + l.name + ": "
		log.Print(style, logs)
	}
}

func (l *Logger) Err(logs ...interface{}) {
	if l.logLevel[ERR] {
		style := "ERR " + l.name + ": "
		log.Print(style, logs)
	}
}

func (l *Logger) Crit(logs ...interface{}) {
	if l.logLevel[CRIT] {
		style := "CRIT " + l.name + ": "
		log.Print(style, logs)
	}
}

func (l *Logger) Alert(logs ...interface{}) {
	if l.logLevel[ALERT] {
		style := "ALERT " + l.name + ": "
		log.Print(style, logs)
	}
}

func (l *Logger) Emerg(logs ...interface{}) {
	if l.logLevel[EMERG] {
		style := "EMERG " + l.name + ": "
		log.Fatal(style, logs)
	}
}

func (l *Logger) Debugf(format string, logs ...interface{}) {
	if l.logLevel[DEBUG] {
		style := "DEBUG " + l.name + ": "
		log.Printf(style+format, logs)
	}
}

func (l *Logger) Infof(format string, logs ...interface{}) {
	if l.logLevel[INFO] {
		style := "INFO " + l.name + ": "
		log.Printf(style+format, logs)
	}
}

func (l *Logger) Noticef(format string, logs ...interface{}) {
	if l.logLevel[NOTICE] {
		style := "NOTICE " + l.name + ": "
		log.Printf(style+format, logs)
	}
}

func (l *Logger) Warningf(format string, logs ...interface{}) {
	if l.logLevel[WARNING] {
		style := "WARNING " + l.name + ": "
		log.Printf(style+format, logs)
	}
}

func (l *Logger) Errf(format string, logs ...interface{}) {
	if l.logLevel[ERR] {
		style := "ERR " + l.name + ": "
		log.Printf(style+format, logs)
	}
}

func (l *Logger) Critf(format string, logs ...interface{}) {
	if l.logLevel[CRIT] {
		style := "CRIT " + l.name + ": "
		log.Printf(style+format, logs)
	}
}

func (l *Logger) Alertf(format string, logs ...interface{}) {
	if l.logLevel[ALERT] {
		style := "ALERT " + l.name + ": "
		log.Printf(style+format, logs)
	}
}

func (l *Logger) Emergf(format string, logs ...interface{}) {
	if l.logLevel[EMERG] {
		style := "EMERG " + l.name + ": "
		log.Fatalf(style+format, logs)
	}
}

func (l *Logger) IsErr(err error) (isErr bool) {
	if isErr = (err != nil); isErr {
		l.Err(err)
	}
	return isErr
}
