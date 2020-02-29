package io

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
)

type gcsBase struct {
	url string
	opts map[string]interface{}
}

type gcsReader struct {
	gcsBase
}

type gcsWriter struct {
	gcsBase
}

func (g *gcsBase) setURL(rawurl string) (err error) {
	g.url = rawurl
	return
}

func (g *gcsBase) setOptions(opts ...Option) (err error) {
	if g.opts == nil {
		g.opts = map[string]interface{}{}
	}
	for _, k := range opts {
		g.opts[k.key()] = k.value()
	}
	return
}

func (g *gcsReader) Read(s beam.Scope) (result beam.PCollection) {
	result = textio.Read(s, g.url)
	return
}

func (g *gcsWriter) Write(s beam.Scope, collection beam.PCollection) {
	textio.Write(s, g.url, collection)
}

