package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SwaggerMessageRef struct {
	FullName string `json:"full_name"`
	Location string `json:"location"`
}

type SwaggerMessage struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Required    []string               `json:"required"`
	Properties  []*SwaggerMessageField `json:"properties"`
	Identify    string                 `json:"identify"`
	Description string                 `json:"description"`
}

func (swaggerMessage *SwaggerMessage) AsRequest(location string) interface{} {
	if location == "param" {
		elements := make([]interface{}, len(swaggerMessage.Properties))
		for i, element := range swaggerMessage.Properties {
			elements[i] = element.AsRequest(location)
		}
		return elements
	}
	return []map[string]interface{}{
		{
			"in":          "body",
			"name":        swaggerMessage.Name,
			"description": swaggerMessage.Description,
			"required":    true,
			"schema": map[string]string{
				"$ref": fmt.Sprintf("#/definitions/%s", swaggerMessage.Identify),
			},
		},
	}
}

func (swaggerMessage *SwaggerMessage) AsResponse() interface{} {
	return map[string]interface{}{
		"schema": map[string]interface{}{
			"$ref": fmt.Sprintf("#/definitions/%s", swaggerMessage.Identify),
		},
		"description": swaggerMessage.Description,
	}
}

func (swaggerMessage *SwaggerMessage) AsDef() interface{} {
	def := map[string]interface{}{
		"type":        swaggerMessage.Type,
		"required":    swaggerMessage.Required,
		"description": swaggerMessage.Description,
	}
	pps := make(map[string]interface{})
	for _, field := range swaggerMessage.Properties {
		pps[field.Name] = field.AsDefinition()
	}
	def["properties"] = pps
	return def
}

type SwaggerMessageField struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Required    bool            `json:"required"`
	Type        string          `json:"-"`
	IsSlice     bool            `json:"-"`
	Enum        []string        `json:"enum,omitempty"`
	Example     string          `json:"example,omitempty"`
	Format      string          `json:"format,omitempty"`
	Identify    string          `json:"-"`
	Refer       *SwaggerMessage `json:"ref,omitempty"`
}

func (swaggerField *SwaggerMessageField) AsRequest(location string) interface{} {
	data := map[string]interface{}{
		"name":        swaggerField.Name,
		"description": swaggerField.Description,
		"required":    swaggerField.Required,
		"in":          location,
		"example":     swaggerField.Example,
		"format":      swaggerField.Format,
	}
	if len(swaggerField.Enum) > 0 {
		data["enum"] = swaggerField.Enum
	}
	if swaggerField.Refer != nil {
		data["schema"] = map[string]interface{}{
			"$ref": fmt.Sprintf("#/definitions/%s", swaggerField.Refer.Identify),
		}
	}
	return data
}

func (swaggerField *SwaggerMessageField) itemsType() map[string]string {
	if swaggerField.Refer != nil {
		return map[string]string{
			"$ref": fmt.Sprintf("#/definitions/%s", swaggerField.Refer.Identify),
		}
	}
	return map[string]string{
		"type": swaggerField.Type,
	}
}

func (swaggerField *SwaggerMessageField) AsDefinition() map[string]interface{} {
	def := map[string]interface{}{
		"description": swaggerField.Description,
		"type":        swaggerField.Type,
	}
	if swaggerField.IsSlice {
		def["type"] = "array"
		def["items"] = swaggerField.itemsType()
	}
	if swaggerField.Example != "" {
		def["example"] = swaggerField.Example
	}
	if swaggerField.Format != "" {
		def["format"] = swaggerField.Format
	}
	if len(swaggerField.Enum) > 0 {
		def["enum"] = swaggerField.Enum
	}
	if swaggerField.Refer != nil {
		def["$ref"] = fmt.Sprintf("#/definitions/%s", swaggerField.Refer.Identify)
	}
	return def
}

