package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type AddRequest struct {
	Addon int
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s", r.URL.Path)
	host, _ := os.Hostname()
	fmt.Fprintf(w, "Hello, world!\n")
	fmt.Fprintf(w, "Version: 1.0.0\n")
	fmt.Fprintf(w, "Hostname: %s\n", host)
}

func GetHeaders(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving request: %s with %s", r.URL.Path, r.Method)

	fmt.Println("Query")
	queries := r.URL.Query()
	if val, ok := queries["key"]; ok {
		fmt.Println(len(queries), val)
	} else {
		fmt.Fprintf(w, "No key in query: %s\n", "key")
	}

	for k, v := range r.URL.Query() {
		fmt.Fprintf(w, "%s: %s\n", k, v)
	}

	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func MultiHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		queries := r.URL.Query()
		if val, ok := queries["a"]; ok {
			fmt.Println(len(queries), val)
			i, _ := strconv.Atoi(val[0])
			fmt.Fprintf(w, "result: %d\n", i*2)
		} else {
			fmt.Fprintf(w, "No key in query: %s\n", "key")
		}

	case "POST":
		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", requestBody)

		//testStr := `{"Addon":123}`
		var addRequest AddRequest
		json.Unmarshal(requestBody, &addRequest)
		fmt.Printf("result: %d\n", addRequest.Addon)

		w.Write([]byte("Receive request"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/headers", GetHeaders)
	mux.HandleFunc("/add", MultiHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
