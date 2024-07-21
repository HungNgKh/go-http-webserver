package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Node struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdon"`
}

var nodes = make(map[string]Node)

var id int = 0

// POST node - /api/nodes
func postNodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "POST /nodes")
}

// GET node - /api/nodes
func getNodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "GET /nodes")
}

// PUT node - /api/nodes/{id}
func putNodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PUT /nodes/{id}")
}

// DELETE node - /api/nodes/{id}
func deleteNodeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "DELETE /nodes/{id}")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/nodes", getNodeHandler).Methods("GET")
	router.HandleFunc("/api/nodes", postNodeHandler).Methods("POST")
	router.HandleFunc("/api/nodes/{id}", putNodeHandler).Methods("PUT")
	router.HandleFunc("/api/nodes/{id}", deleteNodeHandler).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}

	log.Println("Listening on port 8080...")

	server.ListenAndServe()
}
