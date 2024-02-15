// Package private provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package private

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

	"H4sIAAAAAAAC/8xUS2/UPhD/Ktb8/wfQRk1oe8qtEi2thGDVwqnag+tMErPO2IzthaXKd0d29tFueani",
	"wM1xPL/XeHwPyg7OElLwUN8Do3eWPOaPS5Qm9OsLqU1kTDvKUkAKaSmdM1rJoC2Vn7yltOdVj4NMq/8Z",
	"W6jhv3IPX05/fXnObBnGcSygQa9YuwQCNVxvyEVrWUTqJ3rh2N4hjMVWz01UCr3/a3o2ePOJ5jeyDkSN",
	"xYYky5mcJTFNo1O1NHO2DjnolGgrjcdD9JtcLoIVDbaaUIQeBSYgsW0GFOAewCTbTW4HfpWDMwj1aXVa",
	"QGt5kAFq0BROjqGAsHY4fWKHnBIc0HvZPa6FdzaICxup2Zf4wJq6bI/xc9SMDdS3E+8eZTEW8Ci8+v7A",
	"3AfbWLHpQgFIcUgw75eweEo1FqCptU9BzielwrbibH4lGvS6I9Fq9kFI59hK1QtJjYheUyc6ayR1wgep",
	"luLFSjvkGarezjrLw0xZak1ECjNlUJJk1euAKkTGl8m/DjmTm4ky6z+bX0EBK2Q/6Xl1VKUsrUOSTkMN",
	"J0fV0Ulqkgx97k9p9Aq/pVWH4amhHFbus6aWpQ8cswCh/e5+JT+MITKJ46oSX3okIY0RknF7BrIGzlf+",
	"qoEa3mB4m4mLx3N8XFU/m4HdufJguPIUtDKa8Mel23ciD0UcBslrqCEpIvReLOMdMmFAPw2PQGqc1RRS",
	"7LLz6WZMzlSPagmLhFMyymb9vChz6S+DzCd+FOP1xPov5Zgk6ecGmZCQ0xWG+vYww4+bp+01rtBYN2BG",
	"imyghjI9HkzSwLjYoR8CXO7JxC5KDwWQHNIwPRQzLsbvAQAA//8qIdPfdAYAAA==",
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
