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
	err := beamx.Run(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}

