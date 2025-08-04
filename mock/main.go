package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RequestLog struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body,omitempty"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}

	headers := make(map[string]string)
	for name, values := range r.Header {
		headers[name] = values[0]
	}

	logEntry := RequestLog{
		Method:  r.Method,
		Path:    r.URL.Path,
		Headers: headers,
		Body:    string(body),
	}

	jsonLog, err := json.Marshal(logEntry)
	if err != nil {
		log.Printf(`{"error":"failed to encode log","reason":"%s"}`, err)
	} else {
		log.Println(string(jsonLog))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", handler)
	log.Println(`{"message": "Mock service listening on :5678"}`)
	http.ListenAndServe(":5678", nil)
}
