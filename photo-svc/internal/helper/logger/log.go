package logger

import (
	"fmt"
	"runtime/debug"
)

func (l *Log) Log(message interface{}, options ...*Options) {
	msg := []interface{}{message}

	if len(options) > 0 && options[0].IsPrintStack {
		msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
	}

	l.stdout.Println(msg...)
}

func (l *Log) CustomLog(title string, message interface{}, options ...*Options) {
	msg := []interface{}{}
	msg = append(msg, title)
	msg = append(msg, message)

	if len(options) > 0 && options[0].IsPrintStack {
		msg = append(msg, fmt.Sprintf("\n%s", debug.Stack()))
	}

	l.stdout.Println(msg...)
}
