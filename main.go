package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Node struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdon"`
	UpdatedOn   time.Time `json:"updatedon"`
}

var nodesStore = make(map[string]Node)

var id int = 0

// POST node - /api/nodes
func postNodeHandler(w http.ResponseWriter, r *http.Request) {
	var node Node

	err := json.NewDecoder(r.Body).Decode(&node)

	if err != nil {
		panic(err)
	}

	node.CreatedOn = time.Now()
	node.UpdatedOn = node.CreatedOn

	id++
	idString := strconv.Itoa(id)
	nodesStore[idString] = node

	j, err := json.Marshal(node)

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// GET node - /api/nodes
func getNodeHandler(w http.ResponseWriter, r *http.Request) {
	var nodes []Node

	for _, node := range nodesStore {
		nodes = append(nodes, node)
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(nodes)

	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// PUT node - /api/nodes/{id}
func putNodeHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	vars := mux.Vars(r)
	id := vars["id"]
	var updatedNode Node
	var status int

	err = json.NewDecoder(r.Body).Decode(&updatedNode)

	if err != nil {
		panic(err)
	}

	if node, ok := nodesStore[id]; ok {
		updatedNode.CreatedOn = node.CreatedOn
		updatedNode.UpdatedOn = time.Now()
		delete(nodesStore, id)
		nodesStore[id] = updatedNode
		status = http.StatusOK
	} else {
		log.Printf("Could not find node %s to update", id)
		status = http.StatusNotFound
	}

	j, err := json.Marshal(nodesStore[id])

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(j)
}

// DELETE node - /api/nodes/{id}
func deleteNodeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var status int

	if _, ok := nodesStore[id]; ok {
		delete(nodesStore, id)
		status = http.StatusNoContent
	} else {
		log.Printf("Could not find key %s to delete", id)
		status = http.StatusNotFound
	}

	w.WriteHeader(status)
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
