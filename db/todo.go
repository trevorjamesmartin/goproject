package db

import (
	"database/sql"
	"log"

	"github.com/a-h/templ"
	"github.com/trevorjamesmartin/goproject/components"
)

// Todo ...
type Todo struct {
	ID          int
	Name        string
	Description string
	Available   bool
}

// ListItem ... <li>{todo}</li>
func (t *Todo) ListItem() templ.Component {
	return components.ListItemTodo(t.ID, t.Name, t.Description, t.Available)
}

// Save ...
func (t *Todo) Save(store *sql.DB) {
	if len(t.Name) < 1 || len(t.Description) < 1 {
		return
	}

	if t.ID == -1 {
		insertTodo(store, Todo{t.ID, t.Name, t.Description, t.Available})
	} else {
		updateTodo(store, Todo{t.ID, t.Name, t.Description, t.Available})
	}
}

// Delete ... removefrom database
func (t *Todo) Delete(store *sql.DB) {
	if t.ID <= 0 {
		return
	}
	deleteTodo(store, Todo{t.ID, t.Name, t.Description, t.Available})
}

// TodoList Struct ...
type TodoList struct {
	Value []Todo
}

func (tl *TodoList) Read(store *sql.DB) {
	rows, err := store.Query(`SELECT id, name, description, available FROM todo ORDER BY created DESC`)

	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var name, description string
	var id int
	var available bool

	tl.Value = []Todo{}

	for rows.Next() {
		rows.Scan(&id, &name, &description, &available)
		if err != nil {
			log.Fatal(err)
		}
		td := Todo{}
		td.ID = id
		td.Name = name
		td.Description = description
		td.Available = available
		tl.Value = append(tl.Value, td)
	}
}

// GetTodo ... returns a single record
func GetTodo(store *sql.DB, id int) Todo {
	query := `SELECT name, description, available FROM todo WHERE id = $1`
	var name, description string
	var available bool
	err := store.QueryRow(query, id).Scan(&name, &description, &available)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Printf("no rows found with ID %d", id)
		}
		return Todo{-1, "", "", false}

	}
	return Todo{id, name, description, available}
}

func createTodoTable(store *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS todo(
	id SERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	description VARCHAR(255) NOT NULL,
	available BOOLEAN,
	created TIMESTAMP DEFAULT NOW()
	)`
	_, err := store.Exec(query)

	if err != nil {
		log.Fatal(err)
	}
}

func insertTodo(store *sql.DB, todo Todo) int {

	query := `INSERT INTO todo (name, description, available)
	  VALUES ($1, $2, $3) RETURNING id`

	var pk int

	err := store.QueryRow(query, todo.Name, todo.Description, todo.Available).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk
}

func updateTodo(store *sql.DB, todo Todo) Todo {

	query := `UPDATE todo 
	SET name=$2, description=$3, available=$4
	WHERE id=$1
	RETURNING name, description, available`

	var available bool
	var name, description string

	err := store.QueryRow(query, todo.ID, todo.Name, todo.Description, todo.Available).Scan(&name, &description, &available)
	if err != nil {
		log.Fatal(err)
	}

	todo.Name = name
	todo.Description = description
	todo.Available = available

	return todo
}

func deleteTodo(store *sql.DB, todo Todo) {
	query := `DELETE FROM todo WHERE id=$1`
	store.Query(query, todo.ID)
}
