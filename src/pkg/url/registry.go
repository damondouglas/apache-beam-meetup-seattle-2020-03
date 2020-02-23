package url

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
	}
	writeRegistry = map[string]Writer {
		gcs: &gcsWriter{},
	}
)

func parseScheme(rawurl string) (result string, err error) {
	schemeMatcher := regexp.MustCompile("^.*:\\/\\/")
	specialCharMatcher := regexp.MustCompile(":\\/\\/$")
	find := schemeMatcher.FindString(rawurl)
	if find == "" {
		err = fmt.Errorf("scheme missing in %s", rawurl)
		return
	}
	result = specialCharMatcher.ReplaceAllString(find, "")
	result = strings.ToLower(result)
	return
}

func resolveReader(rawurl string) (result Reader, err error) {
	scheme, err := parseScheme(rawurl)
	if err != nil {
		return
	}
	result, ok := readRegistry[scheme]
	if !ok {
		err = fmt.Errorf("%s does not map to a registered reader", rawurl)
		return
	}
	result.setURL(rawurl)
	return
}

func resolveWriter(rawurl string) (result Writer, err error) {
	scheme, err := parseScheme(rawurl)
	if err != nil {
		return
	}
	result, ok := writeRegistry[scheme]
	if !ok {
		err = fmt.Errorf("%s does not map to a registered writer", rawurl)
		return
	}
	result.setURL(rawurl)
	return
}
