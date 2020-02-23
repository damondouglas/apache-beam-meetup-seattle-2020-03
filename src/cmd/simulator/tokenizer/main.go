package main

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"log"
	"os"
	"strconv"
	"strings"
	"temp/pkg/simulator"
)

const (
	inputKey = "INPUT"
	outputKey = "OUTPUT"
	sourceKey = "SOURCE"
	targetKey = "TARGET"
	columnKey = "COLUMNS"
)

var (
	input = os.Getenv(inputKey)
	output = os.Getenv(outputKey)
	source = os.Getenv(sourceKey)
	target = os.Getenv(targetKey)
	columns = os.Getenv(columnKey)
)

func init() {
	for _, k := range []string{
		inputKey,
		outputKey,
		sourceKey,
		targetKey,
		columnKey,
	} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}

func main() {
	var col []int
	for _, k := range strings.Split(columns, ","){
		i, err := strconv.Atoi(k)
		if err != nil {
			log.Fatal(err)
		}
		col = append(col, i)
	}
	beam.Init()
	p, s := beam.NewPipelineWithRoot()
	lines := textio.Read(s, fmt.Sprintf("%s/%s", input, source))
	tokens := simulator.Tokenize(s, lines, col...)
	filtered := simulator.Filter(s, tokens, "_", "-")
	textio.Write(s, fmt.Sprintf("%s/%s", output, target), filtered)
	err := beamx.Run(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}

