package main

import (
	"evaluator/internal"
	"flag"
	"log"
)

func main() {
	labName := flag.String("n", "", "laboratory name")
	flag.Parse()

	ctx, err := internal.SetUpContext(*labName)
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
