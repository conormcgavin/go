package logger

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"
)

const testFilePrefix = "test_log_output"

// test to make sure constructor works successfully
func Test_ConstructorSuccess(t *testing.T) {
	if l, err := newLogger(10, 10, "logs"); err != nil {
		t.Errorf("Constructor should have succeeded")
	} else if l.MAX_FILES != 10 {
		t.Errorf("Max Files Mismatch")
	} else if l.MAX_LINES_PER_FILE != 10 {
		t.Errorf("LPF Mismatch")
	} else if l.FILE_PREFIX != "logs" {
		t.Errorf("Prefix Mismatch")
	}
}

// test to make sure constructor fails with incorrect input
func Test_ConstructorFailure(t *testing.T) {

	if _, err := newLogger(10, -1, "logs"); err == nil {
		t.Errorf("Constructor should have failed")
	} else if err.Error() != "amount and capacity of files must be over 0" {
		t.Errorf("Unexpected error message")
	}
}

// test to make sure constructor implements default values when values initalised to 0
func Test_ConstructorDefaultsSuccess(t *testing.T) {
	if l, err := newLogger(0, 0, "logs"); err != nil {
		t.Errorf("Constructor should have succeeded")
	} else if l.MAX_FILES != DEFAULT_MAX_FILES || l.MAX_LINES_PER_FILE != DEFAULT_MAX_LINES {
		t.Errorf("Constructor should set these to the defult value, not 0")
	}
}

// test for getter methods to make sure they return correct values
func Test_InterfacesSuccess(t *testing.T) {
	R, _ := newLogger(2, 2, "test_log_output")

	if R.GetMaxFiles() != 2 {
		t.Errorf("Max Files is wrong")
	} else if R.GetMaxLinesPerFile() != 2 {
		t.Errorf("Max Lines per File is wrong")
	} else if R.GetFilePrefix() != "test_log_output" {
		t.Errorf("Prefix is wrong")
	}
}

// test to assure basic string is correctly logged
func Test_LoggingMessagesBasicSuccess(t *testing.T) {
	L, _ := newLogger(2, 2, testFilePrefix)
	c, exit := L.StartLogListener()
	defer L.CallExit(exit)
	defer cleanUpTests(*L)

	L.AddMsg(c, "Test")
	L.AddMsg(c, "Test")

	stringList := []string{
		"Test",
		"Test",
	}

	err := assertCorrectLogging(stringList)
	if err != nil {
		t.Errorf(err.Error())
	}

	cleanUpTests(*L)
}

// test to assure basic string is not wronfully classed as correctly logged
func Test_LoggingMessagesBasicFailure(t *testing.T) {
	L, _ := newLogger(2, 2, testFilePrefix)
	c, exit := L.StartLogListener()
	defer L.CallExit(exit)
	defer cleanUpTests(*L)

	L.AddMsg(c, "Test")
	L.AddMsg(c, "Test")

	stringList := []string{
		"BadTest",
		"BadTest",
	}

	err := assertCorrectLogging(stringList)
	if err == nil {
		t.Errorf("should return an error due to mismatched logging")
	}

	cleanUpTests(*L)
}

// test to assure variadic inputs of different types are all correctly turned into strings and logged
func Test_LoggingMessagesTypesSuccess(t *testing.T) {
	L, _ := newLogger(1, 100, testFilePrefix)
	c, exit := L.StartLogListener()
	defer L.CallExit(exit)
	defer cleanUpTests(*L)

	L.AddMsg(c, "Test", "Test Test")
	L.AddMsg(c, 1, 2, 3)
	L.AddMsg(c, "")
	L.AddMsg(c, "", true, false, 100, 432.3)
	L.AddMsg(c, []byte("Hello"))
	L.AddMsg(c, map[string]string{"Hi": "Bye", "Bonjour": "Au Revoir"})
	L.AddMsg(c, L)
	L.AddMsg(c, c)
	L.AddMsg(c, *L)

	stringList := []string{
		"Test",
		"Test Test",
		"1",
		"2",
		"3",
		"",
		"",
		"true",
		"false",
		"100",
		"432.3",
		fmt.Sprintf("%v", []byte("Hello")),
		fmt.Sprintf("%v", map[string]string{"Hi": "Bye", "Bonjour": "Au Revoir"}),
		fmt.Sprintf("%v", L),
		fmt.Sprintf("%v", c),
		fmt.Sprintf("%v", *L),
	}
	err := assertCorrectLogging(stringList)
	if err != nil {
		t.Errorf(err.Error())
	}
}

