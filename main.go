package main

func main() {
	myLogger := NewLogger(8, 8, "log_output")
	//myLogger.LogTrial(46)
	c, exit := myLogger.StartLogListener()
	myLogger.AddMsg(c, "hello")
	myLogger.AddMsg(c, "my")
	myLogger.AddMsg(c, "name")
	myLogger.AddMsg(c, "is")
	myLogger.AddMsg(c, "conor")
	myLogger.CallExit(exit)
}
