package logger

import (
	"fmt"
	"io"
)

type Level int

const (
    DebugLevel Level = iota
    InfoLevel
    WarningLevel
    ErrorLevel
    FatalLevel
)

type Logger struct {
    writer io.Writer
    level  Level
    prefix string
}

func New(writer io.Writer, level Level, prefix string) *Logger {
    return &Logger{
        writer: writer,
        level:  level,
        prefix: prefix,
    }
}

func (l *Logger) Debug(f string, a ...interface{}) {
    if l.level <= DebugLevel {
        fmt.Fprintf(l.writer, "%s [DEBUG] %s\n", l.prefix, fmt.Sprintf(f, a...))
    }
}

func (l *Logger) Info(f string, a ...interface{}) {
    if l.level <= InfoLevel {
        fmt.Fprintf(l.writer, "%s [INFO] %s\n", l.prefix, fmt.Sprintf(f, a...))
    }
}

func (l *Logger) Warning(f string, a ...interface{}) {
		if l.level <= WarningLevel {
				fmt.Fprintf(l.writer, "%s [WARNING] %s\n", l.prefix, fmt.Sprintf(f, a...))
		}
}

func (l *Logger) Error(f string, a ...interface{}) {
		if l.level <= ErrorLevel {
				fmt.Fprintf(l.writer, "%s [ERROR] %s\n", l.prefix, fmt.Sprintf(f, a...))
		}
}

func (l *Logger) Fatal(f string, a ...interface{}) {
		if l.level <= FatalLevel {
				fmt.Fprintf(l.writer, "%s [FATAL] %s\n", l.prefix, fmt.Sprintf(f, a...))
		}
}