package App

import (
	"bytes"
	"fmt"
)

var NickName = "Unkown"

func Log(message string, a ...interface{})  {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("[%s] ", NickName))
	buffer.WriteString(message)
	buffer.WriteString("\n")
	fmt.Printf(buffer.String(), a...)
}