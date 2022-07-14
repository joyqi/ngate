package log

import (
	"fmt"
	sysLog "log"
	"os"
)

type Type uint32

const (
	InfoLog Type = iota
	SuccessLog
	WarningLog
	FatalLog
)

const (
	RESET  = "\033[0m"
	RED    = "\033[31m"
	GREEN  = "\033[32m"
	YELLOW = "\033[33m"
)

var loggers = make(map[Type]*sysLog.Logger)
var DebugMode = false

func getLogger(name Type) *sysLog.Logger {
	if loggers[name] == nil {
		loggers[name] = sysLog.New(os.Stdout, fmt.Sprintf("[%d] ", name), sysLog.Ldate|sysLog.Ltime)
	}

	return loggers[name]
}

// Info prints a normal message with the default color
func Info(format string, v ...any) {
	getLogger(InfoLog).Println(fmt.Sprintf(format, v...))
}

// Debug prints a normal message with the default color
func Debug(format string, v ...any) {
	if DebugMode {
		Info(format, v...)
	}
}

// Success prints a success message in the green color
func Success(format string, v ...any) {
	getLogger(SuccessLog).Println(GREEN + fmt.Sprintf(format, v...) + RESET)
}

// Warning prints a warning message in the yellow color
func Warning(format string, v ...any) {
	getLogger(WarningLog).Println(YELLOW + fmt.Sprintf(format, v...) + RESET)
}

// Error prints a fatal message in the red color and then panic
func Error(format string, v ...any) {
	getLogger(FatalLog).Println(RED + fmt.Sprintf(format, v...) + RESET)
}

// Fatal prints a fatal message in the red color and then panic
func Fatal(format string, v ...any) {
	getLogger(FatalLog).Fatal(RED + fmt.Sprintf(format, v...) + RESET)
}
