package log

import (
	"fmt"
	"strings"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var g_logger *logrus.Logger
var g_logfile string

func InitLog(file string) {
	switch runtime.GOOS {
	case "windows":
		sysdata_path := os.Getenv("ProgramData")
		g_logfile = fmt.Sprintf("%s/ggcal/%s", sysdata_path, file)
	case "linux":
		g_logfile = fmt.Sprintf("/var/log/ggcal/%s", file)
	}
}

type LogFormatter struct{}

func (f *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	msg := entry.Message
	msg = strings.TrimSuffix(msg, "\n")
	funcName := entry.Caller.Function
	line := entry.Caller.Line
	logLine := fmt.Sprintf("%s [%s]\t%s, %s:%d\n", timestamp, level, msg, funcName, line)
	return []byte(logLine), nil
}

func LogService() *logrus.Logger {
	if g_logger != nil {
		return g_logger
	}

	ljack := &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s", g_logfile),
		MaxSize:    3,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}

	writer := io.MultiWriter(os.Stdout, ljack)

	logpath := filepath.Dir(g_logfile)
	err := os.MkdirAll(logpath, os.ModeDir)
	if err != nil && !os.IsExist(err) {
		return nil
	}

	g_logger = logrus.New()
	g_logger.SetFormatter(&LogFormatter{})
	g_logger.SetOutput(writer)
	g_logger.SetReportCaller(true)
	g_logger.SetLevel(logrus.InfoLevel)

	return g_logger
}