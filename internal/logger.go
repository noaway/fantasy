package internal

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"fantasy/iconfig"
	"github.com/sirupsen/logrus"
)

const (
	defaultFormat = "2006/01/02 15:04:05.000000000"
)

type ILog interface {
	SetPrefix(string)
	Debug(...interface{})
	Debugf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
}

type Logger struct {
	file       *os.File
	skip       int
	filename   string
	prefix     string
	timeFormat string
}

func (l *Logger) createLogFile() (err error) {
	l.file, err = os.OpenFile(l.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	return
}

func (l *Logger) SetPrefix(p string) {
	l.prefix = p
}

func (l *Logger) Debug(v ...interface{}) {
	l.width().Debug(v)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.width().Debugf(format, v)
}

func (l *Logger) Info(v ...interface{}) {
	l.width().Info(v)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.width().Infof(format, v)
}

func (l *Logger) Warn(v ...interface{}) {
	l.width().Warn(v)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.width().Warnf(format, v)
}

func (l *Logger) Error(v ...interface{}) {
	l.width().Error(v)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.width().Errorf(format, v)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.width().Fatal(v)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.width().Fatalf(format, v)
}

func (l *Logger) width() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{"": l.prefix, "| ": l.caller(3)})
}

func (l *Logger) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestampFormat := l.timeFormat
	if timestampFormat == "" {
		timestampFormat = logrus.DefaultTimestampFormat
	}

	l.appendKeyValue(b, "", entry.Time.Format(timestampFormat))

	for k := range entry.Data {
		l.appendKeyValue(b, k, entry.Data[k])
	}

	l.appendKeyValue(b, "", entry.Level.String())
	l.appendKeyValue(b, "", entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (l *Logger) short(file string) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return short
}

func (l *Logger) caller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "???"
		line = 0
	}
	f := runtime.FuncForPC(pc)
	return fmt.Sprintf("%v:%v func:%v", l.short(file), line, f.Name())
}

func (l *Logger) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	b.WriteString(key)
	fmt.Fprint(b, value)
	b.WriteByte(' ')
}

func NewLogger(c iconfig.IConfig, log_name string) *Logger {
	var (
		path string
		err  error
	)

	if c != nil {
		path, err = c.String("log", "path")
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(path) == 0 {
		path = fmt.Sprintf("./%s", log_name)
	}

	l := new(Logger)
	l.filename = path
	l.timeFormat = defaultFormat
	logrus.SetFormatter(l)
	logrus.SetLevel(logrus.InfoLevel)
	if err := l.createLogFile(); err != nil {
		panic(err)
	}
	logrus.SetOutput(l.file)
	return l
}
