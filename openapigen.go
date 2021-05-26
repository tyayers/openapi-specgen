// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openapigen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// Generates a simple OpenAPI spec based on a URL get method call
func GenerateSpec(url string) string {
	urlPieces := strings.Split(url, "/")
	resourceName := urlPieces[len(urlPieces)-1]
	baseUrl := strings.Replace(url, "/"+resourceName, "", -1)

	singularResourceName := resourceName
	if last := len(singularResourceName) - 1; last >= 0 && singularResourceName[last] == 's' {
		singularResourceName = singularResourceName[:last]
	}
	singularResourceNameCapitalized := strings.Title(singularResourceName)

	resultSpec := initSpec(resourceName, singularResourceName, singularResourceNameCapitalized, baseUrl)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)

		var result map[string][]map[string]string
		var resourceArray []map[string]string

		json.Unmarshal([]byte(data), &result)

		if len(result[resourceName]) > 0 {
			resourceArray = result[resourceName]
		} else {
			json.Unmarshal([]byte(data), &resourceArray)
		}

		// We have data, now go through properties of first object
		for key, value := range resourceArray[0] {
			//fmt.Println(key, value)

			propType := "string"

			if _, err := strconv.ParseInt(value, 10, 64); err == nil {
				propType = "number"
			}

			resultSpec.Components.Schemas[singularResourceNameCapitalized].Properties[key] = Schema{
				Description: "The " + key + " of the " + singularResourceNameCapitalized,
				Type:        propType,
				Example:     value,
			}
		}
	}

	d, _ := yaml.Marshal(&resultSpec)
	//fmt.Printf("--- t dump:\n%s\n\n", string(d))
	return string(d)
}

