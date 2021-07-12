package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// global variables to keep track of file and log capacity
var totalMessages int
var lineCount = 0
var fileCount = 0

type Logger struct {
	MAX_FILES          int    // max amount of files that should be filled with logs
	MAX_LINES_PER_FILE int    // max amount of logs in a single file
	FILE_PREFIX        string // prefix for log output files
}

// creates a new Logger struct and returns it
func NewLogger(MAX_FILES int, MAX_LINES_PER_FILE int, FILE_PREFIX string) Logger {
	logger := Logger{
		MAX_FILES:          MAX_FILES,
		MAX_LINES_PER_FILE: MAX_LINES_PER_FILE,
		FILE_PREFIX:        FILE_PREFIX,
	}

	return logger
}

// starts listening for logs, returning the logging channel and the exit channel
func (l Logger) StartLogListener() (chan string, chan bool) {
	c := make(chan string)
	exit := make(chan bool)

	totalMessages = l.MAX_FILES * l.MAX_LINES_PER_FILE

	fmt.Println("Listening for logs.")
	go l.waitForLog(c, exit)

	return c, exit
}

// call this to stop listening for logs
func (l Logger) CallExit(exit chan bool) {
	exit <- true
	fmt.Println("Exited.")
}

// adds a message to the logger to be logged
func (l Logger) AddMsg(c chan string, msg string) {
	c <- msg
}

// used by LogTrial() to create random logs, this is mainly for testing
func (l Logger) createRandomLogs(c chan string, amount_messages int) {
	if amount_messages <= 0 {
		log.Fatal("Must have at least one message.")
	} else if amount_messages > l.MAX_FILES*l.MAX_LINES_PER_FILE {
		totalMessages = l.MAX_FILES * l.MAX_LINES_PER_FILE
	} else {
		totalMessages = amount_messages
	}

	for i := 0; i < totalMessages; i++ {
		c <- l.createRandomLog()
	}
}

// rolls the dice and returns a random log based on the roll
func (l Logger) createRandomLog() string {
	var msg string

	switch msg_num := rand.Intn(3); msg_num {
	case 0:
		msg = "THING HAS HAPPENED"
	case 1:
		msg = "STUFF HAS HAPPENED"
	case 2:
		msg = "EVERYTHING IS BROKEN"
	}

	return msg
}

// active listener that waits for an input to either logging channel or exit channel, and responds accordingly
func (l Logger) waitForLog(c chan string, exit chan bool) {
	for {
		select {
		case msg := <-c:
			if l.logsFull() { // if logs have reached full capacity, exit
				l.logIt(msg)
				exit <- true

			} else {
				l.logIt(msg)
			}
		case <-exit:
			return
		}
	}
}

// selects the correct file to output to, then adds the time information to the log string and appends it to output file
func (l Logger) logIt(msg string) {
	l.checkForFullFile()

	fileName := l.FILE_PREFIX + "_" + strconv.Itoa(fileCount) + ".txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log_msg := "[" + time.Time.String(time.Now()) + "] " + msg + "\n"
	file.WriteString(log_msg)

	lineCount += 1
}

// if file is full then increment the file count to start using a different, empty file
func (l Logger) checkForFullFile() {
	if lineCount == l.MAX_LINES_PER_FILE {
		lineCount = 0
		fileCount += 1
	}
}

// checks if logs have reached full capacity
func (l Logger) logsFull() bool {
	return ((lineCount + 1) + (fileCount * 8)) >= totalMessages
}

// makes n random log messages and logs them.
func (l Logger) LogTrial(totalMessages int) {
	c := make(chan string)
	exit := make(chan bool)
	go l.waitForLog(c, exit)

	l.createRandomLogs(c, totalMessages)

	<-exit
	fmt.Println("Done")
}
