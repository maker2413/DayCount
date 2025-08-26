package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/maker2413/daycount/internal/compute"
	"github.com/maker2413/daycount/internal/parse"
	"github.com/maker2413/daycount/internal/server"
)

func main() {
	var serve bool
	var reference string
	var formatOutput bool

	fs := flag.NewFlagSet("daycount", flag.ExitOnError)
	fs.StringVar(&reference, "reference", "", "reference 'today' date (YYYY-MM-DD) instead of now (useful for testing)")
	fs.StringVar(&reference, "r", "", "shorthand for --reference")
	fs.BoolVar(&formatOutput, "format", false, "print a friendly sentence instead of raw integer")
	fs.BoolVar(&formatOutput, "f", false, "shorthand for --format")
	fs.BoolVar(&serve, "serve", false, "run HTTP server instead of CLI")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage: daycount <date> [--reference YYYY-MM-DD] [--format]\n")
		fmt.Fprintln(fs.Output(), "Supported input formats: 2006-01-02, 2006/01/02, 01/02/2006, RFC3339, and common datetime variants.")
		fs.PrintDefaults()
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	if serve {
		// Run server mode and ignore positional args.
		if len(fs.Args()) != 0 {
			fs.Usage()
			os.Exit(2)
		}
		if err := server.Run(context.Background(), ""); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	args := fs.Args()
	if len(args) != 1 {
		fs.Usage()
		os.Exit(2)
	}
	input := args[0]

	target, err := parse.ParseDate(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	var now time.Time
	if reference != "" {
		ref, err := parse.ParseDate(reference)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid reference date: %v\n", err)
			os.Exit(1)
		}
		now = ref
	} else {
		now = time.Now().UTC()
	}

	days := compute.DaysSince(target, now)
	if formatOutput {
		phrase := "days since"
		if days < 0 {
			phrase = "days until"
		}
		fmt.Printf("%d %s %s\n", abs(days), phrase, input)
		return
	}
	fmt.Println(days)
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
