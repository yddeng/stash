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

	customCfg := &config.Config{}

	flagSet := flag.NewFlagSet(exe, flag.ExitOnError)

	bindFlags(customCfg, flagSet)

	flagSet.Usage = func() {
		fmt.Fprintf(flagSet.Output(), "Usage: %s [flag...] file [outfile]\n", exe)
		flagSet.PrintDefaults()
	}

	flagSet.Parse(os.Args[1:])

	if flagSet.NArg() < 1 || flagSet.NArg() > 2 {
		flagSet.Usage()
		os.Exit(1)
	}

	outCfg := &config.Config{}
	if _, err := toml.DecodeFile(flagSet.Arg(0), outCfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	flagSet.Visit(func(f *flag.Flag) {
		if fb := flagBindings[f.Name]; fb == nil {
			panic(fmt.Errorf("flag %s not bound", f.Name))
		} else {
			fb.copy(outCfg, customCfg)
		}
	})

	var outfile string
	if flagSet.NArg() == 1 {
		outfile = flagSet.Arg(0)
	} else {
		outfile = flagSet.Arg(1)
	}

	if file, err := os.Create(outfile); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else {
		encoder := toml.NewEncoder(file)
		if err = encoder.Encode(outCfg); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		file.Close()
	}
}
