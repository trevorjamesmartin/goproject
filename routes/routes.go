// Package routes ...
package routes

import (
	// Go ^1.22
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	// local
	"github.com/trevorjamesmartin/goproject/components"
	"github.com/trevorjamesmartin/goproject/db"
	"github.com/trevorjamesmartin/goproject/library"
)

func pathNumber(key string, r *http.Request) int {
	pathID := r.PathValue(key)
	id64, err := strconv.ParseInt(pathID, 10, 32)
	if err != nil {
		log.Println(err)
		return 0
	}
	return int(id64)
}

func urlKey(key string, r *http.Request) (bool, string) {
	u, err := r.URL.Parse(r.URL.String())

	if err != nil {
		log.Fatal(err)
	}

	m, _ := url.ParseQuery(u.RawQuery)
	if len(m[key]) > 0 {
		return true, m[key][0]
	}
	return false, ""
}

func serveLocalFile(fname string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		f, err := filepath.Abs(fname)
		result, err := os.ReadFile(f)
		if err != nil {
			http.NotFound(w, r)
		}
		fmt.Fprintf(w, "%s", result)
	}
}

// NewRouter ...:
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /class-tools.js", serveLocalFile("./static/class-tools.js"))
	mux.HandleFunc("GET /htmx.min.js", serveLocalFile("./static/htmx.min.js"))

	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/api/todo", todoListHandler)
	mux.HandleFunc("/api/todo/{id}", todoHandler)

	mux.HandleFunc("POST /api/todo/toggle/{id}", todoToggle)
	mux.HandleFunc("GET /api/todo/edit/{id}", todoEdit)
	mux.HandleFunc("PUT /api/todo/update/{id}", todoUpdate)

	return mux
}

func todoUpdate(w http.ResponseWriter, r *http.Request) {
	id := pathNumber("id", r)
	r.ParseForm()
	name := r.FormValue("name")
	description := r.FormValue("description")

	store := db.Connect()
	defer store.Close()

	todo := db.GetTodo(store, id)

	if len(name) > 0 {
		todo.Name = name
	}

	if len(description) > 0 {
		todo.Description = description
	}
	todo.Save(store)

	err := todo.ListItem().Render(context.Background(), w)
	if err != nil {
		http.NotFound(w, r)
	}
}

func todoToggle(w http.ResponseWriter, r *http.Request) {
	id := pathNumber("id", r)
	store := db.Connect()
	defer store.Close()
	todo := db.GetTodo(store, id)
	todo.Available = !todo.Available
	todo.Save(store)
	err := todo.ListItem().Render(context.Background(), w)
	if err != nil {
		http.NotFound(w, r)
	}
}

func todoEdit(w http.ResponseWriter, r *http.Request) {
	id := pathNumber("id", r)
	store := db.Connect()
	defer store.Close()
	todo := db.GetTodo(store, id)
	x := components.TodoEdit(todo.ID, todo.Name, todo.Description, todo.Available)
	err := x.Render(context.Background(), w)
	if err != nil {
		http.NotFound(w, r)
	}
}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	id := pathNumber("id", r)
	store := db.Connect()
	defer store.Close()

	switch r.Method {
	case "DELETE":
		todo := db.GetTodo(store, id)
		if todo.ID > 0 {
			todo.Delete(store)
		}
		break

	case "GET":
		todo := db.GetTodo(store, id)

		_, formatTo := urlKey("format", r)

		switch formatTo {
		case "json":
			w.Write([]byte(library.PrettyPrintJSON(todo)))
			break
		default:
			err := todo.ListItem().Render(context.Background(), w)
			if err != nil {
				http.NotFound(w, r)
			}
		}
	default:
		break
	}
}

func todoListHandler(w http.ResponseWriter, r *http.Request) {
	store := db.Connect()
	defer store.Close()

	if r.Method == "POST" {
		r.ParseForm()
		name := "*new Todo"
		description := "description..."
		t := db.Todo{}
		t.Name = name
		t.Description = description
		t.ID = -1
		t.Available = true
		t.Save(store)
	}

	tl := db.TodoList{}
	tl.Read(store)

	switch r.PathValue("format") {
	case "json":
		fmt.Fprintln(w, library.PrettyPrintJSON(tl.Value))
		break
	default:
		w.Write([]byte("<ul id=\"todo-list\">"))
		for i := 0; i < len(tl.Value); i++ {
			todo := tl.Value[i]
			err := todo.ListItem().Render(context.Background(), w)
			if err != nil {
				log.Fatalf("failed to render TODO: %v", err)
			}
		}
		w.Write([]byte("</ul>"))
		break
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	store := db.Connect()
	defer store.Close()

	switch path {
	case "/":
		btn := components.NewTodo()
		x := components.ContentPage("Todo", btn)
		err := x.Render(context.Background(), w)
		if err != nil {
			log.Fatalf("failed to write output file: %v", err)
		}
		break
	default:
		http.NotFound(w, r)
		return
	}

	tl := db.TodoList{}

	tl.Read(store)
	w.Write([]byte("<ul id=\"todo-list\">"))
	for i := 0; i < len(tl.Value); i++ {
		todo := tl.Value[i]
		err := todo.ListItem().Render(context.Background(), w)
		if err != nil {
			log.Fatalf("failed to render TODO: %v", err)
		}
	}
	w.Write([]byte("</ul>"))

}
