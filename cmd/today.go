package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kristofarkas/today"
)

var unit int
var duration time.Duration
var verbose bool

func main() {

	flag.BoolVar(&verbose, "v", false, "verbose output")

	addSet := flag.NewFlagSet("add", flag.ExitOnError)
	addSet.IntVar(&unit, "unit", 0, "total number of units")
	addSet.DurationVar(&duration, "duration", time.Hour * 24 * 7,  "total allocated time")

	updateSet := flag.NewFlagSet("update", flag.ExitOnError)
	updateSet.IntVar(&unit, "unit", 0, "current unit count")
	updateSet.DurationVar(&duration, "ext", 0,  "extend from today by")

	flag.Parse()

	switch flag.Arg(0) {
	case "add":
		err := addSet.Parse(flag.Args()[1:])
		if err != nil {
			log.Fatal(err)
		}

		name := strings.Join(addSet.Args(), " ")

		err = today.Save(&today.Task{
			Name:    name,
			Total:   unit,
			EndDate: time.Now().Add(duration),
		})

		if err != nil {
			log.Fatalf("failed to save: %v\n", err)
		}

		if verbose {
			today.Print(os.Stdout, verbose)
		}

	case "update":
		err := updateSet.Parse(flag.Args()[1:])
		if err != nil {
			log.Fatal(err)
		}

		name := strings.Join(updateSet.Args(), " ")

		t, err := today.Load(name)
		if err != nil {
			log.Fatalf("failed to load task: %s", err)
		}

		if unit > 0 {
			if err = t.SetCurrent(unit); err != nil {
				log.Fatalf("failed to set current: %s", err)
			}
		}

		if duration > 0 {
			t.EndDate = time.Now().Add(duration)
		}

		if err = today.Save(t); err != nil {
			log.Fatalf("failed to save task: %s", err)
		}

		if verbose {
			today.Print(os.Stdout, verbose)
		}

	case "delete":
		name := strings.Join(flag.Args()[1:], " ")

		if err := today.Delete(name); err != nil {
			log.Fatalf("failed to delete task: %s", err)
		}

		if verbose {
			today.Print(os.Stdout, verbose)
		}

	default:
		today.Print(os.Stdout, verbose)
	}

}
