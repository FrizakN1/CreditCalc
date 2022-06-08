package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var Logger *log.Logger

func init() {
	currentDate := time.Now().String()[0:10]

	logFile, e := os.OpenFile(currentDate+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if e != nil {
		fmt.Println(e)
		return
	}

	Logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
}
