package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
}

type Field struct {
	Key   string
	Value interface{}
}

type loggerTmpl struct {
	level  Level
	logger *log.Logger
	fields []Field
}

type Options struct {
	Level  Level
	Output io.Writer
}

func New(options ...func(*Options)) Logger {
	opts := &Options{
		Level:  LevelInfo,
		Output: os.Stdout,
	}

	for _, option := range options {
		option(opts)
	}

	return &loggerTmpl{
		level:  opts.Level,
		logger: log.New(opts.Output, "", log.LstdFlags|log.Lmicroseconds),
		fields: make([]Field, 0),
	}
}

func WithLevel(level Level) func(*Options) {
	return func(o *Options) {
		o.Level = level
	}
}

func WithOutput(writer io.Writer) func(*Options) {
	return func(o *Options) {
		o.Output = writer
	}
}

func (l *loggerTmpl) Debug(msg string, fields ...Field) {
	if l.level <= LevelDebug {
		l.log("DEBUG", msg, fields)
	}
}

func (l *loggerTmpl) Info(msg string, fields ...Field) {
	if l.level <= LevelInfo {
		l.log("INFO", msg, fields)
	}
}

func (l *loggerTmpl) Warn(msg string, fields ...Field) {
	if l.level <= LevelWarn {
		l.log("WARN", msg, fields)
	}
}

func (l *loggerTmpl) Error(msg string, fields ...Field) {
	if l.level <= LevelError {
		l.log("ERROR", msg, fields)
	}
}

func (l *loggerTmpl) log(level, msg string, fields []Field) {
	allFields := append(l.fields, fields...)
	var builder strings.Builder
	for _, field := range allFields {
		builder.WriteString(l.formatField(field) + " ")
	}
	l.logger.Printf("[%s] %s %s", level, msg, builder.String())
}

func (l *loggerTmpl) formatField(field Field) string {
	return field.Key + "=" + l.formatValue(field.Value)
}

func (l *loggerTmpl) formatValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (l *loggerTmpl) With(fields ...Field) Logger {
	return &loggerTmpl{
		level:  l.level,
		logger: l.logger,
		fields: append(l.fields, fields...),
	}
}

func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Err(err error) Field {
	return Field{Key: "error", Value: err}
}
