package main

// go standard library
import (
	"log"
	"net/http"
)

type apiHandler struct{}

func main() {
	const filePathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	serveMux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir(filePathRoot))))
	serveMux.HandleFunc("/healthz/", healthzHandler)

	server := http.Server{
		Handler: serveMux,
		Addr:    ":" + port,
	}

	log.Printf("Serving files from %v and listening on port %v", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
