package utils

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Logging struct log rotate writer. can be use as io.Writer interface
type Logging struct {
	filename  string // should be set to the actual filename
	fp        *os.File
	logger    *log.Logger
	day       int
	logWorker chan Logs
	quitChan  chan bool
}

var (
	once   sync.Once
	LogObj *Logging
)

func SetLoggerObj(file string) *Logging {
	once.Do(func() {
		// create a new instance

		LogObj = &Logging{
			filename: file,
		}

		LogObj.init()
	})

	return LogObj

}

//InitializeLogger Logging
func (l *Logging) init() {

	logrus.Info("Initializing logger")

	var err error
	_, _, l.day = time.Now().Date()
	l.fp, err = os.OpenFile("Logs/"+time.Now().Format("02-Jan-2006_")+l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatalln("failed to set log file, err: " + err.Error())
	}

	l.logger = log.New(l.fp, "", log.Ldate|log.Lmicroseconds)
}

func (l *Logging) WriteLogs(output Logs) {
	go func() {
		WriteInDB(output)
	}()
}

//Println Logging
func (l *Logging) Info(output ...interface{}) {
	_, _, day := time.Now().Date()
	if l.day != day {
		_ = l.fp.Close()
		l.init()
	}
	go func(l *Logging) {
		_, _, day := time.Now().Date()
		if l.day != day {
			log.Printf(" change in day")
			l.fp.Close()
			l.init()
		}

		l.logger.Println("[INFO]: ", output)
	}(l)
}

func (l *Logging) Error(output ...interface{}) {
	_, _, day := time.Now().Date()
	if l.day != day {
		_ = l.fp.Close()
		l.init()
	}
	go func(l *Logging) {
		_, _, day := time.Now().Date()
		if l.day != day {
			log.Printf(" change in day")
			l.fp.Close()
			l.init()
		}

		l.logger.Println("[ERROR]: ", output)
	}(l)
}

func (l *Logging) Fatal(output ...interface{}) {
	_, _, day := time.Now().Date()
	if l.day != day {
		_ = l.fp.Close()
		l.init()
	}
	go func(l *Logging) {
		_, _, day := time.Now().Date()
		if l.day != day {
			log.Printf(" change in day")
			l.fp.Close()
			l.init()
		}

		l.logger.Fatal("[FATAL]: ", output)
		l.fp.Close()
	}(l)
}
