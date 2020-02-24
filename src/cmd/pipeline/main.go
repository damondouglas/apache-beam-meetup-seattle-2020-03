package main

import (
	"context"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"log"
	"os"
	"strings"
	"temp/pkg/io"
)

const (
	inputKey = "INPUT"
	outputKey = "OUTPUT"
	projectKey = "PROJECT"
	rxnormKey = "RXNORM"
)

var (
	input = os.Getenv(inputKey)
	output = os.Getenv(outputKey)
	project = os.Getenv(projectKey)
	rxnorm = os.Getenv(rxnormKey)
)

func init() {
	for _, k := range []string{
		inputKey,
		outputKey,
		projectKey,
		rxnormKey,
	} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}
func main() {
	p, s := beam.NewPipelineWithRoot()
	var patients, drugs beam.PCollection
	//patients, err := readPatients(s)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//drugs, err := readDrugs(s)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//result := beam.ParDo(s, codeFn, patients, beam.SideInput{Input: drugs})
	//beam.ParDo(s, resultFn, result)

	err := beamx.Run(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}

func handler(primary string, side string, emit func(string)) {

}

type patient struct {
	MRN string
	Allergies string
}

type drug struct {
	Name string
	Rxcui string
}

type coded struct {
	MRN string
	Rxcui string
	Match float64
}

func readPatients(s beam.Scope) (result beam.PCollection, err error) {
	r, err := io.NewReader(input, io.WithProject(project), io.WithType(patient{}))
	if err != nil {
		return
	}
	result = r.Read(s)
	return
}

func readDrugs(s beam.Scope) (result beam.PCollection, err error) {
	r, err := io.NewReader(rxnorm, io.WithProject(project), io.WithType(drug{}))
	if err != nil {
		return
	}
	result = r.Read(s)
	return
}

func codeFn(p patient, d drug, emit func(coded)) {
	allergies := strings.Split(p.Allergies, ",")
	for _, k := range allergies {
		k = strings.TrimSpace(k)
		k = strings.ToLower(k)
		m := d.Name
		m = strings.TrimSpace(m)
		m = strings.ToLower(m)
		if k == m {
			emit(coded{
				MRN: p.MRN,
				Rxcui: d.Rxcui,
				Match: 1.0,
			})
		}
	}
}

func resultFn(c coded, emit func(coded)) {
	log.Println(c)
}