package utils

import (
	"log"
	"os"
	"sync"

	"space/constants"
	"space/helpers"
)

var PrintlnLog = func(logType int, output ...interface{}) {
	logObj.Println(logType, output)
}

const LOG_FILE_NAME = "space"

// Logging struct log rotate writer. can be use as io.Writer interface
type Logging struct {
	filename    string // should be set to the actual filename
	fp          *os.File
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
	day         int
}

var once sync.Once
var logObj *Logging

func GetLoggerObj(file string) *Logging {
	once.Do(func() {
		// create a new instance
		logObj = &Logging{
			filename: file,
		}
		logObj.init()
	})

	return logObj

}

// InitializeLogger Logging
func (l *Logging) init() {

	var err error
	_, _, l.day = helpers.GetCurrentTimeInIST().Date()
	l.fp, err = os.OpenFile("logs/"+helpers.GetCurrentTimeInIST().Format("02-Jan-2006_")+l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	l.infoLogger = log.New(l.fp, "INFO:", log.Ldate|log.Lmicroseconds)
	l.warnLogger = log.New(l.fp, "WARN:", log.Ldate|log.Lmicroseconds)
	l.errorLogger = log.New(l.fp, "ERROR:", log.Ldate|log.Lmicroseconds)

}

// Println Logging
func (l *Logging) Println(logType int, output ...interface{}) {
	_, _, day := helpers.GetCurrentTimeInIST().Date()
	if l.day != day {
		_ = l.fp.Close()
		l.init()
	}
	go func(l *Logging, logType int) {
		_, _, day := helpers.GetCurrentTimeInIST().Date()
		if l.day != day {
			log.Printf(" change in day")
			l.fp.Close()
			l.init()
		}

		switch logType {
		case constants.ERROR:
			l.errorLogger.Println(output)
		case constants.INFO:
			l.infoLogger.Println(output)
		case constants.WARNING:
			l.warnLogger.Println(output)
		}

	}(l, logType)

}
