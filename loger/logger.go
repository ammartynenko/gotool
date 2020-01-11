package loger

import (
	"errors"
	"io"
	"log"
	"log/syslog"
	"os"
	"strings"
)

type ToolLogger struct {
	fileLog *os.File
	sysLog  *syslog.Writer
	lg      *log.Logger
}
type ToolLoggerConfig struct {
	FileLog  string
	SyslogIO *SyslogIOConfig
	LogCfg   *LoggerConfig
}
type SyslogIOConfig struct {
	Tag          string
	Priority     syslog.Priority
	NetConnect   *NetworkedConfig
	LocalConnect bool //имеет приоритет над NetConnect
}
type NetworkedConfig struct {
	Network    string //tcp or udp
	RemoteAddr string //remote address `localhost:514`
}
type LoggerConfig struct {
	Prefix       string
	LoggerBitmap int
	Out          io.Writer
}

var (
	errorLogConfig = errors.New("LoggerConfig is nil. Abort.")
)

const (
	prefixInfo    = "[INFO]"
	prefixWarning = "[WARNING]"
	prefixError   = "[ERROR]"
)

func NewToolLogger(cfg ToolLoggerConfig) (*ToolLogger, error) {
	t := new(ToolLogger)
	//file logging section
	if strings.TrimSpace(cfg.FileLog) != "" {
		f, err := os.OpenFile(cfg.FileLog, os.O_APPEND|os.O_RDWR, 0640)
		switch {
		case errors.Is(err, os.ErrPermission):
			return nil, os.ErrPermission
		case errors.Is(err, os.ErrNotExist):
			f, err = os.Create(cfg.FileLog)
			if err != nil {
				return nil, err
			}
		case errors.Is(err, os.ErrInvalid):
			return nil, os.ErrInvalid
		case errors.Is(err, os.ErrExist):
			return nil, os.ErrExist
		}
		t.fileLog = f
	}
	//syslog section
	if cfg.SyslogIO != nil {
		//создаю вритер локальный если он есть
		if cfg.SyslogIO.LocalConnect && cfg.SyslogIO.NetConnect == nil {
			sl, err := syslog.New(cfg.SyslogIO.Priority, cfg.SyslogIO.Tag)
			if err != nil {
				return t, err
			}
			t.sysLog = sl
		}
		//если локальный выключен использую сетевой если есть его конфигурация
		if cfg.SyslogIO.LocalConnect == false && cfg.SyslogIO.NetConnect != nil {
			sl, err := syslog.Dial(cfg.SyslogIO.NetConnect.Network, cfg.SyslogIO.NetConnect.RemoteAddr,
				cfg.SyslogIO.Priority, cfg.SyslogIO.Tag)
			if err != nil {
				return t, err
			}
			t.sysLog = sl
		}
	}
	//обычный os.stdout
	if cfg.LogCfg == nil {
		return t, errorLogConfig
	}
	t.lg = log.New(cfg.LogCfg.Out, cfg.LogCfg.Prefix, cfg.LogCfg.LoggerBitmap)
	return t, nil
}
func (t *ToolLogger) puts(prefix, msg string) {
	concat := strings.Join([]string{prefix, msg}, " ")
	if t.lg != nil {
		t.lg.Print(concat)
	}
	if t.sysLog != nil {
		_ = t.sysLog.Info(concat)
	}
	if t.fileLog != nil {
		_, _ = t.fileLog.WriteString(concat)
	}
}

func (t *ToolLogger) Info(msg string) {
	t.puts(prefixInfo, msg)
}
func (t *ToolLogger) Warning(msg string) {
	t.puts(prefixWarning, msg)

}
func (t *ToolLogger) Error(msg string) {
	t.puts(prefixError, msg)
}
