package db

import (
	"fmt"
	"strconv"
	"testing"
)

// TodoList tests...

func newList() TodoList {
	tlist := TodoList{}
	values := []string{"foo", "bar", "biz", "fang", "boop"}

	for i, name := range values {
		todo := Todo{}
		todo.Name = name
		todo.Description = strconv.Itoa(i)
		todo.ID = i
		todo.Available = true

		tlist.Value = append(tlist.Value, todo)
	}
	return tlist
}

func TestEvery(t *testing.T) {

	fmt.Println("TodoList.Every")
	tlist := newList()
	total := 0
	expected := 5
	tlist.Every(func(value *Todo) {
		fmt.Printf("[ %d : %s ]\n", value.ID, value.Name)
		total++
	})
	if total != expected {
		t.Fatalf("expected=%d, got=%d", expected, total)
	}
}

func TestFilter(t *testing.T) {
	fmt.Println("TodoList.Filter : Name[0] == 'f'")

	tlist := newList()
	expected := 2
	filtered := tlist.Filter(func(value *Todo) bool {
		return value.Name[0] == 'f'
	})

	filtered.Every(func(value *Todo) {
		fmt.Printf("[ %d : %s ]\n", value.ID, value.Name)
	})

	if len(filtered.Value) != expected {
		t.Fatalf("expected=%d, got=%d", expected, len(filtered.Value))
	}
}
