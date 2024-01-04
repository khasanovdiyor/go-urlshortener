package urlshortener

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	var handler http.Handler
	for _, urlToRedirect := range pathsToUrls {
		handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, ok := pathsToUrls[r.URL.Path]; ok {
				http.Redirect(w, r, urlToRedirect, http.StatusTemporaryRedirect)
			} else {
				fallback.ServeHTTP(w, r)
			}
		})
	}
	return handler.ServeHTTP
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

// JSONHandler will parse the provided JSON and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the JSON, then the
// fallback http.Handler will be called instead.
//
// JSON is expected to be in the format:
//
// [
//
//	{
//		"path": "/some-path",
//		"url": "https://www.some-url.com/demo"
//	}
//
// ]
//
// The only errors that can be returned all related to having
// invalid JSON data.
func JSONHandler(jsonToParse []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJson, err := parseJSON(jsonToParse)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	fmt.Println(pathMap["short"])
	return MapHandler(pathMap, fallback), nil
}

type redirect struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func parseYAML(yamlToParse []byte) ([]redirect, error) {
	return parseData(yamlToParse, yaml.Unmarshal)
}

func parseJSON(jsonToParse []byte) ([]redirect, error) {
	return parseData(jsonToParse, json.Unmarshal)
}

func parseData(dataToParse []byte, parser func([]byte, any) error) ([]redirect, error) {
	var pathUrls []redirect
	err := parser(dataToParse, &pathUrls)

	if err != nil {
		return nil, err
	}

	return pathUrls, nil
}

func buildMap(pathUrls []redirect) map[string]string {
	pathMap := map[string]string{}

	for _, pathUrl := range pathUrls {
		pathMap[pathUrl.Path] = pathUrl.Url
	}

	return pathMap
}
