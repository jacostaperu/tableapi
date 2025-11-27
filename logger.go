package tableapi

import (
	"fmt"
	"log"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	NONE
)

var levelMap = map[string]LogLevel{
	"DEBUG": DEBUG,
	"INFO":  INFO,
	"WARN":  WARN,
	"ERROR": ERROR,
	"NONE":  NONE,
}

type Logger struct {
	Level       LogLevel
	LevelString string
}

func (l *Logger) logf(level LogLevel, format string, args ...interface{}) {

	if level >= l.Level {
		//		_, file, line, ok := runtime.Caller(1) // 1 to get caller of this function

		// if !ok {
		// 	file = "unknown"
		// 	line = 0
		// }

		//prefix := fmt.Sprintf("%s:%d: ", file, line)
		//log.Output(3, prefix+fmt.Sprintf(format, args...))
		log.Output(3, fmt.Sprintf("%s: ", l.LevelString)+fmt.Sprintf(format, args...))
		// 3 means skip 3 stack frames to get the correct caller: log.Output -> logWithCaller -> caller
		//
		//
		//log.Printf(format, args...)
	}
}

func (l *Logger) Debug(msg string)                          { l.logf(DEBUG, "%s", msg) }
func (l *Logger) Info(msg string)                           { l.logf(INFO, "%s", msg) }
func (l *Logger) Warn(msg string)                           { l.logf(WARN, "%s", msg) }
func (l *Logger) Error(msg string)                          { l.logf(ERROR, "%s", msg) }
func (l *Logger) Debugf(format string, args ...interface{}) { l.logf(DEBUG, format, args...) }
func (l *Logger) Infof(format string, args ...interface{})  { l.logf(INFO, format, args...) }
func (l *Logger) Warnf(format string, args ...interface{})  { l.logf(WARN, format, args...) }
func (l *Logger) Errorf(format string, args ...interface{}) { l.logf(ERROR, format, args...) }

func NewLogger(levelStr string) *Logger {
	log.SetFlags(log.LstdFlags) // no file:line from standard logger

	levelStr = strings.ToUpper(levelStr)
	level, ok := levelMap[levelStr]
	if !ok {
		level = INFO // default level
		levelStr = "INFO"
	}
	// log output can also be customized here if needed
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return &Logger{Level: level, LevelString: levelStr}
}
