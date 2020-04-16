package nanolog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
)

// LogLevel determine level of logging
type LogLevel string

const (
	DebugLevel      LogLevel = "DEBUG"
	InfoLevel       LogLevel = "INFO"
	WarnLevel       LogLevel = "WARN"
	ErrorLevel      LogLevel = "ERROR"
	FatalLevel      LogLevel = "FATAL"
	DefaultLogLevel          = ErrorLevel
)

const (
	defaultDebugColor = 32
	defaultInfoColor  = 35
	defaultWarnColor  = 33
	defaultErrorColor = 31
)

var (
	loggers = map[LogLevel]*Logger{}
	DEBUG   *Logger
	INFO    *Logger
	WARN    *Logger
	ERROR   *Logger
	FATAL   *Logger

	DebugColor string = fmt.Sprintf("\x1b[%dm", defaultDebugColor)
	InfoColor  string = fmt.Sprintf("\x1b[%dm", defaultInfoColor)
	WarnColor  string = fmt.Sprintf("\x1b[%dm", defaultWarnColor)
	ErrorColor string = fmt.Sprintf("\x1b[%dm", defaultErrorColor)
)

// Logger wrapper for standart log.Logger
type Logger struct {
	Level  LogLevel
	Prefix string
	Flags  int
	logger *log.Logger
}

func New(lvl LogLevel, writer io.Writer, prefix string, flags int) *Logger {
	return &Logger{
		Level:  lvl,
		Prefix: prefix,
		Flags:  flags,
		logger: log.New(writer, prefix, flags),
	}
}

func (l *Logger) Println(args ...interface{}) {
	if l.Level == FatalLevel {
		l.logger.Fatalln(args...)
		return
	}
	l.logger.Println(args...)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	if l.Level == FatalLevel {
		l.logger.Fatalf(msg, args...)
		return
	}
	l.logger.Printf(msg, args...)
}

// ParseLevel parse a string representation of log level and returns equivalent LogLevel
// useful when set value from config
func ParseLevel(str string) LogLevel {
	logLevel := strings.ToLower(str)
	return map[string]LogLevel{
		"debug": DebugLevel,
		"info":  InfoLevel,
		"warn":  WarnLevel,
		"error": ErrorLevel,
		"fatal": FatalLevel,
	}[logLevel]
}

// LoggingOptions provide basic options for tuning logging
type Options struct {
	Level LogLevel
	Debug io.Writer
	Info  io.Writer
	Warn  io.Writer
	Error io.Writer
	Fatal io.Writer
}

func format(colorCode string, level LogLevel) string {
	if colorCode == "" {
		return string(level)
	}
	return colorCode + string(level) + "\x1b[m"
}

// Init initialize default loggers
func Init(opts Options) {

	if opts.Level == "" {
		opts.Level = DefaultLogLevel
	}

	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		DebugColor = ""
		InfoColor = ""
		WarnColor = ""
		ErrorColor = ""
	}

	logPriority := map[LogLevel]int{
		DebugLevel: 1,
		InfoLevel:  2,
		WarnLevel:  3,
		ErrorLevel: 4,
		FatalLevel: 5,
	}

	lvlPriority := logPriority[opts.Level]

	for lvl, priority := range logPriority {

		var writer io.Writer
		var color string

		loggers[lvl] = New(lvl, ioutil.Discard, "", log.LstdFlags)

		if priority < lvlPriority {
			continue
		}

		switch lvl {
		case DebugLevel:
			writer = opts.Debug
			if writer == nil {
				writer = os.Stdout
			}
			color = DebugColor
		case InfoLevel:
			writer = opts.Info
			if writer == nil {
				writer = os.Stdout
			}
			color = InfoColor
		case WarnLevel:
			writer = opts.Warn
			if writer == nil {
				writer = os.Stdout
			}
			color = WarnColor
		case ErrorLevel:
			writer = opts.Error
			if writer == nil {
				writer = os.Stderr
			}
			color = ErrorColor
		case FatalLevel:
			writer = opts.Fatal
			if writer == nil {
				writer = os.Stderr
			}
			color = ErrorColor
		}
		loggers[lvl] = New(lvl, writer, fmt.Sprintf("[%v] ", format(color, lvl)), log.LstdFlags)
	}

	DEBUG = loggers[DebugLevel]
	INFO = loggers[InfoLevel]
	WARN = loggers[WarnLevel]
	ERROR = loggers[ErrorLevel]
	FATAL = loggers[FatalLevel]
}

// Log logs arguments
func Log(lvl LogLevel, args ...interface{}) {
	loggers[lvl].Println(args...)
}

// Logf logs arguments with formatting
func Logf(lvl LogLevel, msg string, args ...interface{}) {
	loggers[lvl].Printf(msg, args...)
}
