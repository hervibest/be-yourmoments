package logger

import (
	"fmt"
	"os"
	"runtime/debug"
)

func (l *Log) Debug(message interface{}, options ...*Options) {
	if logLevel == "DEBUG" {
		msg := []interface{}{message}

		if len(options) == 0 || options[0].IsPrintStack {
			msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
		}
		l.stderr.Println(msg...)

		if len(options) > 0 && options[0].IsExit {
			exitCode := 1
			if options[0].ExitCode > 1 {
				exitCode = options[0].ExitCode
			}

			os.Exit(exitCode)
		}
	}
}

func (l *Log) CustomDebug(title string, message interface{}, options ...*Options) {
	if logLevel == "DEBUG" {
		msg := []interface{}{}
		msg = append(msg, title)
		msg = append(msg, message)

		if len(options) == 0 || options[0].IsPrintStack {
			msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
		}

		l.stderr.Println(msg...)

		if len(options) > 0 && options[0].IsExit {
			exitCode := 1
			if options[0].ExitCode > 1 {
				exitCode = options[0].ExitCode
			}

			os.Exit(exitCode)
		}
	}
}
