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

const version string = "0.2"

type app struct {
	Version     string
	Stopwatches []stopwatch
}

type stopwatch struct {
	Label string
	Start time.Time
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage for stopwatch version %s:\n", version)
		fmt.Println("   stopwatch       # prints all existing stopwatches")
		fmt.Println("   stopwatch label # starts a stopwatch or stops a stopwatch with that name")
		fmt.Println("\nFlags:")
		flag.PrintDefaults()
	}
}

// There is no file locking, so two processed running at the same time could cause a problem.
// I might also win the lottery.
func main() {
	a, err := load()
	if err != nil {
		fmt.Printf("%s", err)
		return
	}

	a.parseArgs()
}

func (a *app) parseArgs() {
	stopAllPtr := flag.Bool("stopall", false, "Issues a stop command all stopwatches")
	flag.Parse()

	if *stopAllPtr == true {
		a.stopAll()
		return
	}

	if flag.NArg() == 0 {
		a.list()
		return
	}

	label := strings.Join(flag.Args(), " ")

	if i := a.find(label); i >= 0 {
		a.stop(i)
		return
	}

	a.start(label)
}

func (t stopwatch) toString() string {
	d := time.Now().Sub(t.Start)
	d = ((d + time.Second/2) / time.Second) * time.Second
	return fmt.Sprintf("%s %s (%s)\n", t.Label, d, t.Start.Round(time.Second))
}

func (t stopwatch) stop() {
	fmt.Printf("stopped %s", t.toString())
}

func load() (*app, error) {
	if _, err := os.Stat(filepath()); os.IsNotExist(err) {
		return &app{Version: version}, nil
	}

	body, err := ioutil.ReadFile(filepath())

	if err != nil {
		return nil, err
	}

	a := app{}
	err = json.Unmarshal(body, &a)
	a.Version = version

	if err != nil {
		return nil, err
	}

	return &a, err
}

func (a *app) start(label string) error {
	t := stopwatch{label, time.Now()}
	a.Stopwatches = append(a.Stopwatches, t)
	a.save()

	fmt.Printf("started %s\n", label)
	return nil
}

func (a *app) find(label string) int {
	for i, t := range a.Stopwatches {
		if t.Label == label {
			return i
		}
	}

	return -1
}

func (a *app) stop(pos int) {
	a.Stopwatches[pos].stop()
	a.Stopwatches = append(a.Stopwatches[:pos], a.Stopwatches[pos+1:]...)
	a.save()
}

func (a *app) stopAll() {
	for _, t := range a.Stopwatches {
		t.stop()
	}
	a.Stopwatches = []stopwatch{}
	a.save()
}

func (a *app) list() {
	if len(a.Stopwatches) == 0 {
		fmt.Println("No stopwatches exist")
	}

	for _, t := range a.Stopwatches {
		fmt.Printf(t.toString())
	}
}

func (a *app) save() error {
	b, err := json.MarshalIndent(a, "", "  ")
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
