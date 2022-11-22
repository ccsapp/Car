package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Printf("Hello world!")
	http.HandleFunc("/", HelloServer)
	http.ListenAndServe(":80", nil)
	fmt.Printf("done with main")

}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s!", r.URL.Path[1:])
	fmt.Printf("Hello %s!", r.URL.Path[1:])
}

