package logger

import (
	"be-yourmoments/upload-svc/internal/helper/utils"
	"log"
	"os"
)

var logLevel = utils.GetEnv("LOG_LEVEL")

type Log struct {
	stdout *log.Logger
	stderr *log.Logger
}

type Options struct {
	IsPrintStack bool
	IsExit       bool
	ExitCode     int
}

func New(prefix string) *Log {

	return &Log{
		stdout: log.New(os.Stdout, "[LOG]["+prefix+"]", log.Ldate|log.Ltime),
		stderr: log.New(os.Stderr, "[ERROR]["+prefix+"]", log.Ldate|log.Ltime),
	}
}
