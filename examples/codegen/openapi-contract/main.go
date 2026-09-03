// openapi-contract demonstrates consumer-owned API contract generation.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/kucjac/gentools/parser"
	"github.com/kucjac/gentools/types"
)

type document struct {
	OpenAPI    string              `json:"openapi"`
	Info       info                `json:"info"`
	Paths      map[string]pathItem `json:"paths"`
	Components components          `json:"components"`
}
type info struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}
type components struct {
	Schemas map[string]schema `json:"schemas"`
}
type pathItem map[string]operation
type operation struct {
	Summary    string              `json:"summary"`
	Parameters []parameter         `json:"parameters,omitempty"`
	Responses  map[string]response `json:"responses"`
}
type parameter struct {
	Name     string `json:"name"`
	In       string `json:"in"`
	Required bool   `json:"required"`
	Schema   schema `json:"schema"`
}
type response struct {
	Description string               `json:"description"`
	Content     map[string]mediaType `json:"content,omitempty"`
}
type mediaType struct {
	Schema schema `json:"schema"`
}
type schema struct {
	Ref         string            `json:"$ref,omitempty"`
	Type        string            `json:"type,omitempty"`
	Format      string            `json:"format,omitempty"`
	Description string            `json:"description,omitempty"`
	Items       *schema           `json:"items,omitempty"`
	Properties  map[string]schema `json:"properties,omitempty"`
}

type modelRegistry struct {
	byKey      map[string]*types.Struct
	byName     map[string][]*types.Struct
	schemaName map[string]string
}

var routePattern = regexp.MustCompile(`^/(?:[A-Za-z0-9._~-]+|\{[A-Za-z_][A-Za-z0-9_]*\})(?:/(?:[A-Za-z0-9._~-]+|\{[A-Za-z_][A-Za-z0-9_]*\}))*$`)

func main() {
	input := flag.String("input", "", "directory containing annotated Go source")
	output := flag.String("output", "", "generated OpenAPI JSON path")
	flag.Parse()
	if *input == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "-input and -output are required")
		os.Exit(2)
	}
	if err := generate(*input, *output); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type annotation struct{ Method, Route, Summary, Response string }

func parseAnnotation(comment, operationName string) (annotation, error) {
	var line string
	for _, candidate := range strings.Split(comment, "\n") {
		if strings.HasPrefix(strings.TrimSpace(candidate), "@openapi") {
			line = strings.TrimSpace(candidate)
			break
		}
	}
	if line == "" {
		return annotation{}, fmt.Errorf("operation %s: missing @openapi annotation", operationName)
	}
	tokens, err := tokenize(strings.TrimSpace(strings.TrimPrefix(line, "@openapi")))
	if err != nil {
		return annotation{}, fmt.Errorf("operation %s: malformed annotation: %w", operationName, err)
	}
	result := annotation{}
	seen := map[string]bool{}
	for _, token := range tokens {
		parts := strings.SplitN(token, "=", 2)
		if len(parts) != 2 || parts[0] == "" {
			return annotation{}, fmt.Errorf("operation %s: malformed annotation token %q", operationName, token)
		}
		key, value := parts[0], parts[1]
		if seen[key] {
			return annotation{}, fmt.Errorf("operation %s: duplicate annotation key %q", operationName, key)
		}
		seen[key] = true
		switch key {
		case "method":
			result.Method = strings.ToUpper(value)
		case "route":
			result.Route = value
		case "summary":
			result.Summary = value
		case "response":
			result.Response = value
		default:
			return annotation{}, fmt.Errorf("operation %s: unknown annotation key %q", operationName, key)
		}
	}
	if result.Method == "" || result.Route == "" || result.Summary == "" || result.Response == "" {
		return annotation{}, fmt.Errorf("operation %s: annotation requires method, route, summary, and response", operationName)
	}
	switch result.Method {
	case "GET":
	default:
		return annotation{}, fmt.Errorf("operation %s: unsupported HTTP method %q", operationName, result.Method)
	}
	if !routePattern.MatchString(result.Route) {
		return annotation{}, fmt.Errorf("operation %s: malformed route %q", operationName, result.Route)
	}
	if !strings.Contains(result.Response, ":") {
		return annotation{}, fmt.Errorf("operation %s: malformed response %q", operationName, result.Response)
	}
	return result, nil
}

