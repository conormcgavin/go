package main

import "os"

func main() {
	myLogger := NewLogger(0, 0, "log_output")
	c, exit := myLogger.StartLogListener()

	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			myLogger.AddMsg(c, arg)
		}
	}

	/*
		myLogger.LogTrial(12)
		myLogger.AddMsg(c, "hello")
		myLogger.AddMsg(c, "my")
		myLogger.AddMsg(c, "name")
		myLogger.AddMsg(c, "is")
		myLogger.AddMsg(c, "conor")
		myLogger.AddMsg(c, "mcgavin")
	*/

	myLogger.CallExit(exit)
}
