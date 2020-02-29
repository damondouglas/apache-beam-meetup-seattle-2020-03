package main

import (
	"context"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"log"
	"os"
	"strconv"
	"strings"
	"temp/pkg/io"
	"temp/pkg/simulator"
)

const (
	inputKey = "INPUT"
	outputKey = "OUTPUT"
	columnKey = "COLUMNS"
	projectKey = "PROJECT"
)

var (
	input = os.Getenv(inputKey)
	output = os.Getenv(outputKey)
	columns = os.Getenv(columnKey)
	project = os.Getenv(projectKey)
)

func init() {
	for _, k := range []string{
		inputKey,
		outputKey,
		columnKey,
		projectKey,
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
	reader, err := io.NewReader(input)
	if err != nil {
		log.Fatal(err)
	}
	writer, err := io.NewWriter(output, io.WithProject(project), io.WithType(simulator.Drug{}))
	if err != nil {
		log.Fatal(err)
	}
	lines := reader.Read(s)
	tokens := simulator.Tokenize(s, lines, col...)
	filtered := simulator.Filter(s, tokens, "_", "-")
	writer.Write(s, filtered)
	err = beamx.Run(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}

