package url

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
	"github.com/apache/beam/sdks/go/pkg/beam/io/textio"
)

type gcsBase struct {
	url string
}

type gcsReader struct {
	gcsBase
}

type gcsWriter struct {
	gcsBase
}

func (g *gcsBase) setURL(rawurl string) {
	g.url = rawurl
}

func (g *gcsReader) Read(s beam.Scope) (result beam.PCollection) {
	result = textio.Read(s, g.url)
	return
}

func (g *gcsWriter) Write(s beam.Scope, collection beam.PCollection) {
	textio.Write(s, g.url, collection)
}

