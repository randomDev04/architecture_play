package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    // r *http.Request contains EVERYTHING about the incoming request
    fmt.Printf("Method: %s\n", r.Method)        // GET, POST, PUT, DELETE
    fmt.Printf("Path: %s\n", r.URL.Path)        // /api/tasks
    fmt.Printf("Query: %s\n", r.URL.RawQuery)   // ?page=1&limit=10
    fmt.Printf("Headers: %v\n", r.Header)       // All HTTP headers
    
    // w http.ResponseWriter is how you WRITE the response
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK) // 200
    fmt.Fprintf(w, `{"message": "Request received"}`)
}



func main() {
    // This function runs for EVERY request to the server
    http.HandleFunc("/", handler)
    http.ListenAndServe(":8080", nil)
}