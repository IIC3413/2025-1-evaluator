package main

import (
	"evaluator/internal"
	"flag"
	"log"
	"os"
)

func exitWithErr(err error) {
	log.Printf("Error: %s\n", err.Error())
	os.Exit(1)
}

func main() {
	confPath := flag.String("c", "config.yaml", "configuration file path")
	flag.Parse()

	conf, err := internal.OpenConfig(*confPath)
	if err != nil {
		exitWithErr(err)
	}

	ctx, err := internal.SetUpContext(conf)
	if err != nil {
		exitWithErr(err)
	}

	if err = internal.Run(ctx); err != nil {
		exitWithErr(err)
	}
}
