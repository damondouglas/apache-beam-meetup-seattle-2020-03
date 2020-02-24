package io

import (
	"bytes"
	"cloud.google.com/go/bigquery"
	"fmt"
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/bigqueryio"
	"reflect"
	"regexp"
	"strings"
	"text/template"
)

type bqBase struct {
	qualifiedTableName bigqueryio.QualifiedTableName
	opts map[string]interface{}
}
type bqReader struct {
	query string
	bqBase
}
type bqWriter struct {
	bqBase
}

// setURL parses a bigquery url to its project, dataset and table
// rawurl is expected to be:
// bigquery://project:dataset.table
func (bq *bqBase) setURL(rawurl string) (err error) {
	schemeMatcher := regexp.MustCompile("^.*://")
	rawurl = schemeMatcher.ReplaceAllString(rawurl, "")

	bq.qualifiedTableName, err = bigqueryio.NewQualifiedTableName(rawurl)
	return
}

func (bq *bqBase) setOptions(opts ...Option) (err error) {
	if bq.opts == nil {
		bq.opts = map[string]interface{}{}
	}
	for _, k := range opts {
		bq.opts[k.key()] = k.value()
	}
	for k, v := range map[string]string {
		optionType:    "io.WithType",
		optionProject: "io.WithProject",
	} {
		if _, ok := bq.opts[k]; !ok {
			err = fmt.Errorf("io.Option %s missing but expected: use %s to set this", k, v)
		}
	}
	return
}

func (bq *bqBase) project() string {
	return bq.opts[optionProject].(string)
}

func (bq *bqReader) setURL(rawurl string) (err error) {
	err = bq.bqBase.setURL(rawurl)
	if err != nil {
		return
	}
	bq.query, err = bq.q()
	return
}

func (bq *bqReader) t() reflect.Type {
	return reflect.TypeOf(bq.opts[optionType])
}

func (bq *bqReader) Interface() interface{} {
	return bq.opts[optionType]
}

func (bq *bqReader) table() (result string) {
	result = bq.qualifiedTableName.String()
	return
}

func (bq *bqReader) schema() (bigquery.Schema, error) {
	return bigquery.InferSchema(bq.Interface())
}

type queryData struct {
	FieldList string
	Table string
	Where string
}

func (bq *bqReader) q() (result string, err error) {
	schema, err := bq.schema()
	if err != nil {
		return
	}

	fields := make([]string, len(schema))
	clauses := make([]string, len(schema))
	for i := range fields {
		name := schema[i].Name
		fields[i] = name
		clauses[i] = fmt.Sprintf("%s is not null", name)
	}

	tmpl, err := template.New("q").Parse("select {{.FieldList}} from [{{.Table}}] where {{.Where}}")
	if err != nil {
		return
	}
	buf := bytes.Buffer{}
	data := &queryData{
		FieldList: strings.Join(fields, ","),
		Table: bq.table(),
		Where: strings.Join(clauses, " or "),
	}
	err = tmpl.Execute(&buf, data)

	if err != nil {
		return
	}

	result = buf.String()
	return
}

func (bq *bqReader) Read(s beam.Scope) (result beam.PCollection) {
	result = bigqueryio.Query(s, bq.project(), bq.query, bq.t())
	//result = bigqueryio.Read(s, bq.project(), bq.qualifiedTableName.String(), bq.t())
	return
}

func (bq *bqWriter) Write(s beam.Scope, collection beam.PCollection) {
	bigqueryio.Write(s, bq.project(), bq.qualifiedTableName.String(), collection)
	return
}