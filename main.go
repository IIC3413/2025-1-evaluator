package main

import (
	"evaluator/internal"
	"flag"
	"log"
	"os"
)

const errPrefix = "Error: %s"

func main() {
	labName := flag.String("n", "", "laboratory name")
	flag.Parse()

	ctx, err := internal.SetUpContext(*labName)
	if err != nil {
		log.Printf(errPrefix, err)
		os.Exit(1)
	}

	e, err := internal.NewEvaluator(ctx)
	if err != nil {
		log.Printf(errPrefix, err)
		os.Exit(1)
	}

	if err = e.Eval(); err != nil {
		log.Printf(errPrefix, err)
		os.Exit(1)
	}
}
