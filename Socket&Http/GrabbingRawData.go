package main

import (
	"fmt"
	"os"
	"syscall"
)
/*in linux successed*/
func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_TCP)
	if err != nil {
		fmt.Println("Error in syscall.Socket:", err)
		return
	}

	f := os.NewFile(uintptr(fd), "IPPROTO_TCP")
	if f == nil {
		fmt.Println("Error in os.NewFile:", err)
		return
	}
	err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 2048)
	if err != nil {
		fmt.Println("Error in syscall.Socket:", err)
		return
	}

	for {
		buf := make([]byte, 2048)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println("Read failed:", err)
		}
		fmt.Printf("% X\n", buf[:numRead])
	}
}
