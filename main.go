package main

// go standard library
import (
	"fmt"
	"net/http"
)

func main() {
	serveMux := http.NewServeMux()
	serveMux.Handle("/", http.FileServer(http.Dir(".")))
	server := http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
	}
}
