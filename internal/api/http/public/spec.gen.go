// Package public provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package public

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9SVQW8bLRCG/wqa7zuuYqe97S2tpWhvUeueIisiy9gm2gUKs2kti/9eDThx1ruR46pV",
	"3NNiDMy8z8wLW6ht66xBQwHKLYR6ja1Mw7mdWf46bx160phmFYbaa0faGv5JG4dQQiCvzQpiAarDOyUJ",
	"+c+l9a0kKCFNFMPFpKnB0WPIKnunVe+UrtMKCrApuGygXMom4ODYWIDH7532qKC83cVYPC+z9w9YE8dg",
	"gZVxHb2TytMTvZFUr49muwdEvsPitOyPbn5Wc2RlHAjhKW2WdpAyXN1UYmm9aKWRK21WghULTdgGBipX",
	"Acrbw01zq6zwGGzna87cyJajcedAXDwnmuiJq5sKCnhEH/Lmy4spi7EOjXQaSvh4Mb2YQgFO0jpBnfBB",
	"abRCynrRS45dKSjhGumqaeZpDdcxOGtCLseH6ZQ/tTWEJm2VzjW6TpsnDyFXKTuNR1lnuYX/PS6hhP8m",
	"e09OdoacJDfumUrv5SYjPUApGh1I2OVLhrwsdG0r/QZK+ILkNT6ikE0zSjozXMQCnA0j2j97lISsHXIL",
	"Y6BPVm1OUn1MbHZm7LuEOy0OcF/+0cBjVOdrFHVSrfbEDrBmKkIKgz9eLBpQjcWutyZb/lQqZkM0SDhk",
	"PUvzO9ZOetkioR/zQzXjqtMaU1xBVjwVGth3UKbe7hulUnCIt3iB6vDqPbT4olcK0zVNv89a+7jL5n4j",
	"qtlYh71mrn9C8ele/93me63prpGOAHZPb0YfcXpKzhnyX7xV8jP6pqtlehZXS0pYSCPwpw7UeyBHS96N",
	"eOqbU/K8L5J3f0fOo9hfu/tAmjp+S95W8RjjrwAAAP//LaWrysgLAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