func tokenize(input string) ([]string, error) {
	var result []string
	for len(strings.TrimSpace(input)) > 0 {
		input = strings.TrimSpace(input)
		start := 0
		quoted := false
		escaped := false
		for i, r := range input {
			if escaped {
				escaped = false
				continue
			}
			if r == '\\' && quoted {
				escaped = true
				continue
			}
			if r == '"' {
				quoted = !quoted
				continue
			}
			if r == ' ' && !quoted {
				start = i
				break
			}
			start = len(input)
		}
		if quoted {
			return nil, errors.New("unterminated quoted value")
		}
		token := input[:start]
		input = input[start:]
		if token != "" {
			if strings.Contains(token, "\"") {
				value, err := strconv.Unquote(strings.SplitN(token, "=", 2)[1])
				if err != nil {
					return nil, err
				}
				token = strings.SplitN(token, "=", 2)[0] + "=" + value
			}
			result = append(result, token)
		}
	}
	return result, nil
}

func generate(input, output string) error {
	pkgs, err := parser.LoadPackages(parser.LoadConfig{Paths: []string{input}, WithComments: true})
	if err != nil {
		return err
	}
	doc := document{OpenAPI: "3.0.3", Info: info{Title: "Gentools example API", Version: "1.0.0"}, Paths: map[string]pathItem{}, Components: components{Schemas: map[string]schema{}}}
	models := &modelRegistry{byKey: map[string]*types.Struct{}, byName: map[string][]*types.Struct{}, schemaName: map[string]string{}}
	for _, pkg := range pkgs {
		for name, value := range pkg.Types {
			if model, ok := value.(*types.Struct); ok && !strings.HasPrefix(name, "_") {
				key := modelKey(model)
				models.byKey[key] = model
				models.byName[name] = append(models.byName[name], model)
			}
		}
	}
	for name, candidates := range models.byName {
		for _, model := range candidates {
			componentName := name
			if len(candidates) > 1 {
				componentName = sanitizeComponentName(model.Pkg.Identifier + "_" + name)
			}
			models.schemaName[modelKey(model)] = componentName
		}
	}
	seenRoutes := map[string]string{}
	for _, pkg := range pkgs {
		for name, value := range pkg.Types {
			fn, ok := value.(*types.Function)
			if !ok {
				continue
			}
			if !strings.Contains(fn.Comment, "@openapi") {
				continue
			}
			ann, parseErr := parseAnnotation(fn.Comment, name)
			if parseErr != nil {
				return parseErr
			}
			key := ann.Method + " " + ann.Route
			if previous, exists := seenRoutes[key]; exists {
				return fmt.Errorf("duplicate operation route %s claimed by %s and %s", key, previous, name)
			}
			seenRoutes[key] = name
			statusModel := strings.SplitN(ann.Response, ":", 2)
			status, modelName := statusModel[0], statusModel[1]
			if _, err := strconv.Atoi(status); err != nil {
				return fmt.Errorf("operation %s: invalid response status %q", name, status)
			}
			model, resolveErr := models.resolve(modelName, pkg)
			if resolveErr != nil {
				return fmt.Errorf("operation %s: unresolved model %q", name, modelName)
			}
			modelSchema, schemaErr := buildModel(model, doc.Components.Schemas, models)
			if schemaErr != nil {
				return schemaErr
			}
			doc.Components.Schemas[models.componentName(model)] = modelSchema
			op := operation{Summary: ann.Summary, Responses: map[string]response{status: {Description: "Successful response", Content: map[string]mediaType{"application/json": {Schema: schema{Ref: "#/components/schemas/" + modelName}}}}}}
			for _, segment := range strings.Split(ann.Route, "/") {
				if strings.HasPrefix(segment, "{") {
					param := strings.Trim(segment, "{}")
					op.Parameters = append(op.Parameters, parameter{Name: param, In: "path", Required: true, Schema: schema{Type: "integer", Format: "int64"}})
				}
			}
			if doc.Paths[ann.Route] == nil {
				doc.Paths[ann.Route] = pathItem{}
			}
			doc.Paths[ann.Route][strings.ToLower(ann.Method)] = op
		}
	}
	return writeDocument(output, doc)
}

