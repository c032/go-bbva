package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/c032/go-bbva"
)

func main() {
	flagset := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagset.Usage = func() {
		fmt.Fprintf(flagset.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flagset.Output(), "\n")
		fmt.Fprintf(flagset.Output(), "    %s FILE\n", os.Args[0])

		flagset.PrintDefaults()
	}
	flagset.Parse(os.Args[1:])

	args := flagset.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "Expected to receive one file.\n")

		os.Exit(2)
	} else if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Only one file was expected, but received %d files.\n", len(args))

		os.Exit(2)
	}

	file := args[0]
	enc := json.NewEncoder(os.Stdout)

	xlsx, err := bbva.ParseXLSXFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err.Error())

		os.Exit(1)
	}

	for _, item := range xlsx.Items {
		outputItem := map[string]string{}

		for _, cell := range item.Cells {
			key := cell[0]
			value := cell[1]

			outputItem[key] = value
		}

		err := enc.Encode(outputItem)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}
	}
}
