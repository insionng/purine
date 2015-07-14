package log

import "strings"

// implement to tango.Logger, so receive tango's logs
type TangoLogger struct {
	log *Logger
}

func (l *Logger) ToTangoLogger() *TangoLogger {
	lg := &TangoLogger{l}
	return lg
}

func (tl *TangoLogger) Debugf(format string, v ...interface{}) {
	tl.log.Print(DEBUG, format, v...)
}

func (tl *TangoLogger) Debug(v ...interface{}) {
	format := strings.Repeat("%v ", len(v))
	tl.log.Print(DEBUG, format, v...)
}

func (tl *TangoLogger) Infof(format string, v ...interface{}) {
	tl.log.Print(INFO, format, v...)
}

func (tl *TangoLogger) Info(v ...interface{}) {
	format := strings.Repeat("%v ", len(v))
	tl.log.Print(INFO, format, v...)
}

func (tl *TangoLogger) Warnf(format string, v ...interface{}) {
	tl.log.Print(WARNING, format, v...)
}

func (tl *TangoLogger) Warn(v ...interface{}) {
	format := strings.Repeat("%v ", len(v))
	tl.log.Print(WARNING, format, v...)
}

func (tl *TangoLogger) Errorf(format string, v ...interface{}) {
	tl.log.Print(ERROR, format, v...)
}

func (tl *TangoLogger) Error(v ...interface{}) {
	format := strings.Repeat("%v ", len(v))
	tl.log.Print(ERROR, format, v...)
}
