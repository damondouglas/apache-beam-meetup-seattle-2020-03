// Package url handles parsing of URLs into the respective beam IO workloads
package url

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
)

type Reader interface {
	urlSetter
	Read(s beam.Scope) beam.PCollection
}

type Writer interface {
	urlSetter
	Write(s beam.Scope, collection beam.PCollection)
}

type urlSetter interface {
	setURL(rawurl string)
}

func NewReader(rawurl string) (result Reader, err error) {
	return resolveReader(rawurl)
}

func NewWriter(rawurl string) (result Writer, err error) {
	return resolveWriter(rawurl)
}

