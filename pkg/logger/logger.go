package logger

import (
	"fmt"
	sysLog "log"
	"os"
)

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
)

var loggers = make(map[string]*sysLog.Logger)

func getLogger(name string) *sysLog.Logger {
	if loggers[name] == nil {
		loggers[name] = sysLog.New(os.Stdout, "["+name+"] ", sysLog.Ldate|sysLog.Ltime)
	}

	return loggers[name]
}

// Info prints a normal message with the default color
func Info(format string, v ...any) {
	getLogger("info").Println(fmt.Sprintf(format, v...))
}

// Success prints a success message in the green color
func Success(format string, v ...any) {
	getLogger("info").Println(GREEN + fmt.Sprintf(format, v...) + RESET)
}

// Warning prints a warning message in the yellow color
func Warning(format string, v ...any) {
	getLogger("info").Println(YELLOW + fmt.Sprintf(format, v...) + RESET)
}

// Fatal prints a fatal message in the red color and then panic
func Fatal(format string, v ...any) {
	getLogger("fatal").Fatal(RED + fmt.Sprintf(format, v...) + RESET)
}
