// Package log provides logging functions
package log

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	TIME_FORMAT = "01-02 15:04:05"
)

var (
	LEVEL_FLAGS  = [...]string{"DEBUG", " INFO", " WARN", "ERROR", "FATAL"}
	globalLogger = NewLogger(os.Stderr, "[P]", true, false)
)

type Logger struct {
	Writer      io.Writer
	Prefix      string
	NonColor    bool
	ShowDepth   bool
	CallerDepth int
	Level       int
}

func NewLogger(w io.Writer, prefix string, enableColor bool, enableDepth bool) *Logger {
	l := new(Logger)
	l.Writer = w
	l.Prefix = prefix
	l.NonColor = !enableColor
	l.Level = DEBUG

	// 颜色输出不支持windows
	if runtime.GOOS == "windows" {
		l.NonColor = true
	}

	l.ShowDepth = enableDepth
	l.CallerDepth = 2
	return l
}

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

func (lg *Logger) Print(level int, format string, args ...interface{}) {
	if level < lg.Level {
		return
	}
	var depthInfo string
	if lg.ShowDepth {
		pc, file, line, ok := runtime.Caller(lg.CallerDepth)
		if ok {
			// Get caller function name.
			fn := runtime.FuncForPC(pc)
			var fnName string
			if fn == nil {
				fnName = "?()"
			} else {
				fnName = strings.TrimLeft(filepath.Ext(fn.Name()), ".") + "()"
			}
			depthInfo = fmt.Sprintf("[%s:%d %s] ", filepath.Base(file), line, fnName)
		}
	}
	if lg.NonColor {
		fmt.Fprintf(lg.Writer, "%s %s [%s] %s%s\n",
			lg.Prefix, time.Now().Format(TIME_FORMAT), LEVEL_FLAGS[level], depthInfo,
			fmt.Sprintf(format, args...))
		if level == FATAL {
			os.Exit(1)
		}
		return
	}

	switch level {
	case DEBUG:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m [\033[34m%s\033[0m] %s%s\n",
			lg.Prefix, time.Now().Format(TIME_FORMAT), LEVEL_FLAGS[level], depthInfo,
			fmt.Sprintf(format, args...))
	case INFO:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m [\033[32m%s\033[0m] %s%s\n",
			lg.Prefix, time.Now().Format(TIME_FORMAT), LEVEL_FLAGS[level], depthInfo,
			fmt.Sprintf(format, args...))
	case WARNING:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m [\033[33m%s\033[0m] %s%s\n",
			lg.Prefix, time.Now().Format(TIME_FORMAT), LEVEL_FLAGS[level], depthInfo,
			fmt.Sprintf(format, args...))
	case ERROR:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m [\033[31m%s\033[0m] %s%s\n",
			lg.Prefix, time.Now().Format(TIME_FORMAT), LEVEL_FLAGS[level], depthInfo,
			fmt.Sprintf(format, args...))
	case FATAL:
		fmt.Fprintf(lg.Writer, "%s \033[36m%s\033[0m [\033[35m%s\033[0m] %s%s\n",
			lg.Prefix, time.Now().Format(TIME_FORMAT), LEVEL_FLAGS[level], depthInfo,
			fmt.Sprintf(format, args...))
		os.Exit(1)
	default:
		fmt.Fprintf(lg.Writer, "%s %s [%s] %s%s\n",
			lg.Prefix, time.Now().Format(TIME_FORMAT), LEVEL_FLAGS[level], depthInfo,
			fmt.Sprintf(format, args...))
	}
}

func Debug(format string, args ...interface{}) {
	globalLogger.Print(DEBUG, format, args...)
}

func Warn(format string, args ...interface{}) {
	globalLogger.Print(WARNING, format, args...)
}

func Info(format string, args ...interface{}) {
	globalLogger.Print(INFO, format, args...)
}

func Error(format string, args ...interface{}) {
	globalLogger.Print(ERROR, format, args...)
}

func Fatal(format string, args ...interface{}) {
	globalLogger.Print(FATAL, format, args...)
}

func Get() *Logger {
	return globalLogger
}