func initSpec(resourceName string, singularResourceName string, singularResourceNameCapitalized string, baseUrl string) OpenApiSchema {

	resultSpec := OpenApiSchema{
		Openapi: "3.0.3",
		Info: SpecInfo{
			Title:       strings.Title(resourceName + " API"),
			Description: "API for managing " + singularResourceNameCapitalized + " resources.",
			Version:     "0.0.1",
		},
		Servers: []Server{
			{
				Url: baseUrl,
			},
		},
		Paths: map[string]OpsCollection{
			"/" + resourceName: {
				Get: Operation{
					Summary:     "List '" + singularResourceNameCapitalized + "' objects.",
					Description: "Retrieve a page of '" + singularResourceNameCapitalized + "' objects from the server.  Follows the standards for parameters from the [List AIP](  https://aip.dev/132).",
					Parameters: []Parameter{
						{
							Name:        "pageSize",
							In:          "query",
							Description: "Max size of returned list.",
							Schema: Schema{
								Type:    "integer",
								Default: "25",
							},
						},
						{
							Name:        "pageToken",
							In:          "query",
							Description: "A page token recieved from the previous list call. Provide this to retrieve the next page.",
							Schema: Schema{
								Type: "string",
							},
						},
						{
							Name:        "orderBy",
							In:          "query",
							Description: "The ordering of the returned list. See the [List Ordering API]( https://aip.dev/132) for details on the formatting of this field.",
							Schema: Schema{
								Type:    "string",
								Default: "displayName",
							},
						},
						{
							Name:        "filter",
							In:          "query",
							Description: "Filter that will be used to select " + singularResourceNameCapitalized + " objects to return. See the [Filtering AIP](https://aip.dev/160) for usage and details on the filtering grammar.",
							Schema: Schema{
								Type: "string",
							},
						},
					},
					Responses: map[string]Response{
						"200": {
							Description: "Successful response",
							Content: Content{
								ApplicationJson: ContentSchema{
									Schema: Schema{
										Type: "object",
										Properties: map[string]Schema{
											resourceName: {
												Type: "array",
												Items: Item{
													Ref: "#/components/schemas/ListOf" + strings.Title(resourceName),
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Post: Operation{
					Summary:     "Creates a new '" + singularResourceNameCapitalized + "' object.",
					Description: "Creates a new '" + singularResourceNameCapitalized + "' object.",
					RequestBody: Body{
						Description: "The " + singularResourceNameCapitalized + " object to create.",
						Content: Content{
							ApplicationJson: ContentSchema{
								Schema: Schema{
									Ref: "#/components/schemas/" + singularResourceNameCapitalized,
								},
							},
						},
					},
					Responses: map[string]Response{
						"201": {
							Description: "Successful response",
							Content: Content{
								ApplicationJson: ContentSchema{
									Schema: Schema{
										Ref: "#/components/schemas/" + singularResourceNameCapitalized,
									},
								},
							},
						},
					},
				},
			},
			"/" + resourceName + "/{" + singularResourceName + "}": {
				Get: Operation{
					Summary:     "Retrieve " + singularResourceNameCapitalized + " object.",
					Description: "Retrieve a single " + singularResourceNameCapitalized + " object.",
					Parameters: []Parameter{
						{
							Name:        singularResourceName,
							Description: "Unique identifier of the desired " + singularResourceNameCapitalized + " object.",
							In:          "path",
							Required:    true,
							Schema: Schema{
								Type: "string",
							},
						},
					},
					Responses: map[string]Response{
						"200": {
							Description: "Successful response",
							Content: Content{
								ApplicationJson: ContentSchema{
									Schema: Schema{
										Ref: "#/components/schemas/" + singularResourceNameCapitalized,
									},
								},
							},
						},
						"404": {
							Description: singularResourceNameCapitalized + " was not found.",
						},
					},
				},
				Put: Operation{
					Summary:     "Update " + singularResourceNameCapitalized + " object.",
					Description: "Update a single " + singularResourceNameCapitalized + " object.",
					Parameters: []Parameter{
						{
							Name:        singularResourceName,
							Description: "Unique identifier of the desired " + singularResourceNameCapitalized + " object.",
							In:          "path",
							Required:    true,
							Schema: Schema{
								Type: "string",
							},
						},
					},
					RequestBody: Body{
						Description: "The " + singularResourceNameCapitalized + " object to update.",
						Content: Content{
							ApplicationJson: ContentSchema{
								Schema: Schema{
									Ref: "#/components/schemas/" + singularResourceNameCapitalized,
								},
							},
						},
					},
					Responses: map[string]Response{
						"200": {
							Description: "Successful response",
							Content: Content{
								ApplicationJson: ContentSchema{
									Schema: Schema{
										Ref: "#/components/schemas/" + singularResourceNameCapitalized,
									},
								},
							},
						},
						"404": {
							Description: singularResourceNameCapitalized + " was not found.",
						},
					},
				},
				Delete: Operation{
					Summary:     "Delete " + singularResourceNameCapitalized + " object.",
					Description: "Delete a single " + singularResourceNameCapitalized + " object.",
					Parameters: []Parameter{
						{
							Name:        singularResourceName,
							Description: "Unique identifier of the desired " + singularResourceNameCapitalized + " object.",
							In:          "path",
							Required:    true,
							Schema: Schema{
								Type: "string",
							},
						},
					},
					Responses: map[string]Response{
						"200": {
							Description: "Successful response",
						},
						"404": {
							Description: singularResourceNameCapitalized + " was not found.",
						},
					},
				},
			},
		},
		Components: Components{
			Schemas: map[string]Schema{
				"ListOf" + strings.Title(resourceName): {
					Title: "List of " + singularResourceNameCapitalized + " objects",
					Type:  "array",
					Items: Item{
						Ref: "#/components/schemas/" + singularResourceNameCapitalized,
					},
				},
				singularResourceNameCapitalized: {
					Title:      singularResourceNameCapitalized,
					Type:       "object",
					Properties: map[string]Schema{},
				},
			},
		},
	}

	return resultSpec
}

type OpenApiSchema struct {
	Openapi    string                   `yaml:"openapi"`
	Info       SpecInfo                 `yaml:"info"`
	Servers    []Server                 `yaml:"servers"`
	Paths      map[string]OpsCollection `yaml:"paths"`
	Components Components               `yaml:"components"`
}

type SpecInfo struct {
	Description string `yaml:"description"`
	Version     string `yaml:"version"`
	Title       string `yaml:"title"`
}

type Server struct {
	Url string `yaml:"url"`
}

type OpsCollection struct {
	Get    Operation `yaml:"get"`
	Post   Operation `yaml:"post,omitempty"`
	Put    Operation `yaml:"put,omitempty"`
	Delete Operation `yaml:"delete,omitempty"`
}

type Operation struct {
	Summary     string              `yaml:"summary"`
	Description string              `yaml:"description"`
	Parameters  []Parameter         `yaml:"parameters,omitempty"`
	RequestBody Body                `yaml:"requestBody,omitempty"`
	Responses   map[string]Response `yaml:"responses"`
}

type Parameter struct {
	Name        string `yaml:"name"`
	In          string `yaml:"in"`
	Required    bool   `yaml:"required,omitempty"`
	Description string `yaml:"description"`
	Schema      Schema `yaml:"schema"`
}

type Body struct {
	Description string  `yaml:"description"`
	Required    bool    `yaml:"required"`
	Content     Content `yaml:"content"`
}

type Response struct {
	Description string  `yaml:"description"`
	Content     Content `yaml:"content,omitempty"`
}

type Content struct {
	ApplicationJson ContentSchema `yaml:"application/json"`
}

type ContentSchema struct {
	Schema Schema `yaml:"schema"`
}

type Schema struct {
	Title       string            `yaml:"title,omitempty"`
	Description string            `yaml:"description,omitempty"`
	Type        string            `yaml:"type,omitempty"`
	Default     string            `yaml:"default,omitempty"`
	Items       Item              `yaml:"items,omitempty"`
	Ref         string            `yaml:"$ref,omitempty"`
	Properties  map[string]Schema `yaml:"properties,omitempty"`
	Example     string            `yaml:"example,omitempty"`
}

type Item struct {
	Ref string `yaml:"$ref"`
}

type Components struct {
	Schemas map[string]Schema `yaml:"schemas"`
}
