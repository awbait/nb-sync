package log

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Trace   *log.Logger
	Error   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	wrt := io.MultiWriter(os.Stdout, file)

	Info = log.New(wrt, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(wrt, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Trace = log.New(wrt, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(wrt, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
