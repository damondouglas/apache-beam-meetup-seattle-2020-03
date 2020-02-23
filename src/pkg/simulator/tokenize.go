package simulator

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"strings"
)

var (
	tokenColumns = []int{1,2,4}
)

// Tokenize a beam.PCollection expected as
// dbid    drug    known   misspellingPhon edOne   edTwo   pillMark        google_ms       google_title    google_snippet  ud_slang        missed
// See https://zenodo.org/record/3238718#.Xk9hOpNKjOQ
// Tab delimited columns 2, 3 and 5 are extracted and split by comma
func Tokenize(s beam.Scope, lines beam.PCollection) beam.PCollection {
	tabbedTokens := beam.ParDo(s, tabbedLineHandler, lines)
	return beam.ParDo(s, commaDelimitedHandler, tabbedTokens)
}

func tabbedLineHandler(line string, emit func(string)) {
	if line == "" {
		return
	}
	tabbedTokens := strings.Split(line, "\t")
	for _, k := range tokenColumns {
		emit(tabbedTokens[k])
	}
}

func commaDelimitedHandler(segment string, emit func(string)) {
	tokens := strings.Split(segment, ",")
	for _, k := range tokens {
		emit(k)
	}
}