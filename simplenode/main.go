package main

import (
	"fmt"
	"io"
	"net/http"
)

// fill at build-time using linker flags
var strToAppend string

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := fmt.Sprintf("%s%s", string(body), strToAppend)
	w.Write([]byte(response))
}

func main() {
	fmt.Println(strToAppend)
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
