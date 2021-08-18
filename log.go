package nanolog

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

// LogLevel determines the level of logging (priority)
// Separate loggers exists for all types of log levels
type LogLevel string

// passthrough constants from standart log package
const (
	Ldate         = log.Ldate
	Ltime         = log.Ltime
	Lmicroseconds = log.Lmicroseconds
	Llongfile     = log.Llongfile
	Lshortfile    = log.Lshortfile
	LUTC          = log.LUTC
	LstdFlags     = log.LstdFlags
)

const (
	DebugLevel      LogLevel = "DEBUG"
	InfoLevel       LogLevel = "INFO"
	WarnLevel       LogLevel = "WARN"
	ErrorLevel      LogLevel = "ERROR"
	FatalLevel      LogLevel = "FATAL"
	DefaultLogLevel          = ErrorLevel
	DefaultPrefix            = "[%v] "
	DefaultFlags             = log.LstdFlags
)

const (
	defaultDebugColor = 32
	defaultInfoColor  = 35
	defaultWarnColor  = 33
	defaultErrorColor = 31
)

var (
	DebugColor = defaultDebugColor
	InfoColor  = defaultInfoColor
	WarnColor  = defaultWarnColor
	ErrorColor = defaultErrorColor
)

var mutex = &sync.Mutex{}
var loggers = createLoggers(Options{Level: DefaultLogLevel})

// Logger wrapper for standart log.Logger
type Logger struct {
	level  LogLevel
	prefix string
	flags  int
	logger *log.Logger
}

func New(lvl LogLevel, writer io.Writer, prefix string, flags int) *Logger {
	return &Logger{
		level:  lvl,
		prefix: prefix,
		flags:  flags,
		logger: log.New(writer, prefix, flags),
	}
}

func (l *Logger) Println(args ...interface{}) {
	if l.level == FatalLevel {
		l.logger.Fatalln(args...)
		return
	}
	l.logger.Println(args...)
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	if l.level == FatalLevel {
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

type LoggerOptions struct {
	Writer io.Writer
	Color  int
	Prefix string
	Flags  int
}

type internalOptions struct {
	writer   io.Writer
	color    int
	prefix   string
	flags    int
	priority int
}

// Options provide basic options for tuning logging
type Options struct {
	Level LogLevel
	Debug LoggerOptions
	Info  LoggerOptions
	Warn  LoggerOptions
	Error LoggerOptions
	Fatal LoggerOptions
}

func FormatPrefix(prefix string, colorCode int, level LogLevel) string {

	if colorCode == 0 {
		return string(level)
	}
	formattedPrefix := fmt.Sprintf("\x1b[%dm%s\x1b[m", colorCode, level)
	return fmt.Sprintf(prefix, formattedPrefix)
}

// NoColor turn off colored output
func NoColor() {
	DebugColor = 0
	InfoColor = 0
	WarnColor = 0
	ErrorColor = 0
}

func internalDefaults(priority int, color int) internalOptions {

	return internalOptions{
		priority: priority,
		color:    color,
		writer:   os.Stdout,
		prefix:   DefaultPrefix,
		flags:    DefaultFlags,
	}
}

func getDefaultOptions() map[LogLevel]internalOptions {

	return map[LogLevel]internalOptions{
		DebugLevel: internalDefaults(1, DebugColor),
		InfoLevel:  internalDefaults(2, InfoColor),
		WarnLevel:  internalDefaults(3, WarnColor),
		ErrorLevel: internalDefaults(4, ErrorColor),
		FatalLevel: internalDefaults(5, ErrorColor),
	}
}

func createLoggers(opts Options) map[LogLevel]*Logger {

	mutex.Lock()
	defer mutex.Unlock()

	loggers := make(map[LogLevel]*Logger)

	if opts.Level == "" {
		opts.Level = DefaultLogLevel
	}

	if runtime.GOOS != "linux" && runtime.GOOS != "darwin" {
		NoColor()
	}

	optionsOverride := map[LogLevel]LoggerOptions{
		DebugLevel: opts.Debug,
		InfoLevel:  opts.Info,
		WarnLevel:  opts.Warn,
		ErrorLevel: opts.Error,
		FatalLevel: opts.Fatal,
	}

	defaultOptions := getDefaultOptions()

	lvlPriority := defaultOptions[opts.Level].priority

	for lvl, defaultOptions := range defaultOptions {
		if defaultOptions.priority < lvlPriority {
			loggers[lvl] = New(lvl, ioutil.Discard, "", log.LstdFlags)
			continue
		}
		override := optionsOverride[lvl]
		writer := override.Writer
		color := override.Color
		prefix := override.Prefix
		flags := override.Flags

		if writer == nil {
			writer = defaultOptions.writer
		}
		if color == 0 {
			color = defaultOptions.color
		}
		if prefix == "" {
			prefix = defaultOptions.prefix
		}
		if flags == 0 {
			flags = defaultOptions.flags
		}
		if flags == -1 {
			flags = 0
		}

		prefixFormatted := FormatPrefix(prefix, color, lvl)
		loggers[lvl] = New(lvl, writer, prefixFormatted, flags)
	}

	return loggers
}

// Init initialize default loggers
func Init(opts Options) {
	loggers = createLoggers(opts)
}

func getLogger(lvl LogLevel) *Logger {
	mutex.Lock()
	defer mutex.Unlock()
	return loggers[lvl]
}

// Log logs arguments
func Log(lvl LogLevel, args ...interface{}) {
	getLogger(lvl).Println(args...)
}

// Logf logs arguments with formatting
func Logf(lvl LogLevel, msg string, args ...interface{}) {
	getLogger(lvl).Printf(msg, args...)
}

func Debug() *Logger {
	return getLogger(DebugLevel)
}

func Info() *Logger {
	return getLogger(InfoLevel)
}

func Warn() *Logger {
	return getLogger(WarnLevel)
}

func Error() *Logger {
	return getLogger(ErrorLevel)
}

func Fatal() *Logger {
	return getLogger(FatalLevel)
}

type NanoLogger struct {
}

func (l *NanoLogger) Debug() *Logger {
	return Debug()
}

func (l *NanoLogger) Info() *Logger {
	return Info()
}

func (l *NanoLogger) Warn() *Logger {
	return Warn()
}

func (l *NanoLogger) Error() *Logger {
	return Error()
}

func (l *NanoLogger) Fatal() *Logger {
	return Fatal()
}

var std = &NanoLogger{}

func DefaultLogger() *NanoLogger {
	return std
}
