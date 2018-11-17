package logger

import (
	"fmt"
	"io"
	"log"
)

/* Get a configured logger */
func GetLogger(tag string, dest io.Writer) *log.Logger {
	var logger *log.Logger = log.New(dest,
		fmt.Sprintf("%v: ", tag),
		log.Ldate|log.Ltime|log.Lshortfile)
	return logger
}
