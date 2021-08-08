package util

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"syscall"
	"time"
)

func writeTTYRecHeader(fd io.Writer, length int) {
	t := time.Now()
	tv := syscall.NsecToTimeval(t.UnixNano())
	binary.Write(fd, binary.LittleEndian, int32(tv.Sec))
	binary.Write(fd, binary.LittleEndian, int32(tv.Usec))
	binary.Write(fd, binary.LittleEndian, int32(length))
}

func WriteTtyRecData(fd io.Writer, data []byte) {
	if len(data) > 0 {
		writeTTYRecHeader(fd, len(data))
		fd.Write(data)
	}
}

func GetHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return hostname
}

func GetUserName() string {
	user, err := user.Current()
	if err != nil {
		log.Println(err)
		return ""
	}
	return user.Username
}
