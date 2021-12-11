package eqlog

import (
	"fmt"

	"github.com/akfarrington/earthquake_bot/cwbapi"
)

// LogWriter custom log writer
type LogWriter struct {
}

func (writer LogWriter) Write(bytes []byte) (int, error) {
	return fmt.Print(cwbapi.GetTwTime() + ": " + string(bytes))
}
