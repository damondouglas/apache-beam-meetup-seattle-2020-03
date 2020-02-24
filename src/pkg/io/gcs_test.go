package io

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
	"github.com/apache/beam/sdks/go/pkg/beam/testing/passert"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/top"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	gcsOutKey = "GCS_OUTPUT"
)

var (
	gcsOut = os.Getenv(gcsOutKey)
)

func init() {
	for _, k := range []string{
		gcsOutKey,
	} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}

func Test_gcsReader_Read(t *testing.T) {
	p, s := beam.NewPipelineWithRoot()
	r, err := NewReader("gs://apache-beam-samples/shakespeare/kinglear.txt")
	if err != nil {
		t.Fatal(err)
	}
	lines := r.Read(s)
	largest := top.Largest(s, lines, 1, func(a, b string) bool {
		return strings.Compare(a, b) < 0
	})
	passert.True(s, largest, func(element []string) bool {
		return strings.Contains(element[0], "Third Servant")
	})
	err = beamx.Run(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_gcsWriter_Write(t *testing.T) {
	p, s := beam.NewPipelineWithRoot()
	rawurl := fmt.Sprintf("%s/%s", gcsOut, "/test_gcsWriter_Write")
	w, err := NewWriter(rawurl)
	if err != nil {
		t.Fatal(err)
	}
	col := beam.Create(s, "foo", "bar")
	w.Write(s, col)

	err = beamx.Run(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
	p, s = beam.NewPipelineWithRoot()
	lines := textio.Read(s, rawurl)
	passert.Equals(s, lines, "foo", "bar")
	err = beamx.Run(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
}

type station struct {
	Name string
}

