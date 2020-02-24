package io

import (
	"context"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/bigqueryio"
	"github.com/apache/beam/sdks/go/pkg/beam/testing/passert"
	"github.com/apache/beam/sdks/go/pkg/beam/transforms/top"
	"github.com/apache/beam/sdks/go/pkg/beam/x/beamx"
	"log"
	"os"
	"strings"
	"testing"
)

const (
	projectKey = "PROJECT"
)

var (
	project = os.Getenv(projectKey)
)

func init() {
	for _, k := range []string{
		projectKey,
	} {
		if os.Getenv(k) == "" {
			log.Fatalf("%s empty but expected from environment variables", k)
		}
	}
}

func TestBqReader_Read(t *testing.T) {
	st := station{}
	p, s := beam.NewPipelineWithRoot()
	r, err := NewReader("bigquery://bigquery-public-data:noaa_gsod.stations", WithProject(project), WithType(st))
	if err != nil {
		t.Error(err)
	}
	data := r.Read(s)
	largest := top.Largest(s, data, 1, func(a, b station) bool {
		return strings.Compare(a.Name, b.Name) < 0
	})
	passert.Equals(s, largest, []station{
		{
			Name: "ZYRYANKA",
		},
	})
	err = beamx.Run(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBqWriter_Write(t *testing.T) {
	p, s := beam.NewPipelineWithRoot()
	data := beam.Create(s, station{
		Name: "foo",
	})
	rawurl := fmt.Sprintf("bigquery://%s:beam.foo", project)
	w, err := NewWriter(rawurl, WithProject(project), WithType(station{}))
	if err != nil {
		t.Error(err)
	}
	w.Write(s, data)
	err = beamx.Run(context.Background(), p)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_bqReader_q(t *testing.T) {
	type foo struct {
		Foo bool
		Bar string
		Baz int
	}
	type fields struct {
		bqBase bqBase
	}
	tests := []struct {
		name       string
		fields     fields
		wantResult string
		wantErr    bool
	}{
		{
			name: "single field",
			fields: fields{
				bqBase: bqBase{
					qualifiedTableName: bigqueryio.QualifiedTableName{
						Project: "otherproject",
						Dataset: "dataset",
						Table: "table",
					},
					opts: map[string]interface{}{
							optionType: station{},
							optionProject: "project",
					},
				},
			},
			wantResult: "select Name from [otherproject:dataset.table] where Name is not null",
			wantErr: false,
		},
		{
			name: "multiple fields",
			fields: fields{
				bqBase: bqBase{
					qualifiedTableName: bigqueryio.QualifiedTableName{
						Project: "otherproject",
						Dataset: "dataset",
						Table: "table",
					},
					opts: map[string]interface{}{
						optionType: foo{},
						optionProject: "project",
					},
				},
			},
			wantResult: "select Foo,Bar,Baz from [otherproject:dataset.table] where Foo is not null or Bar is not null or Baz is not null",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bq := &bqReader{
				bqBase: tt.fields.bqBase,
			}
			gotResult, err := bq.q()
			if (err != nil) != tt.wantErr {
				t.Errorf("q() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResult != tt.wantResult {
				t.Errorf("q() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

