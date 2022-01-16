package main

import (
	"notion/generate"
	"notion/insert"
	"notion/update"

	"flag"
	"log"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("not enough arguments given")
	}
	args := os.Args[1:]

	// flag setting
	var flags [3]bool
	f := flag.FlagSet{}
	f.BoolVar(&flags[0], "g", false, "runs generate")
	f.BoolVar(&flags[1], "i", false, "runs insert")
	f.BoolVar(&flags[2], "u", false, "runs update")
	// parse flags from args
	f.Parse(args)

	tmpC := 0
	for _, v := range flags {
		if v {
			tmpC++
		}
	}

	if tmpC > 1 {
		log.Fatal("more than one method chosen")
	}

	// remove flag
	args = args[1:]

	if flags[0] {
		err := generate.Generate(args)
		if err != nil {
			log.Fatal(err)
		}
	} else if flags[1] {
		err := insert.Insert(args)
		if err != nil {
			log.Fatal(err)
		}
	} else if flags[2] {
		err := update.Update(args)
		if err != nil {
			log.Fatal(err)
		}
	}
}
