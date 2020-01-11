package loger

import (
	"log"
	"log/syslog"
	"os"
	"testing"
)

func TestNewToolLogger(t *testing.T) {
	cfg := ToolLoggerConfig{
		FileLog: "/tmp/testgotoollogger.log",
		SyslogIO: &SyslogIOConfig{
			Tag:          "testgotoolog",
			Priority:     syslog.LOG_WARNING | syslog.LOG_INFO,
			NetConnect:   nil,
			LocalConnect: true,
		},
		LogCfg: &LoggerConfig{
			Prefix:       "[TESTGOTOOLLOGGING]",
			LoggerBitmap: log.Ltime | log.Ldate | log.Lshortfile,
			Out:          os.Stdout,
		},
	}
	_, err := NewToolLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}
}
func TestToolLogger_Info(t *testing.T) {
	cfg := ToolLoggerConfig{
		FileLog: "/tmp/testgotoollogger.log",
		SyslogIO: &SyslogIOConfig{
			Tag:          "testgotoolog",
			Priority:     syslog.LOG_WARNING | syslog.LOG_INFO,
			NetConnect:   nil,
			LocalConnect: true,
		},
		LogCfg: &LoggerConfig{
			Prefix:       "[TESTGOTOOLLOGGING]",
			LoggerBitmap: log.Ltime | log.Ldate | log.Lshortfile,
			Out:          os.Stdout,
		},
	}
	tl, err := NewToolLogger(cfg)
	if err != nil {
		t.Fatal(err)
	}
	//info testing
	tl.Info(" SIMPLE MESSAGE TEST\n")
	tl.Warning(" SIMPLE MESSAGE TEST\n")
	tl.Error(" SIMPLE MESSAGE TEST\n")
}