func modelKey(model *types.Struct) string { return model.Pkg.Path + "." + model.TypeName }

func sanitizeComponentName(name string) string {
	return strings.NewReplacer("/", "_", "-", "_", ".", "_").Replace(name)
}

func (r *modelRegistry) componentName(model *types.Struct) string {
	return r.schemaName[modelKey(model)]
}

func (r *modelRegistry) resolve(name string, context *types.Package) (*types.Struct, error) {
	if local, ok := context.GetStruct(name); ok {
		return local, nil
	}
	if candidates := r.byName[name]; len(candidates) == 1 {
		return candidates[0], nil
	}
	for key, model := range r.byKey {
		if key == name || strings.HasSuffix(key, "."+name) || model.Pkg.Identifier+"."+model.TypeName == name {
			return model, nil
		}
	}
	return nil, errors.New("model not found")
}

func buildModel(model *types.Struct, definitions map[string]schema, models *modelRegistry) (schema, error) {
	result := schema{Type: "object", Description: strings.TrimSpace(model.Comment), Properties: map[string]schema{}}
	for _, field := range model.Fields {
		if field.Name == "" || !astExported(field.Name) {
			continue
		}
		jsonName, ok := field.Tag.Lookup("json")
		if ok {
			jsonName = strings.Split(jsonName, ",")[0]
			if jsonName == "-" {
				continue
			}
			if jsonName == "" {
				jsonName = field.Name
			}
		} else {
			jsonName = field.Name
		}
		fieldSchema, err := schemaFor(field.Type, definitions, models)
		if err != nil {
			return schema{}, fmt.Errorf("model %s field %s: %w", model.TypeName, field.Name, err)
		}
		fieldSchema.Description = strings.TrimSpace(field.Comment)
		result.Properties[jsonName] = fieldSchema
	}
	return result, nil
}

func astExported(name string) bool { return name != "" && strings.ToUpper(name[:1]) == name[:1] }

func schemaFor(tp types.Type, definitions map[string]schema, models *modelRegistry) (schema, error) {
	switch value := tp.(type) {
	case *types.Pointer:
		return schemaFor(value.Elem(), definitions, models)
	case *types.Array:
		inner, err := schemaFor(value.Elem(), definitions, models)
		return schema{Type: "array", Items: &inner}, err
	case *types.Struct:
		componentName := models.componentName(value)
		if _, ok := definitions[componentName]; !ok {
			built, err := buildModel(value, definitions, models)
			if err != nil {
				return schema{}, err
			}
			definitions[componentName] = built
		}
		return schema{Ref: "#/components/schemas/" + componentName}, nil
	case *types.Map, *types.Interface, *types.Chan, *types.Function:
		return schema{}, fmt.Errorf("unsupported field shape %s", tp.String())
	default:
		kind := tp.Kind()
		if !kind.IsBuiltin() {
			return schema{}, fmt.Errorf("unsupported field shape %s", tp.String())
		}
		if kind == types.KindString {
			return schema{Type: "string"}, nil
		}
		if kind.IsNumber() {
			format := ""
			if kind == types.KindInt64 {
				format = "int64"
			}
			return schema{Type: "integer", Format: format}, nil
		}
		if kind == types.KindBool {
			return schema{Type: "boolean"}, nil
		}
		return schema{}, fmt.Errorf("unsupported field shape %s", tp.String())
	}
}

func writeDocument(output string, doc document) error {
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	dir := filepath.Dir(output)
	temp, err := os.CreateTemp(dir, ".openapi-*")
	if err != nil {
		return err
	}
	tempName := temp.Name()
	defer os.Remove(tempName)
	if _, err = temp.Write(data); err != nil {
		temp.Close()
		return err
	}
	if err = temp.Close(); err != nil {
		return err
	}
	return os.Rename(tempName, output)
}
