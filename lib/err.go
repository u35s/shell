package lib

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
)

func CheckError(err error) bool {
	if err != nil {
		log.Printf("Error:%v", err.Error())
		return false
	}
	return true
}

func DumpStack(srvName string, id uint) {
	if err := recover(); err != nil {
		var buf bytes.Buffer
		bs := make([]byte, 1<<12)
		num := runtime.Stack(bs, false)
		buf.WriteString(fmt.Sprintf("Panic: %s\n", err))
		buf.Write(bs[:num])
		log.Printf(buf.String())
	}
}
