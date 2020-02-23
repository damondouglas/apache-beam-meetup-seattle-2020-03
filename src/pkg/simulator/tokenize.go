package simulator

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/filter"
	"strings"
)

// Tokenize a beam.PCollection by splitting at indicated columns.
// Column numbers a zero based.
func Tokenize(s beam.Scope, lines beam.PCollection, column ...int) beam.PCollection {
	tabbedTokens := beam.ParDo(s, tabbedLineHandler(column...), lines)
	tokens := beam.ParDo(s, commaDelimitedHandler, tabbedTokens)
	return filter.Distinct(s, tokens)
}

func tabbedLineHandler(column ...int) func(string, func(string)) {
	return func(line string, emit func(string)) {
		if line == "" {
			return
		}
		tokens := strings.Split(line, "\t")
		for _, k := range column {
			emit(tokens[k])
		}
	}
}

func commaDelimitedHandler(segment string, emit func(string)) {
	tokens := strings.Split(segment, ",")
	for _, k := range tokens {
		emit(k)
	}
}