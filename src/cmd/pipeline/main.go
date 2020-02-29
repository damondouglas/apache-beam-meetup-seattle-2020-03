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
	inputKey   = "INPUT"
	outputKey  = "OUTPUT"
	projectKey = "PROJECT"
	snomedKey  = "SNOMED"
)

var (
	input   = os.Getenv(inputKey)
	output  = os.Getenv(outputKey)
	project = os.Getenv(projectKey)
	snomed  = os.Getenv(snomedKey)
)

func init() {
	for _, k := range []string{
		inputKey,
		outputKey,
		projectKey,
		snomedKey,
	} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}

type patient struct {
	MRN          string
	RawAllergies string `bigquery:"allergies"`
}

type Term string

type splitPatientAllergy struct {
	Term Term
	MRN string
}

type drug struct {
	Term      Term
	ConceptID int
}

type coded struct {
	MRN       string
	ConceptID int
	Match     float64
}

func main() {
	var err error
	p, s := beam.NewPipelineWithRoot()
	patients := beam.Create(s,
		patient{
			MRN:          "M12345",
			RawAllergies: "qwert,asdf,zxcv",
		},
		patient{
			MRN:          "M54321",
			RawAllergies: "yuio,fdsa,qree",
		},
	)
	drugs := beam.Create(
		s,
		drug{
			Term:      "qwert",
			ConceptID: 1,
		},
		drug{
			Term:      "asdf",
			ConceptID: 2,
		},
		drug{
			Term:      "ydydu",
			ConceptID: 3,
		},
		drug{
			Term:      "gree",
			ConceptID: 3,
		},
	)

	splitAllergies := beam.ParDo(s, splitFn, patients)

	keyedDrugs := beam.ParDo(s, keyDrugFn, drugs)

	joined := beam.CoGroupByKey(s, splitAllergies, keyedDrugs)

	beam.ParDo(s, func(t Term, allergies func(*string) bool, drugs func(*int) bool, emit func(Term, string)) {
		log.Println(t)
	}, joined)

	err = beamx.Run(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}

func keyDrugFn(d drug) (Term, int) {
	return d.Term, d.ConceptID
}

func splitFn(p patient, emit func(Term, string)) {
	allergies := strings.Split(p.RawAllergies, ",")
	for _, k := range allergies {
		emit(Term(k), p.MRN)
	}
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
	r, err := io.NewReader(snomed, io.WithProject(project), io.WithType(drug{}))
	if err != nil {
		return
	}
	result = r.Read(s)
	return
}
