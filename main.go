package main

import (
	"evaluator/internal"
	"flag"
	"log"
)

func main() {
	labName := flag.String("n", "", "laboratory name")
	mode := flag.String("m", "Release", "compilation mode")
	flag.Parse()

	if labName == nil || *labName == "" {
		log.Fatal("Missing lab name flag")
	}

	ctx, err := internal.SetUpContext(*labName, *mode)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	e, err := internal.NewEvaluator(ctx)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	if err = e.Eval(); err != nil {
		log.Fatalf("Error: %s", err)
	}

	if err = e.FreeLogs(); err != nil {
		log.Fatalf("Error: %s", err)
	}
}
