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

// LogLevel determines the level of logging (priority)
// Separate loggers exists for all types of log levels
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

	DebugColor int = defaultDebugColor
	InfoColor  int = defaultInfoColor
	WarnColor  int = defaultWarnColor
	ErrorColor int = defaultErrorColor
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

// Options provide basic options for tuning logging
type Options struct {
	Level LogLevel
	Debug io.Writer
	Info  io.Writer
	Warn  io.Writer
	Error io.Writer
	Fatal io.Writer
}

func format(colorCode int, level LogLevel) string {
	if colorCode == 0 {
		return string(level)
	}
	return fmt.Sprintf("\x1b[%dm", colorCode) + string(level) + "\x1b[m"
}

func NoColor() {
	DebugColor = 0
	InfoColor = 0
	WarnColor = 0
	ErrorColor = 0
}

// Init initialize default loggers
func Init(opts Options) {

	if opts.Level == "" {
		opts.Level = DefaultLogLevel
	}

	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		NoColor()
	}

	type options struct {
		priority      int
		color         int
		writer        io.Writer
		defaultWriter io.Writer
	}

	params := map[LogLevel]options{
		DebugLevel: {
			priority:      1,
			color:         DebugColor,
			writer:        opts.Debug,
			defaultWriter: os.Stdout,
		},
		InfoLevel: {
			priority:      2,
			color:         InfoColor,
			writer:        opts.Info,
			defaultWriter: os.Stdout,
		},
		WarnLevel: {
			priority:      3,
			color:         WarnColor,
			writer:        opts.Warn,
			defaultWriter: os.Stdout,
		},
		ErrorLevel: {
			priority:      4,
			color:         ErrorColor,
			writer:        opts.Error,
			defaultWriter: os.Stderr,
		},
		FatalLevel: {
			priority:      5,
			color:         ErrorColor,
			writer:        opts.Fatal,
			defaultWriter: os.Stderr,
		},
	}

	lvlPriority := params[opts.Level].priority

	for lvl, cfg := range params {
		loggers[lvl] = New(lvl, ioutil.Discard, "", log.LstdFlags)
		if cfg.priority < lvlPriority {
			continue
		}
		writer := cfg.writer
		if writer == nil {
			writer = cfg.defaultWriter
		}
		loggers[lvl] = New(lvl, writer, fmt.Sprintf("[%v] ", format(cfg.color, lvl)), log.LstdFlags)
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
