package logger

import (
	"fmt"
)

func main() {
	myLogger, err := NewLogger(0, 0, "log_output")
	if err != nil {
		fmt.Println(err.Error())
	}
	c, exit := myLogger.StartLogListener()
	myLogger.AddMsg(c, "hello", 4, true, "no")

	/*
		myLogTrialler, err := NewLogTrialler(myLogger)

		if err != nil {
			fmt.Println(err.Error())
		}
	*/
	//myLogTrialler.LogTrial(50)
	/*
		if len(os.Args) > 1 {
			for _, arg := range os.Args[1:] {
				myLogger.AddMsg(c, arg)
			}
		}
	*/

	/*
		myLogger.AddMsg(c, "hello")
		myLogger.AddMsg(c, "my")
		myLogger.AddMsg(c, "name")
		myLogger.AddMsg(c, "is")
		myLogger.AddMsg(c, "conor")
		myLogger.AddMsg(c, "mcgavin")
		myLogger.AddMsg(c, "mcgavin")
		myLogger.AddMsg(c, "hello")
		myLogger.AddMsg(c, "my")
		myLogger.AddMsg(c, "name")
		myLogger.AddMsg(c, "is")
		myLogger.AddMsg(c, "conor")
		myLogger.AddMsg(c, "mcgavin")
		myLogger.AddMsg(c, "hello")
		myLogger.AddMsg(c, "my")
		myLogger.AddMsg(c, "name")
		myLogger.AddMsg(c, "is")
		myLogger.AddMsg(c, "conor")
		myLogger.AddMsg(c, "mcgavin")
		myLogger.AddMsg(c, "hello")
		myLogger.AddMsg(c, "my")
		myLogger.AddMsg(c, "name")
		myLogger.AddMsg(c, "is")
		myLogger.AddMsg(c, "conor")
		myLogger.AddMsg(c, "mcgavin")
	*/

	myLogger.CallExit(exit)
}

// next, implement it so anything that has a toString function can enter a log msg in
// consider Write interface? no?
