package main

import (
	"fmt"
	"log"
	"math/rand"
)

type LogTrialler interface {
	LogTrial(int)
	waitForLog_trial(chan string, chan bool)
	trialFinished() bool
	createRandomLogs(chan string, int)
}

type logTrialler struct {
	logger Logger
}

func newLogTrialler(l Logger) (*logTrialler, error) {
	var trialler logTrialler

	trialler.logger = l
	return &trialler, nil
}

func NewLogTrialler(l Logger) (LogTrialler, error) {
	return newLogTrialler(l)
}

// makes n random log messages and logs them.
func (lt logTrialler) LogTrial(amountMessages int) {
	totalTrialMessages = amountMessages

	c := make(chan string)
	exit := make(chan bool)
	go lt.waitForLog_trial(c, exit)

	lt.createRandomLogs(c, totalTrialMessages)

	<-exit
	fmt.Println("Done")
}

// log listener for trial and testing purposes, difference is that it calls exit when the chosen amount of messages are logged
func (lt logTrialler) waitForLog_trial(c chan string, exit chan bool) {
	for {
		lt.logger.logIt(<-c)
		if lt.logger.logsFull() { // if logs have reached full capacity, exit
			lt.logger.rotate()
		}
		if lt.trialFinished() {
			exit <- true
			return
		}
	}
}

// checks if all the messages in the log trial have been logged
func (lt logTrialler) trialFinished() bool {
	rotateAmount := rotateCount * (lt.logger.GetMaxFiles() * lt.logger.GetMaxLinesPerFile())
	fmt.Println((lineCount + (fileCount * lt.logger.GetMaxLinesPerFile()) + rotateAmount), (lineCount+(fileCount*lt.logger.GetMaxLinesPerFile())+rotateAmount) == totalTrialMessages)
	return (lineCount + (fileCount * lt.logger.GetMaxLinesPerFile()) + rotateAmount) == totalTrialMessages
}

// used by LogTrial() to create random logs, this is mainly for testing
func (lt logTrialler) createRandomLogs(c chan string, amountMessages int) {
	if amountMessages <= 0 {
		log.Fatal("Must have at least one message.") // change
	}
	for i := 0; i < amountMessages; i++ {
		c <- createRandomLog()
	}
}

// rolls the dice and returns a random log based on the roll
func createRandomLog() string {
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
