// stopwatch records when a named event is started and prints out the duration with it is stopped.
// The data is stored as a json file in ~/.stopwatch.json
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"
)

type Stopwatch struct {
	Timers []Timer
}

type Timer struct {
	Label string
	Start time.Time
}

// There is no file locking, so two processed running at the same time could cause a problem.
// I might also win the lottery.
func main() {
	st, err := newStopwatch()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	l := len(os.Args)
	if l < 2 {
		st.list()
		return
	}

	label := os.Args[1]

	i, _ := st.find(label)
	if i >= 0 {
		st.stop(i)
		return
	}

	st.start(label)
}

func (t Timer) toString() string {
	d := time.Now().Sub(t.Start)
	return fmt.Sprintf("'%s' running: %s started: %s\n", t.Label, d, t.Start)
}

func newStopwatch() (*Stopwatch, error) {
	if _, err := os.Stat(filepath()); os.IsNotExist(err) {
		return &Stopwatch{}, nil
	}

	body, err := ioutil.ReadFile(filepath())

	if err != nil {
		return nil, err
	}

	st := Stopwatch{}
	err = json.Unmarshal(body, &st)

	if err != nil {
		return nil, err
	}

	return &st, err
}

func (st *Stopwatch) start(label string) error {
	t := Timer{label, time.Now()}
	st.Timers = append(st.Timers, t)
	st.write()

	fmt.Printf(t.toString())
	return nil
}

func (st *Stopwatch) find(label string) (int, *Timer) {
	for i, t := range st.Timers {
		if t.Label == label {
			return i, &t
		}
	}

	return -1, nil
}

func (st *Stopwatch) stop(pos int) {
	t := st.Timers[pos]
	st.Timers = append(st.Timers[:pos], st.Timers[pos+1:]...)
	st.write()
	fmt.Printf("Stopped %s", t.toString())
}

func (st *Stopwatch) list() {
	if len(st.Timers) == 0 {
		fmt.Println("No stopwatches exist")
	}

	for _, t := range st.Timers {
		fmt.Printf(t.toString())
	}
}

func (st *Stopwatch) write() error {
	b, err := json.Marshal(st)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	err = ioutil.WriteFile(filepath(), b, 0644)

	if err != nil {
		fmt.Printf("%s", err)
		return err
	}

	return nil
}

func filepath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir + "/.stopwatch.json"
}
