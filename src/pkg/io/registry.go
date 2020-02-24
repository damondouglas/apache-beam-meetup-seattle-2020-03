package io

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	gcs = "gs"

	bigQuery = "bigquery"

	local = "file"
)

var (
	readRegistry = map[string]Reader {
		gcs: &gcsReader{},
		bigQuery: &bqReader{},
	}
	writeRegistry = map[string]Writer {
		gcs: &gcsWriter{},
		bigQuery: &bqWriter{},
	}
)

func parseScheme(rawurl string) (result string, err error) {
	schemeMatcher := regexp.MustCompile("^.*://")
	specialCharMatcher := regexp.MustCompile("://$")
	find := schemeMatcher.FindString(rawurl)
	if find == "" {
		err = fmt.Errorf("scheme missing in %s", rawurl)
		return
	}
	result = specialCharMatcher.ReplaceAllString(find, "")
	result = strings.ToLower(result)
	return
}

func resolveReader(rawurl string, opts ...Option) (result Reader, err error) {
	scheme, err := parseScheme(rawurl)
	if err != nil {
		return
	}
	result, ok := readRegistry[scheme]
	if !ok {
		err = fmt.Errorf("%s does not map to a registered reader", rawurl)
		return
	}
	err = result.setOptions(opts...)
	if err != nil {
		return
	}
	err = result.setURL(rawurl)
	return
}

func resolveWriter(rawurl string, opts ...Option) (result Writer, err error) {
	scheme, err := parseScheme(rawurl)
	if err != nil {
		return
	}
	result, ok := writeRegistry[scheme]
	if !ok {
		err = fmt.Errorf("%s does not map to a registered writer", rawurl)
		return
	}
	err = result.setOptions(opts...)
	if err != nil {
		return
	}
	err = result.setURL(rawurl)
	return
}
