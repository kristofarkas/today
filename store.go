package today

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os/user"
	"strings"
	"text/tabwriter"
)

// Save saves a task to store.
func Save(t *Task) error {

	ts, err := LoadAll()
	if err != nil {
		return err
	}

	for i := range ts {
		if ts[i].Name == t.Name {
			ts = append(ts[:i], ts[i+1:]...)
			break
		}
	}

	ts = append(ts, t)

	data, err := json.Marshal(ts)
	if err != nil {
		return err
	}

	usr, err := user.Current()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/today/tasks.json", usr.HomeDir), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a named task if present
func Delete(name string) error {

	usr, err := user.Current()
	if err != nil {
		return err
	}

	ts, err := LoadAll()
	if err != nil {
		return err
	}

	for i := range ts {
		if ts[i].Name == name {
			ts = append(ts[:i], ts[i+1:]...)
			break
		}
	}

	data, err := json.Marshal(ts)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/today/tasks.json", usr.HomeDir), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Load loads a task if present
func Load(name string) (*Task, error) {
	ts, err := LoadAll()
	if err != nil {
		return nil, err
	}

	var fuzzyMatch []*Task

	for _, t := range ts {
		if t.Name == name {
			return t, nil
		} else if strings.HasPrefix(strings.ToLower(t.Name), strings.ToLower(name)) {
			fuzzyMatch = append(fuzzyMatch, t)
		}
	}

	if len(fuzzyMatch) == 1 {
		return fuzzyMatch[0], nil
	}

	return nil, fmt.Errorf("%s not found", name)
}

// LoadAll loads all the tasks in store.
func LoadAll() ([]*Task, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/today/tasks.json", usr.HomeDir))
	if err != nil {
		return nil, err
	}

	var t []*Task

	if err = json.Unmarshal(data, &t); err != nil {
		return nil, err
	}

	return t, nil
}

func Print(output io.Writer, verbose bool) {
	ts, err := LoadAll()

	if err != nil {
		return
	}

	w := tabwriter.NewWriter(output, 0, 8, 2, ' ', 0)

	if verbose {
		const format = "%v\t%v\t%v\t%v\t%v\t%v\n"

		fmt.Fprintf(w, format, "Name", "Current", "Today", "To go", "Progress", "P/D")
		fmt.Fprintf(w, format, "----", "-------", "-----", "-----", "--------", "---")

		for i := len(ts) - 1; i >= 0; i-- {
			t := ts[i]
			fmt.Fprintf(w, format, t.Name, t.Current, t.Today(), t.ToGo(), fmt.Sprintf("%d%%", t.Progress()), t.UnitsPerDay())
		}
	} else {
		const format = "%v\t%v\t%v\n"

		fmt.Fprintf(w, format, "Name", "Current", "Today")
		fmt.Fprintf(w, format, "----", "-------", "-----")

		for i := len(ts) - 1; i >= 0; i-- {
			t := ts[i]
			fmt.Fprintf(w, format, t.Name, t.Current, t.Today())
		}
	}

	w.Flush()
}
