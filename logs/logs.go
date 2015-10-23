// logs provides various utility functions for logging.
package logs

import (
	"log"
	"os"
	"time"
	"fmt"
)

var logger *log.Logger

// Initialize instantiates the logger.
func Initialize(prefix string) {
	logger = log.New(os.Stdout, fmt.Sprintf("%s: ", prefix), log.Ldate|log.Ltime)
}

// Log logs anything that is passed in.
func Log(msg interface{}) {
	logger.Println(msg)
}

// CheckFatal checks for an error and terminates if one is found.
func CheckFatal(err error) {
	if err != nil {
		logger.Fatalln("ERROR:", err)
	}
}

// CheckErr logs the error and returnns true if there is an error
// otherwise it returns false.
func CheckErr(err error) bool {
	if err != nil {
		logger.Println(err)
		return true
	}
	return false
}

// TimerBegin records a time to be used with TimerEnd.
func TimerBegin(s string) (string, time.Time) {
	logger.Printf("STARTING %s", s)
	return s, time.Now()
}

// TimerEnd logs the time differene between now and the given startTime.
func TimerEnd(s string, startTime time.Time) {
	logger.Printf("FINISHED %s (%s)", s, time.Now().Sub(startTime))
}
