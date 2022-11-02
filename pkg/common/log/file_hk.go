package log

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	logFilePrefix       = "./logs/log_"
	DisplayFileTimeMode = "2006-01-02"
)

var FileWriter *logFileWriter

// setupFileLog
func setupFileLog() {
	FileWriter = &logFileWriter{}
	logrus.SetOutput(FileWriter)
	gin.DefaultWriter = FileWriter
}

// logFileWriter
type logFileWriter struct {
	file *os.File
	size int64
}

// Write
func (fw *logFileWriter) Write(data []byte) (n int, err error) {
	if fw == nil {
		return 0, errors.New("logFileWriter is nil")
	}
	if fw.file == nil {
		setLogFileWriter()
	}

	var bb bytes.Buffer
	bb.Write(data)
	bb.WriteByte('\n')
	n, err = fw.file.Write(bb.Bytes())
	fw.size += int64(n)
	checkLogFile()
	return
}

// Print
func (fw *logFileWriter) Print(v ...interface{}) {
	if fw == nil {
		return
	}
	if fw.file == nil {
		setLogFileWriter()
	}

	var bb bytes.Buffer
	bb.WriteString("gorm: ")
	switch v[0] {
	case "sql":
		bb.WriteString(fmt.Sprintf("src[%+v]\n", v[1]))
		bb.WriteString(fmt.Sprintf("duration[%+v] ", v[2]))
		bb.WriteString(fmt.Sprintf("sql[%+v]\n", v[3]))
		bb.WriteString(fmt.Sprintf("values[%+v]\n", v[4]))
		bb.WriteString(fmt.Sprintf("rows_returned[%v]", v[5]))
	case "log":
		bb.WriteString(fmt.Sprintf("log[%+v]", v[2]))
	}
	bb.WriteByte('\n')
	bb.WriteByte('\n')
	n, _ := fw.file.Write(bb.Bytes())
	fw.size += int64(n)
	checkLogFile()
	return
}

// setLogFileWriter
func setLogFileWriter() {
	file, err := os.OpenFile(logFilePrefix+convertTimeToCustomString(time.Now(), DisplayFileTimeMode),
		os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
	if err != nil {
		log.Fatal("log  init failed")
	}

	var info os.FileInfo
	info, err = file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	FileWriter.file = file
	FileWriter.size = info.Size()
}

// checkLogFile
func checkLogFile() {
	if FileWriter.size > 1024*64 {
		_ = FileWriter.file.Close()
		FileWriter.file, _ = os.OpenFile(logFilePrefix+convertTimeToCustomString(time.Now(), DisplayFileTimeMode),
			os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0600)
		FileWriter.size = 0
	}

	historyFilePath := logFilePrefix + convertTimeToCustomString(time.Now().Add(time.Hour*-1*24*7), DisplayFileTimeMode)
	_, err := os.Stat(historyFilePath)
	if os.IsExist(err) {
		_ = os.Remove(historyFilePath)
	}
}

// convertTimeToCustomString
func convertTimeToCustomString(val time.Time, model string) string {
	return val.Format(model)
}
