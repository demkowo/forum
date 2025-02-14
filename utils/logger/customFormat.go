package logger

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *log.Entry) ([]byte, error) {

	time := entry.Time.Format("MST 2006-01-02 15:04:05.9999999")
	lvl := strings.ToUpper(fmt.Sprint(entry.Level))
	function := filepath.Base(entry.Caller.Function)
	file := entry.Caller.File
	line := entry.Caller.Line
	msg := entry.Message

	parts := strings.Split(file, "/")
	if len(parts) >= 4 {
		file = strings.Join(parts[len(parts)-3:], "/")
	}

	str := fmt.Sprintf("%s:%d", file, line)

	if msg != "" {
		return []byte(fmt.Sprintf("%-31s   [%-7s]   %-55s   %-50s\n    === %s\n\n", time, lvl, function, str, msg)), nil
	}
	return []byte(fmt.Sprintf("%-31s   [%-7s]   %-55s   %-50s\n", time, lvl, function, str)), nil
}
