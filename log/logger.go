package log

import (
	"fmt"
	_log "log"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	color "gopkg.in/gookit/color.v1"
)

var (
	colorDebug = color.New(color.Magenta)
	colorWarn  = color.New(color.Yellow)
	colorError = color.New(color.Red)
	colorInfo  = color.New()

	colorCaller = color.New()
	colorLine   = color.New(color.Bold)

	colorFieldKey   = color.New(color.Gray.Darken())
	colorFieldValue = color.New(color.Blue.Light().Light())
)

type Fields = map[string]interface{}
type F = Fields

type Context interface {
	WithValue(k string, v interface{}) Context
	WithValues(f Fields) Context
	Map() Fields
}

type contextImpl struct {
	parent Context
	m      Fields
}

var _ Context = (*contextImpl)(nil)

func emptyContextImpl() *contextImpl {
	ci := new(contextImpl)
	ci.m = make(Fields, 0)
	ci.parent = nil
	return ci
}

func (ci *contextImpl) WithValue(k string, v interface{}) Context {
	childCi := new(contextImpl)
	childCi.parent = ci
	m := make(Fields, 1)
	m[k] = v
	childCi.m = m
	return childCi
}

func (ci *contextImpl) WithValues(f Fields) Context {
	childCi := new(contextImpl)
	childCi.parent = ci
	childCi.m = f
	return childCi
}

func (ci *contextImpl) Map() Fields {
	if ci.parent != nil {
		pm := ci.parent.Map()
		nm := make(Fields, len(ci.m))
		for k, v := range pm {
			nm[k] = v
		}
		for k, v := range ci.m {
			nm[k] = v
		}
		return nm
	} else {
		nm := make(Fields, len(ci.m))
		for k, v := range ci.m {
			nm[k] = v
		}
		return nm
	}
}

type Level int32

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
)

type Logger struct {
	level  Level
	logger *_log.Logger
	fields Context
}

func New(level Level) *Logger {
	l := new(Logger)
	l.level = level
	l.logger = _log.Default()
	l.fields = emptyContextImpl()
	return l
}

func getCaller() string {
	if _, path, line, ok := runtime.Caller(3); ok {
		_, file := filepath.Split(path)
		return colorCaller.Sprint(file) + ":" + colorLine.Sprint(line)
	}
	return "<???>"
}

func getSeverity(level Level) string {
	severity := ""
	switch level {
	case LevelDebug:
		severity = colorDebug.Render("DEBUG")
	case LevelInfo:
		severity = colorInfo.Render("INFO")
	case LevelWarn:
		severity = colorWarn.Render("WARN")
	case LevelError:
		severity = colorError.Render("ERROR")
	case LevelPanic:
		severity = colorError.Render("FATAL")
	}
	return "[" + severity + "]"
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
	newLogger := new(Logger)
	newLogger.level = l.level
	newLogger.logger = l.logger
	newLogger.fields = l.fields.WithValue(key, value)
	return newLogger
}

func (l *Logger) WithFields(fields Fields) *Logger {
	newLogger := new(Logger)
	newLogger.level = l.level
	newLogger.logger = l.logger
	newLogger.fields = l.fields.WithValues(fields)
	return newLogger
}

func (l *Logger) structure(level Level, msgs ...interface{}) []interface{} {
	//timestamp := "[" + time.Now().Format(time.RFC3339) + "]"

	fieldMap := l.fields.Map()
	_fields := make([]string, 0, len(fieldMap))
	for k, v := range fieldMap {
		_fields = append(_fields, "\t"+colorFieldKey.Render(k)+"="+colorFieldValue.Render(v))
	}
	sort.Strings(_fields)
	//fields := strings.Join(_fields, ", ")
	var fields string
	if len(_fields) > 0 {
		fields = "\n" + strings.Join(_fields, "\n")
	}

	ctx := make([]interface{}, 3+len(msgs))
	//ctx[0] = timestamp
	ctx[0] = getSeverity(level)
	ctx[1] = getCaller() + ":"
	for i, v := range msgs {
		ctx[2+i] = v
	}
	ctx[len(ctx)-1] = fields
	return ctx
}

func (l *Logger) Debug(msgs ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.Println(l.structure(LevelDebug, msgs...)...)
	}
}

func (l *Logger) Debugf(f string, msgs ...interface{}) {
	if l.level <= LevelDebug {
		l.logger.Println(l.structure(LevelDebug, fmt.Sprintf(f, msgs...))...)
	}
}

func (l *Logger) Info(msgs ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.Println(l.structure(LevelInfo, msgs...)...)
	}
}

func (l *Logger) Infof(f string, msgs ...interface{}) {
	if l.level <= LevelInfo {
		l.logger.Println(l.structure(LevelInfo, fmt.Sprintf(f, msgs...))...)
	}
}

func (l *Logger) Warn(msgs ...interface{}) {
	if l.level <= LevelWarn {
		l.logger.Println(l.structure(LevelWarn, msgs...)...)
	}
}

func (l *Logger) Warnf(f string, msgs ...interface{}) {
	if l.level <= LevelWarn {
		l.logger.Println(l.structure(LevelWarn, fmt.Sprintf(f, msgs...))...)
	}
}

func (l *Logger) Error(msgs ...interface{}) {
	if l.level <= LevelError {
		l.logger.Println(l.structure(LevelError, msgs...)...)
	}
}

func (l *Logger) Errorf(f string, msgs ...interface{}) {
	if l.level <= LevelError {
		l.logger.Println(l.structure(LevelError, fmt.Sprintf(f, msgs...))...)
	}
}

func (l *Logger) Panic(msgs ...interface{}) {
	if l.level <= LevelPanic {
		l.logger.Panicln(l.structure(LevelPanic, msgs...)...)
	}
}

func (l *Logger) Panicf(f string, msgs ...interface{}) {
	if l.level <= LevelPanic {
		l.logger.Panicln(l.structure(LevelPanic, fmt.Sprintf(f, msgs...))...)
	}
}

var defaultLogger = New(LevelInfo)

func WithField(key string, value interface{}) *Logger {
	return defaultLogger.WithField(key, value)
}

func WithFields(fields Fields) *Logger {
	return defaultLogger.WithFields(fields)
}

func Debug(msgs ...interface{}) {
	defaultLogger.Debug(msgs...)
}

func Debugf(f string, msgs ...interface{}) {
	defaultLogger.Debugf(f, msgs...)
}

func Info(msgs ...interface{}) {
	defaultLogger.Info(msgs...)
}

func Infof(f string, msgs ...interface{}) {
	defaultLogger.Infof(f, msgs...)
}

func Warn(msgs ...interface{}) {
	defaultLogger.Warn(msgs...)
}

func Warnf(f string, msgs ...interface{}) {
	defaultLogger.Warnf(f, msgs...)
}

func Error(msgs ...interface{}) {
	defaultLogger.Error(msgs...)
}

func Errorf(f string, msgs ...interface{}) {
	defaultLogger.Errorf(f, msgs...)
}

func Panic(msgs ...interface{}) {
	defaultLogger.Panic(msgs...)
}

func Panicf(f string, msgs ...interface{}) {
	defaultLogger.Panicf(f, msgs...)
}

func SetLevel(level Level) {
	defaultLogger.level = level
}
