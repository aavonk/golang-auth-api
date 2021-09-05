package logger

import (
	"log"
	"os"
)

var (
	Info  = log.New(os.Stdout, "[ AUTH-SERVICE :: INFO ] \t", log.Ldate|log.Ltime)
	Error = log.New(os.Stderr, "[ AUTH-SERVICE :: ERROR ] \t", log.Ldate|log.Ltime|log.Lshortfile)
)
