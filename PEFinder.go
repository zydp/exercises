package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

var (
	input    = flag.String("i", "", "Specify the input file")
	datafile = flag.String("w", "pefinder.txt", "change to directory DIR")
	logfile  = flag.String("log", "", "change to directory DIR")
	sigs     = make(chan os.Signal, 1)
	exit     = false
)

const (
	GUESS_SUFFIX      = "  ---Perhaps\n"
	FILE_NAME_MIN_LIN = 5
)

func init() {
	flag.CommandLine.Usage = help
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	if flag.NFlag() < 1 {
		help()
		os.Exit(0)
	}
	if *logfile != "" {
		SetOutPutLog(*logfile)
	}
	go listenSignal()

	listfiles := strings.Split(*input, ",")
	peoutfd, err := os.OpenFile(*datafile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if nil != err {
		log.Println("Can't write to file, Please check", err)
		os.Exit(0)
	}
	defer peoutfd.Close()
	pewriter := bufio.NewWriter(peoutfd)

	zipoutfd, err := os.OpenFile(*datafile+"_zip.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if nil != err {
		log.Println("Can't write to file, Please check", err)
		os.Exit(0)
	}
	defer zipoutfd.Close()
	zipwriter := bufio.NewWriter(zipoutfd)

	for _, listfile := range listfiles {
		if exit {
			fmt.Printf("Received exit signal, ready to exit")
			break
		}
		if err := isExist(listfile); nil != err {
			log.Println("Check File Failed:", err)
			continue
		}
		listfd, err := os.Open(listfile)
		if err != nil {
			log.Println("Open File Failed:", err)
			continue
		}
		defer listfd.Close()
		log.Println("Checking File", listfile, "..")

		reader := bufio.NewReader(listfd)
		for {
			if exit {
				fmt.Printf("Received exit signal, ready to exit")
				break
			}
			file, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			} else if err != nil {
				log.Println("Read File Failed:", err)
				continue
			}
			if len(file) < FILE_NAME_MIN_LIN {
				continue
			}
			file = strings.TrimSpace(file)
			/*is zip file?*/
			if strings.HasSuffix(file, "zip") {
				zipwriter.WriteString(file + "\n")
				zipwriter.Flush()
				continue
			}
			if err := isExist(file); nil != err {
				log.Println("Check File Failed:", err)
				continue
			}
			filefd, err := os.Open(file)
			if err != nil {
				log.Println("Open File Failed:", err)
				continue
			}
			defer filefd.Close()
			buffer := make([]byte, 3)
			n, err := filefd.Read(buffer)
			if n < 2 {
				log.Println(file, "is a empty file")
				continue
			}
			// fmt.Printf("%X -- %X  -- %X -- %X \n", buffer[0], buffer[1], buffer[240], buffer[241])
			// continue
			if 0x4d == buffer[0] && 0x5a == buffer[1] {
				if _, err = filefd.Seek(240, 0); nil != err {
					pewriter.WriteString(file + GUESS_SUFFIX)
					pewriter.Flush()
					continue
				}
				buffer2 := make([]byte, 3)
				filefd.Read(buffer2)
				if n < 2 {
					pewriter.WriteString(file + GUESS_SUFFIX)
				} else if 0x50 == buffer2[0] && 0x45 == buffer2[1] {
					pewriter.WriteString(file + "\n")
				} else {
					pewriter.WriteString(file + GUESS_SUFFIX)
				}
				pewriter.Flush()
			}
		}
	}
	log.Println("Check Over!")
}

func isExist(filename string) (err error) {
	_, err = os.Stat(filename)
	return
}

func help() {
	fmt.Printf("Provide find PE file interface based on laboratory environment.\n")
	fmt.Printf("Usage: %s [OPTION]...\n", filepath.Base(os.Args[0]))
	fmt.Println("     -i\t\tSpecify the intput files, required")
	fmt.Println("     -w\t\tSpecify the output file, The default output is pefinder.txt")
	fmt.Println("     -log\tSpecify the logfile,The default output is on the screen")
	fmt.Println("     -help\tdisplay help info and exit")
	fmt.Println("example:")
	fmt.Printf("	%s -i file1,file2..", filepath.Base(os.Args[0]))
	fmt.Println()
}

func handleSignals(signal os.Signal) {
	log.Println("Recv a signal:", signal)
	exit = true
	os.Exit(0)
}

func listenSignal() {
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	for {
		sig := <-sigs
		handleSignals(sig)
	}
}

func SetOutPutLog(logfilename string) {
	if logout, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		log.SetOutput(logout)
		log.SetPrefix("[Info] ")
		//log.SetFlags(log.Ldate | log.Ltime /*| log.Lshortfile*/)
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println(time.Now().Format(time.RFC3339), strings.Title(runtime.GOARCH), strings.Title(runtime.GOOS))
	} else {
		fmt.Println(time.Now().Format(time.RFC3339), strings.Title(runtime.GOARCH), strings.Title(runtime.GOOS))
	}
}
