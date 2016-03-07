// stopwatch records when a named event is started and prints out the duration with it is stopped.
// The data is stored as a json file in ~/.stopwatch.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

const version string = "0.1"

type stopwatch struct {
	Version string
	Timers  []timer
}

type timer struct {
	Label string
	Start time.Time
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage for stopwatch version %s:\n", version)
		fmt.Println("stopwatch       # prints all existing stopwatches")
		fmt.Println("stopwatch label # starts a stopwatch or stops a stopwatch with that name")
		fmt.Println("")
		flag.PrintDefaults()
	}
}

// There is no file locking, so two processed running at the same time could cause a problem.
// I might also win the lottery.
func main() {
	st, err := newStopwatch()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	stopAllPtr := flag.Bool("stopall", false, "Stops all stopwatches")
	flag.Parse()

	if *stopAllPtr == true {
		st.stopAll()
		return
	}

	l := len(os.Args)
	if l < 2 {
		st.list()
		return
	}

	label := strings.Join(os.Args[1:], " ")

	i, _ := st.find(label)
	if i >= 0 {
		st.stop(i)
		return
	}

	st.start(label)
}

func (t timer) toString() string {
	d := time.Now().Sub(t.Start)
	d = ((d + time.Second/2) / time.Second) * time.Second
	return fmt.Sprintf("%s %s (%s)\n", t.Label, d, t.Start.Round(time.Second))
}

func (t timer) stop() {
	fmt.Printf("stopped %s", t.toString())
}

func newStopwatch() (*stopwatch, error) {
	if _, err := os.Stat(filepath()); os.IsNotExist(err) {
		return &stopwatch{Version: version}, nil
	}

	body, err := ioutil.ReadFile(filepath())

	if err != nil {
		return nil, err
	}

	st := stopwatch{}
	err = json.Unmarshal(body, &st)
	st.Version = version

	if err != nil {
		return nil, err
	}

	return &st, err
}

func (st *stopwatch) start(label string) error {
	t := timer{label, time.Now()}
	st.Timers = append(st.Timers, t)
	st.write()

	fmt.Printf("started %s\n", label)
	return nil
}

func (st *stopwatch) find(label string) (int, *timer) {
	for i, t := range st.Timers {
		if t.Label == label {
			return i, &t
		}
	}

	return -1, nil
}

func (st *stopwatch) stop(pos int) {
	t := st.Timers[pos]
	st.Timers = append(st.Timers[:pos], st.Timers[pos+1:]...)
	t.stop()
	st.write()
}

func (st *stopwatch) stopAll() {
	for _, t := range st.Timers {
		t.stop()
	}
	st.Timers = []timer{}
	st.write()
}

func (st *stopwatch) list() {
	if len(st.Timers) == 0 {
		fmt.Println("No stopwatches exist")
	}

	for _, t := range st.Timers {
		fmt.Printf(t.toString())
	}
}

func (st *stopwatch) write() error {
	b, err := json.MarshalIndent(st, "", "  ")
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
