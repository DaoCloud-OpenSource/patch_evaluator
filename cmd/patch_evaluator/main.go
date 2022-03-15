package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DaoCloud-OpenSource/patch_evaluator"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("need args")
	}

	patch, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer patch.Close()

	evaluator := patch_evaluator.Evaluator{}
	files, reasons, err := evaluator.Evaluate(patch)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("valuable: %d\n", len(files))
	for _, file := range files {
		fmt.Printf("\t%s\n", file.NewName)
	}

	fmt.Printf("valueless: %d\n", len(reasons))
	for _, reason := range reasons {
		fmt.Printf("\t%s: %s\n", reason.File, reason.Message)
	}
}
