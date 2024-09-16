package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type Node struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createdon"`
	UpdatedOn   time.Time `json:"updatedon"`
}

type EditNode struct {
	Node
	Id string
}

var nodesStore = make(map[string]Node)

var id int = 0

var templates map[string]*template.Template

// GET nodes - /
func getNodes(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index", "base", nodesStore)
}

// GET nodes - /add
func addNode(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "add", "base", nil)
}

// node/save for new item
func saveNode(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	title := r.PostFormValue("title")
	desc := r.PostFormValue("description")
	node := Node{Title: title, Description: desc, CreatedOn: time.Now(), UpdatedOn: time.Now()}

	id++
	k := strconv.Itoa(id)
	nodesStore[k] = node
	http.Redirect(w, r, "/", http.StatusFound)
}

// handler for /node/edit/{id} to editing existing item
func editNode(w http.ResponseWriter, r *http.Request) {
	var viewModel EditNode

	vars := mux.Vars(r)
	k := vars["id"]
	if node, ok := nodesStore[k]; ok {
		viewModel = EditNode{node, k}
	} else {
		http.Error(w, "Could not find resource to edit", http.StatusBadRequest)
	}

	renderTemplate(w, "edit", "base", viewModel)
}

// handler for  /node/update/{id} to updating existing item
func updateNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	var nodeToUpdate Node

	if node, ok := nodesStore[k]; ok {
		r.ParseForm()
		nodeToUpdate.Title = r.PostFormValue("title")
		nodeToUpdate.Description = r.PostFormValue("description")
		nodeToUpdate.CreatedOn = node.CreatedOn
		nodeToUpdate.UpdatedOn = time.Now()

		delete(nodesStore, k)
		nodesStore[k] = nodeToUpdate
	} else {
		http.Error(w, "Could not find resource to update", http.StatusBadRequest)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

// DELETE node handler - /node/delete/{id}
func deleteNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var status int

	if _, ok := nodesStore[id]; ok {
		delete(nodesStore, id)
		status = http.StatusFound
	} else {
		log.Printf("Could not find key %s to delete", id)
		status = http.StatusNotFound
	}

	http.Redirect(w, r, "/", status)
}

func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templates["index"] = template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))
	templates["add"] = template.Must(template.ParseFiles("templates/add.html", "templates/base.html"))
	templates["edit"] = template.Must(template.ParseFiles("templates/edit.html", "templates/base.html"))
}

func renderTemplate(w http.ResponseWriter, name string, template string, viewModel interface{}) {
	tmpl, ok := templates[name]
	if !ok {
		http.Error(w, "Template does not exists", http.StatusInternalServerError)
	}

	err := tmpl.ExecuteTemplate(w, template, viewModel)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func LoggingHandler(next http.Handler) http.Handler {
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return handlers.LoggingHandler(logFile, next)
}

func iconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

func main() {
	router := mux.NewRouter()

	commonHandlers := alice.New(LoggingHandler, handlers.CompressHandler)

	router.Handle("/", alice.New(LoggingHandler, handlers.CompressHandler).ThenFunc(http.HandlerFunc(getNodes))).Methods("GET")
	router.HandleFunc("/favicon.ico", iconHandler)
	router.Handle("/nodes/add", commonHandlers.ThenFunc(http.HandlerFunc(addNode))).Methods("GET")
	router.Handle("/nodes/save", commonHandlers.ThenFunc(http.HandlerFunc(saveNode))).Methods("POST")
	router.Handle("/nodes/edit/{id}", commonHandlers.ThenFunc(http.HandlerFunc(editNode))).Methods("GET")
	router.Handle("/nodes/update/{id}", commonHandlers.ThenFunc(http.HandlerFunc(updateNode))).Methods("POST")
	router.Handle("/nodes/delete/{id}", commonHandlers.ThenFunc(http.HandlerFunc(deleteNode)))

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
