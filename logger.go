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
var lineCount = 0
var fileCount = 0
var rotateCount = 0
var totalTrialMessages int

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

	go l.waitForLog(c, exit)

	return c, exit
}

// call this to stop listening for logs
func (l Logger) CallExit(exit chan bool) {
	exit <- true
	fmt.Println("Exit Invoked.")
}

// adds a message to the logger to be logged
func (l Logger) AddMsg(c chan string, msg string) {
	c <- msg
}

// active listener that waits for an input to either logging channel or exit channel, and responds accordingly
func (l Logger) waitForLog(c chan string, exit chan bool) {
	fmt.Println("Listening for logs.")
	for {
		select {
		case msg := <-c:
			l.logIt(msg)
			if l.logsFull() { // if logs have reached full capacity, exit
				l.rotate()
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

	fi, _ := file.Stat()
	if lineCount == 0 && fi.Size() > 0 { // overwrite file contents if we're starting a new file and the file already has contents.
		file.Truncate(0)
		file.Seek(0, 0)
	}

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
	return (lineCount + (fileCount * l.MAX_LINES_PER_FILE)) >= l.MAX_FILES*l.MAX_LINES_PER_FILE
}

// rotates the log files to use old ones when the current one is full
func (l Logger) rotate() {
	fmt.Println("Rotating")
	fileCount = 0
	lineCount = 0
	rotateCount += 1
}




// makes n random log messages and logs them.
func (l Logger) LogTrial(amountMessages int) {
	totalTrialMessages = amountMessages

	c := make(chan string)
	exit := make(chan bool)
	go l.waitForLog_trial(c, exit)

	l.createRandomLogs(c, totalTrialMessages)

	<-exit
	fmt.Println("Done")
}

// log listener for trial and testing purposes, difference is that it calls exit when the chosen amount of messages are logged
func (l Logger) waitForLog_trial(c chan string, exit chan bool) {
	for {
		l.logIt(<-c)
		if l.logsFull() { // if logs have reached full capacity, exit
			l.rotate()
		}
		if l.trialFinished() {
			exit <- true
			return
		}
	}
}

// checks if all the messages in the log trial have been logged
func (l Logger) trialFinished() bool {
	rotateAmount := rotateCount * (l.MAX_FILES * l.MAX_LINES_PER_FILE)
	fmt.Println((lineCount + (fileCount * l.MAX_LINES_PER_FILE) + rotateAmount), (lineCount+(fileCount*l.MAX_LINES_PER_FILE)+rotateAmount) == totalTrialMessages)
	return (lineCount + (fileCount * l.MAX_LINES_PER_FILE) + rotateAmount) == totalTrialMessages
}

// used by LogTrial() to create random logs, this is mainly for testing
func (l Logger) createRandomLogs(c chan string, amountMessages int) {
	if amountMessages <= 0 {
		log.Fatal("Must have at least one message.") // change
	}
	for i := 0; i < amountMessages; i++ {
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