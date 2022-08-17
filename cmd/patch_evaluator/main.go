package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DaoCloud-OpenSource/patch_evaluator"
	pflag "github.com/spf13/pflag"
)

var (
	kind string // code or doc
)

func init() {
	pflag.StringVar(&kind, "kind", "code", "code or doc")
	pflag.Parse()
}

type ValueDefine struct {
	LowValue []patch_evaluator.Filterer
	NoValue  []patch_evaluator.Filterer
}

var (
	define = map[string]ValueDefine{
		"code": {
			LowValue: []patch_evaluator.Filterer{
				patch_evaluator.FocusSuffixFilterer{".sh", ".bash", ".c", ".go", ".py", ".java", ".cpp", ".h", ".hpp", ".yaml", ".yml"},
				patch_evaluator.PrefixFilterer{"test/", "tests/"},
				patch_evaluator.SuffixFilterer{"_test.go"},
			},
			NoValue: []patch_evaluator.Filterer{
				patch_evaluator.SuffixFilterer{".md"},
				patch_evaluator.PrefixFilterer{"vendor/"},
				patch_evaluator.ContainsFilterer{"generated", "testdata"},
				patch_evaluator.CommentFilterer{},
				patch_evaluator.EmptyLineFilterer{},
			},
		},
		"doc": {
			NoValue: []patch_evaluator.Filterer{
				patch_evaluator.PrefixFilterer{
					"content/de", "content/es", "content/fr", "content/hi", "content/id", "content/it", "content/ja", "content/ko", "content/no", "content/pl", "content/pt-br", "content/ru", "content/uk", "content/vi",
				},
				patch_evaluator.CommentFilterer{},
				patch_evaluator.EmptyLineFilterer{},
			},
		},
	}
)

func main() {
	valueDefine, ok := define[kind]
	if !ok {
		log.Fatalf("kind %q not found", kind)
	}

	args := pflag.Args()
	if len(args) < 1 {
		log.Fatal("need args")
	}

	patch, err := os.Open(args[0])
	if err != nil {
		log.Fatal(err)
	}
	defer patch.Close()

	evaluator := patch_evaluator.Evaluator{}
	files, reasonsLowValue, reasonsNoValue, err := evaluator.Evaluate(patch, valueDefine.LowValue, valueDefine.NoValue)
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
