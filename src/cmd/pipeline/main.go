package main

import (
	"fmt"
	"log"
	"os"
)

const (
	inputKey = "INPUT"
	outputKey = "OUTPUT"
	sourceKey = "SOURCE"
	targetKey = "TARGET"
)

var (
	input = os.Getenv(inputKey)
	output = os.Getenv(outputKey)
	source = os.Getenv(sourceKey)
	target = os.Getenv(targetKey)
)

func init() {
	for _, k := range []string{
		inputKey,
		outputKey,
		sourceKey,
		targetKey,
	} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}
func main() {
	fmt.Println(input, output, source, target)
}