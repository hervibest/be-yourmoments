package logger

func (l *Log) Panic(message interface{}, options ...Options) {
	l.stderr.Panicln(message)
}

func (l *Log) CustomPanic(title string, message interface{}, options ...Options) {
	msg := []interface{}{}
	msg = append(msg, title)
	msg = append(msg, message)

	l.stderr.Panicln(msg...)
}
