package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served: %s\n", r.Host)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is:"
	fmt.Fprintf(w, "<h1 align=\"center\">%s</h1>", Body)
	fmt.Fprintf(w, "<h2 align=\"center\">%s</h2>n", t)
	fmt.Fprintf(w, "Serving: %s\n", r.URL.Path)
	fmt.Printf("Served time for: %s\n", r.Host)
}

/* the test cmd on linux: time nc localhost 3244 */

func main() {

	PORT := ":3244"
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Printf("Listening on https://0.0.0.0%s\n", PORT)
	} else {
		PORT = ":" + arguments[1]
		fmt.Printf("Listening on https://0.0.0.0%s\n", PORT)
	}

	m := http.NewServeMux()
	srv := &http.Server{
		Addr:         PORT,
		Handler:      m,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
	m.HandleFunc("/time", timeHandler)
	m.HandleFunc("/", myHandler)

	//err := srv.ListenAndServe()

	//in tls
	err := srv.ListenAndServeTLS("server.crt", "server.key")
	
	if err != nil {
		fmt.Println(err)
		return
	}
}
