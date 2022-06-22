package utilities

import (
	"log"
	"runtime"
	"time"
)

func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, _, line, _ := runtime.Caller(1)

		log.SetFlags(log.Lmsgprefix)
		log.Printf("------------------------------------------------------------------------------------\n[Error]\t\t: %s\n[Line]\t\t: %d\n[Time]\t\t: %s\n[Message]\t: %v", runtime.FuncForPC(pc).Name(), line, time.Now().Format(time.RFC850), err)
		// PostToWebHook(err)
		b = true
	}
	return
}