type SwaggerTag struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SwaggerAuth struct {
	Name    string   `json:"name"`
	Scope   []string `json:"scope"`
	Schemes []string `json:"schemes"`
}

type SwaggerApis struct {
	Tags       []*SwaggerTag `json:"tags"`
	BasePath   string        `json:"basePath"`
	Containers map[string]map[string]*SwaggerRoute
}

func (sa *SwaggerApis) convertTags() []string {
	ss := make([]string, len(sa.Tags))
	for i, tag := range sa.Tags {
		ss[i] = tag.Name
	}
	return ss
}

func (sa *SwaggerApis) AddRoute(path, method string, route *SwaggerRoute) {
	if sa.Containers == nil {
		sa.Containers = make(map[string]map[string]*SwaggerRoute)
	}
	fullPath := sa.BasePath + path
	v, ok := sa.Containers[fullPath]
	if !ok {
		v = make(map[string]*SwaggerRoute)
	}
	v[method] = route
	sa.Containers[fullPath] = v
}

func (sa *SwaggerApis) AsDoc(paramsMap map[string]*SwaggerMessage) map[string]interface{} {
	apis := make(map[string]interface{})
	for path, routes := range sa.Containers {
		target := make(map[string]interface{})
		for method, route := range routes {
			apiDoc := route.AsDoc(paramsMap)
			apiDoc["tags"] = sa.convertTags()
			target[strings.ToLower(method)] = apiDoc
		}
		apis[path] = target
	}
	return apis
}

type SwaggerDoc struct {
	Tags []*SwaggerTag           `json:"tags"`
	Auth []*SwaggerAuth          `json:"auth"`
	Host string                  `json:"host"`
	Apis map[string]*SwaggerApis `json:"apis"`
	Loc  string                  `json:"loc"`
}

func (sd *SwaggerDoc) AddService(name string, service *SwaggerApis) {
	if sd.Apis == nil {
		sd.Apis = make(map[string]*SwaggerApis)
	}
	_, ok := sd.Apis[name]
	if !ok {
		sd.Apis[name] = service
		sd.Tags = append(sd.Tags, service.Tags...)
	}
}

func (sd *SwaggerDoc) AsDoc(paramsMap map[string]*SwaggerMessage) map[string]interface{} {
	doc := map[string]interface{}{
		"swagger": "2.0",
		"host":    sd.Host,
		"info": map[string]interface{}{
			"version": strconv.FormatInt(time.Now().Unix(), 10),
			"title":   "Generated By Cato",
		},
		"schemes": []string{"http", "https"},
	}
	paths := make(map[string]interface{})
	for _, apis := range sd.Apis {
		for key, value := range apis.AsDoc(paramsMap) {
			paths[key] = value
		}
	}
	doc["paths"] = paths
	defs := make(map[string]interface{})
	for name, param := range paramsMap {
		defs[name] = param.AsDef()
	}
	doc["definitions"] = defs
	return doc
}

type SwaggerRoute struct {
	Description string   `json:"description"`
	OperationId string   `json:"operationId"`
	Consumes    []string `json:"consumes"`
	Produces    []string `json:"produces"`

	Requests  *SwaggerMessageRef `json:"requests"`
	Responses *SwaggerMessageRef `json:"responses"`
}

func (sr *SwaggerRoute) AsDoc(paramsMap map[string]*SwaggerMessage) map[string]interface{} {
	data := map[string]interface{}{
		"description": sr.Description,
		"summary":     sr.Description,
		"operationId": sr.OperationId,
		"produces":    sr.Produces,
	}
	if len(sr.Consumes) > 0 {
		data["consumes"] = sr.Consumes
	}
	param := paramsMap[sr.Requests.FullName]
	data["parameters"] = param.AsRequest(sr.Requests.Location)
	response := paramsMap[sr.Responses.FullName]
	data["responses"] = map[string]interface{}{
		"200": response.AsResponse(),
	}
	return data
}
