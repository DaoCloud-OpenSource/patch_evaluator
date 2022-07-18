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
	files, reasonsLowValue, reasonsNoValue, err := evaluator.Evaluate(patch)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("valuable: %d\n", len(files))
	for _, file := range files {
		fmt.Printf("\t%s\n", file.NewName)
	}

	fmt.Printf("low-value: %d\n", len(reasonsLowValue))
	for _, reason := range reasonsLowValue {
		fmt.Printf("\t%s: %s\n", reason.File, reason.Message)
	}

	fmt.Printf("no-value: %d\n", len(reasonsNoValue))
	for _, reason := range reasonsNoValue {
		fmt.Printf("\t%s: %s\n", reason.File, reason.Message)
	}
}
