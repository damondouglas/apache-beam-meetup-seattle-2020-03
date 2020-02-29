package main

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"log"
	"os"
	"strings"
	"temp/pkg/io"
	"time"
)

const (
	patientsKey = "PATIENTS"
	outputKey  = "OUTPUT"
	projectKey = "PROJECT"
	snomedKey  = "SNOMED"
)

var (
	patients = os.Getenv(patientsKey)
	output  = os.Getenv(outputKey)
	project = os.Getenv(projectKey)
	snomed  = os.Getenv(snomedKey)
)

func init() {
	for _, k := range []string{
		patientsKey,
		outputKey,
		projectKey,
		snomedKey,
	} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}

type Term string
type ConceptID string
type MRN string

type patient struct {
	MRN         MRN
	RawAllergies string `bigquery:"allergies"`
}

type drug struct {
	Term      Term
	ConceptID ConceptID
}

type coded struct {
	MRN      MRN
	ConceptID ConceptID
}

func main() {
	var err error
	p, s := beam.NewPipelineWithRoot()

	patients, err := readPatients(s)
	if err != nil {
		log.Fatal(err)
	}

	drugs, err := readDrugs(s)
	if err != nil {
		log.Fatal(err)
	}

	splitAllergies := beam.ParDo(s, splitFn, patients)

	keyedDrugs := beam.ParDo(s, keyDrugFn, drugs)

	joined := beam.CoGroupByKey(s, splitAllergies, keyedDrugs)

	coded := beam.ParDo(s, func(t Term, allergies func(*MRN) bool, drugs func(*ConceptID) bool, emit func(coded)) {
		var mrn MRN
		var conceptID ConceptID
		hasMRN := allergies(&mrn)
		hasDrug := drugs(&conceptID)

		if hasMRN && hasDrug {
			emit(coded{
				MRN: mrn,
				ConceptID: conceptID,
			})
		}

	}, joined)

	err = writeCoded(s, coded)
	if err != nil {
		log.Fatal(err)
	}

	err = beamx.Run(context.Background(), p)
	if err != nil {
		log.Fatal(err)
	}
}

func keyDrugFn(d drug) (Term, ConceptID) {
	term := string(d.Term)
	term = strings.ToLower(term)
	return Term(term), d.ConceptID
}

func splitFn(p patient, emit func(Term, MRN)) {
	allergies := strings.Split(p.RawAllergies, ",")
	for _, k := range allergies {
		k = strings.ToLower(k)
		emit(Term(k), p.MRN)
	}
}

func readPatients(s beam.Scope) (result beam.PCollection, err error) {
	r, err := io.NewReader(patients, io.WithProject(project), io.WithType(patient{}))
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

func writeCoded(s beam.Scope, pCoded beam.PCollection) (err error) {
	table := fmt.Sprintf("%s_%v", output, time.Now().Unix())
	w, err := io.NewWriter(table, io.WithProject(project), io.WithType(coded{}))
	if err != nil {
		return
	}
	w.Write(s, pCoded)
	return
}