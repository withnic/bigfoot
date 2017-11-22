package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	t   = flag.Int("t", 5, "It is openning time.")
	m   = flag.Bool("m", false, "It is multi invoke.")
	par = flag.Int("par", 10, "It is browser num.")
)

func main() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage of this:
		-t  :It is openning time.
		-m  :It is multi invoke.
		-par  :it is browser num.`)
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Not enough Args:%d please input file name", len(args))
		os.Exit(1)
	}
	os.Exit(run(args[0]))
}

func run(name string) int {
	po, err := NewBigfoot(name, *m, *par, time.Duration(*t)*time.Second)
	if err != nil {
		fmt.Fprintf(os.Stderr, "New Bigfoot Error: %v", err)
		return 1
	}
	err = po.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Start Error: %v", err)
		return 1
	}
	defer po.Stop()

	err = po.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Run Error: %v", err)
		return 1
	}

	return 0
}
