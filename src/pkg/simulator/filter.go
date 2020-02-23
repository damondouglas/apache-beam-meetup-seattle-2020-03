package simulator

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"regexp"
)

func Filter(s beam.Scope, tokens beam.PCollection, ignore ...string) beam.PCollection {
	return beam.ParDo(s, handler(ignore...), tokens)
}

func handler(ignore ...string) func(string, func(string)) {
	p := make([]*regexp.Regexp, len(ignore))
	for i := range ignore {
		p[i] = regexp.MustCompile(ignore[i])
	}
	return func(token string, emit func(string)) {
		for _, k := range p {
			if k.MatchString(token) {
				return
			}
		}
		emit(token)
	}
}

