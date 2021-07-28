package logging

import (
	"fmt"
	"log"
	"os"
)

var (
	infoLog  = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
	errorLog = log.New(os.Stderr, "[ERROR] ", log.LstdFlags)
)

func Info(format string, a ...interface{}) {
	infoLog.Println(fmt.Sprintf(format, a...))
}

func Error(format string, a ...interface{}) {
	errorLog.Println(fmt.Sprintf(format, a...))
}
