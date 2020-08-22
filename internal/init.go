package internal

import (
	"io"
	"log"
	"os"
)

var logger = log.New(os.Stdout, "Parser:\t", log.Ldate|log.Ltime)

//SetLoggerOutput 设置日志输出
func SetLoggerOutput(out io.Writer) {
	logger.SetOutput(out)
}