// helper function that takes in a file and a slice of messages and compares the lines (logs) in the file to the strings one by one. Assures they were logged correctly.
func assertCorrectLogging(logMsg []string) error {
	fileName := testFilePrefix + "_0.txt"
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("Unable to read file contents")
	}
	string_content := string(content)
	lines := strings.Split(string_content, "\n")
	for i := 0; i < len(lines)-1; i++ {

		stuff := lines[i][len(lines[i])-(len(logMsg[i])):]
		if stuff != logMsg[i] {
			return fmt.Errorf("strings not being logged correctly")
		}
	}

	return nil
}

// test to check that the logger can detect when a file is "full"
func Test_checkForFullFile(t *testing.T) {
	L, _ := newLogger(2, 2, testFilePrefix)
	fileCount = 0
	lineCount = 2
	L.checkForFullFile()
	if fileCount != 1 || lineCount != 0 {
		t.Errorf("Error dealing with full files")
	}

	fileCount = 0
	lineCount = 1
	L.checkForFullFile()
	if fileCount != 0 || lineCount != 1 {
		t.Errorf("Classing non-full files as full.")
	}

	reset()
}

// test to check that the logger can detect when the logs are "full" and hence need to be rotated
func Test_logsFull(t *testing.T) {
	L, _ := newLogger(2, 2, testFilePrefix)
	fileCount = 0
	lineCount = 2
	if L.logsFull() {
		t.Errorf("Error dealing with full files")
	}

	fileCount = 0
	lineCount = 1
	L.checkForFullFile()
	if fileCount != 0 || lineCount != 1 {
		t.Errorf("Classing non-full files as full.")
	}
}

