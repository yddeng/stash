package main

import (
	"flag"
	"fmt"
	"initialthree/robot/config"
	"os"

	"github.com/BurntSushi/toml"
)

func main() {
	exe := os.Args[0]

	cfg := &config.Config{}

	flagSet := flag.NewFlagSet(exe, flag.ExitOnError)

	bindFlags(cfg, flagSet)

	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Usage: %s [flag...] outfile\n", exe)
		flagSet.PrintDefaults()
	}

	flagSet.Parse(os.Args[1:])

	if flagSet.NArg() < 1 {
		flagSet.Usage()
		os.Exit(1)
	}

	file, err := os.Create(flagSet.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	encoder := toml.NewEncoder(file)
	if err = encoder.Encode(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	file.Close()
}
