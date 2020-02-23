package main

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/stats"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"log"
	"os"
	"regexp"
)

const (
	inputKey = "INPUT"
	outputKey = "OUTPUT"
)

var (
	wordRE  = regexp.MustCompile(`[a-zA-Z]+('[a-z])?`)
	input = os.Getenv(inputKey)
	output = os.Getenv(outputKey)
)

func init() {
	for _, k := range []string{} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}

func main() {
	beam.Init()
	p, s := beam.NewPipelineWithRoot()
	lines := textio.Read(s, input)
	words := beam.ParDo(s, func(line string, emit func(string)) {
		for _, word := range wordRE.FindAllString(line, -1) {
			emit(word)
		}
	}, lines)
	counted := stats.Count(s, words)
	formatted := beam.ParDo(s, func(w string, c int) string {
		return fmt.Sprintf("%s: %v", w, c)
	}, counted)
	textio.Write(s, output, formatted)
	err := beamx.Run(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}

