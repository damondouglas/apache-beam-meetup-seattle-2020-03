// Package url handles parsing of URLs into the respective beam IO workloads
package io

import (
	"github.com/apache/beam/sdks/go/pkg/beam"
)

type Reader interface {
	urlSetter
	optionsSetter
	Read(s beam.Scope) beam.PCollection
}

type Writer interface {
	urlSetter
	optionsSetter
	Write(s beam.Scope, collection beam.PCollection)
}

type urlSetter interface {
	setURL(rawurl string) error
}

type optionsSetter interface {
	setOptions(opts ...Option) error
}

func NewReader(rawurl string, opt ...Option) (result Reader, err error) {
	return resolveReader(rawurl, opt...)
}

func NewWriter(rawurl string, opt ...Option) (result Writer, err error) {
	return resolveWriter(rawurl, opt...)
}