// test that the logger will still log messages correctly when it rotates and that it accounts for rotations correctly.
func Test_correctRotation(t *testing.T) {
	L, _ := newLogger(2, 8, testFilePrefix)
	LT, _ := newLogTrialler(L)
	LT.LogTrial(20, true)
	filename := testFilePrefix + "_" + strconv.Itoa(fileCount) + ".txt"
	if rotateCount != 1 || lineCount != 4 || fileCount != 0 { // test it is performing rotations correctly.
		t.Errorf("not handling rotations correctly.")
	} else if _, err := os.Open(filename); err != nil { // test that it is creating new files instead of repeating the same one
		t.Errorf("Not creating files correctly.")
	} else if getLineCount(filename) != lineCount { // test that the last file has the correct number of lines
		t.Errorf("incorrect number of lines in files")
	} else if _, err := os.Open(testFilePrefix + "_2.txt"); err == nil { // test it is not creating more than the necessary files
		t.Errorf("creating files that shouldn't exist")
	} else if _, err := os.Open(testFilePrefix + "_1.txt"); err != nil { // test it is not deleting files before it is necessary
		t.Errorf("Deleting files before needed.")
	}
	cleanUpTests(*L)

	L2, _ := newLogger(2, 2, testFilePrefix)
	LT2, _ := newLogTrialler(L2)
	LT2.LogTrial(5, true)
	filename = testFilePrefix + "_" + strconv.Itoa(fileCount) + ".txt"
	if rotateCount != 1 || lineCount != 1 || fileCount != 0 {
		t.Errorf("not handling rotations correctly.")
	} else if _, err := os.Open(filename); err != nil {
		t.Errorf("Not creating files correctly.")
	} else if getLineCount(filename) != lineCount {
		t.Errorf("incorrect number of lines in files")
	} else if _, err := os.Open(testFilePrefix + "_2.txt"); err == nil {
		t.Errorf("creating files that shouldn't exist")
	} else if _, err := os.Open(testFilePrefix + "_1.txt"); err != nil {
		t.Errorf("Deleting files before needed.")
	}
	cleanUpTests(*L2)

	L3, _ := newLogger(1, 11, testFilePrefix)
	LT3, _ := newLogTrialler(L3)
	LT3.LogTrial(10, true)
	filename = testFilePrefix + "_" + strconv.Itoa(fileCount) + ".txt"
	if rotateCount != 0 || lineCount != 10 || fileCount != 0 {
		t.Errorf("not handling rotations correctly.")
	} else if _, err := os.Open(filename); err != nil {
		t.Errorf("Not creating files correctly.")
	} else if getLineCount(filename) != lineCount {
		t.Errorf("incorrect number of lines in files")
	} else if _, err := os.Open(testFilePrefix + "_1.txt"); err == nil {
		t.Errorf("creating files that shouldn't exist")
	} else if _, err := os.Open(testFilePrefix + "_0.txt"); err != nil {
		t.Errorf("Deleting files before needed.")
	}
	cleanUpTests(*L3)

	L4, _ := newLogger(8, 1, testFilePrefix)
	LT4, _ := newLogTrialler(L4)
	LT4.LogTrial(10, true)
	filename = testFilePrefix + "_" + strconv.Itoa(fileCount) + ".txt"
	if rotateCount != 1 || lineCount != 1 || fileCount != 1 {
		t.Errorf("not handling rotations correctly.")
	} else if _, err := os.Open(filename); err != nil {
		t.Errorf("Not creating files correctly.")
	} else if getLineCount(filename) != lineCount {
		t.Errorf("incorrect number of lines in files")
	} else if _, err := os.Open(testFilePrefix + "_8.txt"); err == nil {
		t.Errorf("creating files that shouldn't exist")
	} else if _, err := os.Open(testFilePrefix + "_0.txt"); err != nil {
		t.Errorf("Deleting files before needed.")
	}
	cleanUpTests(*L4)

	L5, _ := newLogger(1, 1, testFilePrefix)
	LT5, _ := newLogTrialler(L5)
	LT5.LogTrial(1, true)
	fmt.Println(rotateCount, lineCount, fileCount)
	filename = testFilePrefix + "_" + strconv.Itoa(fileCount) + ".txt"
	if rotateCount != 1 || lineCount != 0 || fileCount != 0 {
		t.Errorf("not handling rotations correctly.")
	} else if _, err := os.Open(filename); err != nil {
		t.Errorf("Not creating files correctly.")
	} else if getLineCount(filename) != 1 { // 1 here as filename never changes and it hasn't been reset yet
		t.Errorf("incorrect number of lines in files")
	}
	cleanUpTests(*L5)
}

// helper function that deletes used test files, to reset
func cleanUpTests(L logger) {
	var amountOfFiles int
	if rotateCount == 0 {
		amountOfFiles = fileCount
	} else {
		amountOfFiles = L.MAX_FILES
	}

	for i := 0; i < amountOfFiles; i++ {
		fileName := L.FILE_PREFIX + "_" + strconv.Itoa(i) + ".txt"
		err := os.Remove(fileName)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	reset()
}

// helper function that gets the amount of lines in a file
func getLineCount(filename string) int {
	file, _ := os.Open(filename)
	fileScanner := bufio.NewScanner(file)
	lineCountInFile := 0
	for fileScanner.Scan() {
		lineCountInFile++
	}
	return lineCountInFile
}

// are logTrialling and logger similar enough to test using log trialler? should i change it to just add x amount of messages then call exit? no, wouldnt work, would exit before
