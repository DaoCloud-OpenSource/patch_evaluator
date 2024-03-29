package patch_evaluator

import (
	"io"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type Evaluator struct {
}

func (e Evaluator) Evaluate(r io.Reader, lowValue []Filterer, noValue []Filterer) ([]*gitdiff.File, []*Reasons, []*Reasons, error) {
	files, _, err := gitdiff.Parse(r)
	if err != nil {
		return nil, nil, nil, err
	}

	filtered := []*gitdiff.File{}

	reasonsLowValue := []*Reasons{}

	reasonsNoValue := []*Reasons{}
loop:
	for _, file := range files {
		for _, filterer := range noValue {
			if rea := filterer.Filter(file); rea != nil {
				reasonsNoValue = append(reasonsNoValue, rea)
				continue loop
			}
		}
		for _, filterer := range lowValue {
			if rea := filterer.Filter(file); rea != nil {
				reasonsLowValue = append(reasonsLowValue, rea)
				continue loop
			}
		}
		filtered = append(filtered, file)
	}
	return filtered, reasonsLowValue, reasonsNoValue, nil
}

var contentDiff = diffmatchpatch.New()

type Filterer interface {
	Filter(file *gitdiff.File) *Reasons
}

type Reasons struct {
	File    string
	Message string
}
type StringsModifyFilterer struct{}

func (s StringsModifyFilterer) Filter(file *gitdiff.File) *Reasons {

	total := 0
	stringModiry := 0
	for _, fragments := range file.TextFragments {
		for i, line := range fragments.Lines {
			if line.Op != gitdiff.OpAdd {
				continue
			}
			if i == 0 {
				continue
			}
			prvLine := fragments.Lines[i-1]
			if prvLine.Op != gitdiff.OpDelete {
				continue
			}

			diffs := contentDiff.DiffMain(prvLine.Line, line.Line, false)
			if len(diffs) != 3 {
				continue
			}

			total++
			if strings.ContainsAny(diffs[0].Text, `"'`) && strings.ContainsAny(diffs[2].Text, `"'`) {
				stringModiry++
			}
		}
	}
	if total != 0 && total == stringModiry {
		return &Reasons{
			File:    file.NewName,
			Message: "only string modified",
		}
	}
	return nil
}

type EmptyLineFilterer struct{}

func (s EmptyLineFilterer) Filter(file *gitdiff.File) *Reasons {
	if file.IsNew {
		return nil
	}
	for _, fragments := range file.TextFragments {
		for _, line := range fragments.Lines {
			if line.Op != gitdiff.OpContext {
				continue
			}
			if strings.TrimSpace(line.Line) == "" {
				continue
			}
			return nil
		}
	}
	return &Reasons{
		File:    file.NewName,
		Message: "not modified",
	}
}

type CommentFilterer struct{}

func (s CommentFilterer) Filter(file *gitdiff.File) *Reasons {
	for _, fragments := range file.TextFragments {
		for _, line := range fragments.Lines {
			if line.Op != gitdiff.OpAdd {
				continue
			}
			if strings.TrimSpace(line.Line) == "" {
				continue
			}
			if strings.HasPrefix(strings.TrimSpace(line.Line), "//") {
				continue
			}
			if strings.HasPrefix(strings.TrimSpace(line.Line), "#") {
				continue
			}
			return nil
		}
	}
	return &Reasons{
		File:    file.NewName,
		Message: "All lines are comments",
	}
}

type FocusSuffixFilterer []string

func (s FocusSuffixFilterer) Filter(file *gitdiff.File) *Reasons {
	for _, v := range s {
		if strings.HasSuffix(file.NewName, v) {
			return nil
		}
	}
	return &Reasons{
		File:    file.NewName,
		Message: "not focused",
	}
}

type SuffixFilterer []string

func (s SuffixFilterer) Filter(file *gitdiff.File) *Reasons {
	for _, v := range s {
		if strings.HasSuffix(file.NewName, v) {
			return &Reasons{
				File:    file.NewName,
				Message: "suffix " + v,
			}
		}
	}
	return nil
}

type PrefixFilterer []string

func (s PrefixFilterer) Filter(file *gitdiff.File) *Reasons {
	for _, v := range s {
		if strings.HasPrefix(file.NewName, v) {
			return &Reasons{
				File:    file.NewName,
				Message: "prefix " + v,
			}
		}
	}
	return nil
}

type ContainsFilterer []string

func (s ContainsFilterer) Filter(file *gitdiff.File) *Reasons {
	for _, v := range s {
		if strings.Contains(file.NewName, v) {
			return &Reasons{
				File:    file.NewName,
				Message: "contains " + v,
			}
		}
	}
	return nil
}

type NotFilter struct {
	Filterer Filterer
}

func (s NotFilter) Filter(file *gitdiff.File) *Reasons {
	s.Filterer.Filter(file)
	return nil
}
