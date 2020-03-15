package todo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// item struct ToDo item
type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

// list of ToDo items
type List []item

// Add creates new todo and append to list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, t)
}

// Complete marks a todo as complete
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}
	// Adjusting index for 0 based index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()
	return nil
}

// Delete deletes a todo from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("Item %d does not exist", i)
	}
	// Adjusting index for 0 based index
	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

// Save encodes the List as JSON and saves it
func (l *List) Save(filename string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, js, 0644)
}

// Get opens the file, decodes the JSON
func (l *List) Get(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, l)
}

//String prints formatted list
func (l *List) String(v bool, c bool) string {
	formatted := ""
	for k, t := range *l {
		if c && t.Done {
			continue
		}
		prefix := " "
		if t.Done {
			prefix = "X "
		}
		// Adjust the item index
		if !v {
			formatted += fmt.Sprintf("%s%d: %s\n", prefix, k+1, t.Task)
		} else {
			formatted += fmt.Sprintf("%s%d: %s - %s\n", prefix, k+1, t.Task, t.CreatedAt.Format(time.ANSIC))
		}
	}
	return formatted
}
