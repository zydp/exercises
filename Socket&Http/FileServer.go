package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	LOG_FILE = "./logs/"
)

type RespMsg struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

var (
	iservices = flag.Bool("s", false, "To running as a services")
	port      = flag.Int("port", 3256, "The TCP port that the server listens on")
	address   = flag.String("address", "", "The net address that the server listens")
	url_path  = flag.String("path", "/file/", "Specify url path")
	ssl_crt   = flag.String("ssl_ctr", "", "Specify SSL/TLS certificate files")
	ssl_key   = flag.String("ssl_key", "", "Specify SSL/TLS key files")
	dataPath  = flag.String("data", "./data", "Specify the voice file path")
	sigs      = make(chan os.Signal, 1)
	exit      = make(chan bool, 1)
)

func init() {
	flag.CommandLine.Usage = help

	if logout, err := os.OpenFile(LOG_FILE+strings.TrimSuffix(filepath.Base(os.Args[0]), ".exe")+".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		log.SetOutput(logout)
		log.SetPrefix("[Info] ")
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		log.Println(time.Now().Format(time.RFC3339), strings.Title(runtime.GOARCH), strings.Title(runtime.GOOS))
	} else {
		fmt.Println(time.Now().Format(time.RFC3339), strings.Title(runtime.GOARCH), strings.Title(runtime.GOOS))
		fmt.Printf("Not found the [ logs ] directory, the log will be displayed on the terminal\n")
	}
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	runAsServices()
	if flag.NFlag() <= 0 {
		fmt.Printf("Using default setting, listen on %s:%d\n", *address, *port)
		log.Printf("Using default setting, listen on %s:%d\n", *address, *port)
	}

	router := http.NewServeMux()

	router.Handle(*url_path, http.StripPrefix(*url_path, http.FileServer(http.Dir(*dataPath))))

	fmt.Printf("Service listen on %s:%d:%s\n", *address, *port, *url_path)
	log.Printf("Service listen on %s:%d:%s\n", *address, *port, *url_path)

	go listenSignal()
	if len(*ssl_crt) < 4 || len(*ssl_key) < 4 {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", *address, *port), router); err != nil {
			log.Println(err)
		}
	} else {
		if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%d", *address, *port), "server.crt", "server.key", router); err != nil {
			log.Println(err)
		}
	}
}

func safe_http_handle(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err, ok := recover().(error)
			if ok {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}()
		fn(w, r)
	}
}

func help() {
	fmt.Printf("Provide data upload interface based on laboratory environment.\n")
	fmt.Printf("Usage: %s [OPTION]...\n", filepath.Base(os.Args[0]))
	fmt.Println("     -s\t\tSet process running as a services, using [false] by default")
	fmt.Println("     -address\tSet the listener address, using [0.0.0.0] by default")
	fmt.Println("     -port\tSet the listener port, using port [3244] by default")
	fmt.Println("     -path\tSpecify url path, using [ /file/ ] by default")
	fmt.Println("     -ssl_ctr\tSpecify SSL/TLS certificate files")
	fmt.Println("     -ssl_key\tSpecify SSL/TLS key files")
	fmt.Println("     -data\tSpecify the voice file path, default ./data/")
	fmt.Println("     -help\tdisplay help info and exit")
	fmt.Println()
}

func runAsServices() {
	if *iservices {
		cmd := exec.Command(os.Args[0], flag.Args()...)
		cmd.Start()
		fmt.Printf("%s [PID] %d running...\n", filepath.Base(os.Args[0]), cmd.Process.Pid)
		log.Printf("%s [PID] %d running...\n", filepath.Base(os.Args[0]), cmd.Process.Pid)
		*iservices = false
		os.Exit(0)
	}
}

func handleSignals(signal os.Signal) {
	log.Println("Recv a signal:", signal)
	exit <- true
	os.Exit(0)
}

func listenSignal() {
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGABRT)
	for {
		sig := <-sigs
		handleSignals(sig)
	}
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
