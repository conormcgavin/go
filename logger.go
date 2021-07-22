package logger

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const DEFAULT_MAX_FILES = 8
const DEFAULT_MAX_LINES = 512

// global variables to keep track of file and log capacity
var lineCount = 0
var fileCount = 0
var rotateCount = 0

type Logger interface {
	StartLogListener() (chan string, chan bool)
	CallExit(chan bool)
	AddMsg(chan string, ...interface{})

	GetMaxFiles() int
	GetMaxLinesPerFile() int
	GetFilePrefix() string

	waitForLog(chan string, chan bool)
	checkForFullFile()
	logIt(string)
	rotate()
	logsFull() bool
}

type logger struct {
	MAX_FILES          int    // max amount of files that should be filled with logs
	MAX_LINES_PER_FILE int    // max amount of logs in a single file
	FILE_PREFIX        string // prefix for log output files
}

// creates a new logger struct and returns it. Set to 0 or lower to work with default config
func newLogger(MAX_FILES int, MAX_LINES_PER_FILE int, FILE_PREFIX string) (*logger, error) {
	if MAX_FILES < 0 || MAX_LINES_PER_FILE < 0 {
		return nil, fmt.Errorf("amount and capacity of files must be over 0")
	}
	if MAX_FILES == 0 {
		MAX_FILES = DEFAULT_MAX_FILES
	}
	if MAX_LINES_PER_FILE == 0 {
		MAX_LINES_PER_FILE = DEFAULT_MAX_LINES
	}

	var l logger

	l.MAX_FILES = MAX_FILES
	l.MAX_LINES_PER_FILE = MAX_LINES_PER_FILE
	l.FILE_PREFIX = FILE_PREFIX

	return &l, nil
}

// creates a new logger struct and returns it. Set to 0 or lower to work with default config
func NewLogger(MAX_FILES int, MAX_LINES_PER_FILE int, FILE_PREFIX string) (Logger, error) {
	return newLogger(MAX_FILES, MAX_LINES_PER_FILE, FILE_PREFIX)
}

// starts listening for logs, returning the logging channel and the exit channel
func (l logger) StartLogListener() (chan string, chan bool) {
	c := make(chan string)
	exit := make(chan bool)

	go l.waitForLog(c, exit)

	return c, exit
}

// call this to stop listening for logs
func (l logger) CallExit(exit chan bool) {
	exit <- true
	fmt.Println("Exit Invoked.")
	reset()
}

// adds a message to the logger to be logged
func (l logger) AddMsg(c chan string, messages ...interface{}) { // change channel c to accept byte string!!!! and change
	for _, arg := range messages {
		c <- fmt.Sprintf("%v", arg)
	}
}

// active listener that waits for an input to either logging channel or exit channel, and responds accordingly
func (l logger) waitForLog(c chan string, exit chan bool) {
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
func (l logger) logIt(msg string) {
	l.checkForFullFile()
	fileName := l.FILE_PREFIX + "_" + strconv.Itoa(fileCount) + ".txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err.Error())
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
func (l logger) checkForFullFile() {
	if lineCount == l.MAX_LINES_PER_FILE {
		lineCount = 0
		fileCount += 1
	}
}

// checks if logs have reached full capacity
func (l logger) logsFull() bool {
	return (lineCount + (fileCount * l.MAX_LINES_PER_FILE)) >= l.MAX_FILES*l.MAX_LINES_PER_FILE
}

// rotates the log files to use old ones when the current one is full
func (l logger) rotate() {
	fmt.Println("Rotating")
	fileCount = 0
	lineCount = 0
	rotateCount += 1
}

func reset() {
	lineCount = 0
	fileCount = 0
	rotateCount = 0
}

func (l *logger) GetMaxFiles() int {
	return l.MAX_FILES
}

func (l *logger) GetMaxLinesPerFile() int {
	return l.MAX_LINES_PER_FILE
}

func (l *logger) GetFilePrefix() string {
	return l.FILE_PREFIX
}
