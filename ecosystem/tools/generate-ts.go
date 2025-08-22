package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/CodeClarityCE/utility-types/ecosystem"
)

// TypeScriptType represents a TypeScript type definition
type TypeScriptType struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Fields      []TypeScriptField      `json:"fields,omitempty"`
	Values      []string               `json:"values,omitempty"`      // for enums
	Extends     string                 `json:"extends,omitempty"`     // for inheritance
	Implements  []string               `json:"implements,omitempty"`  // for interface implementation
	IsInterface bool                   `json:"isInterface"`
	IsEnum      bool                   `json:"isEnum"`
	Comment     string                 `json:"comment,omitempty"`
}

// TypeScriptField represents a field in a TypeScript interface/type
type TypeScriptField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Optional bool   `json:"optional"`
	Comment  string `json:"comment,omitempty"`
}

// GoToTypeScriptConverter converts Go structs to TypeScript definitions
type GoToTypeScriptConverter struct {
	types map[string]TypeScriptType
}

// NewConverter creates a new converter
func NewConverter() *GoToTypeScriptConverter {
	return &GoToTypeScriptConverter{
		types: make(map[string]TypeScriptType),
	}
}

// ConvertStruct converts a Go struct to TypeScript type definition
func (c *GoToTypeScriptConverter) ConvertStruct(structType reflect.Type, name string) TypeScriptType {
	tsType := TypeScriptType{
		Name:        name,
		Type:        "interface",
		IsInterface: true,
		Fields:      []TypeScriptField{},
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		
		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		tsField := TypeScriptField{
			Name:     c.getJsonFieldName(field),
			Type:     c.convertGoTypeToTypescript(field.Type),
			Optional: false,
		}

		// Check if field is optional (pointer type)
		if field.Type.Kind() == reflect.Ptr {
			tsField.Optional = true
			tsField.Type = c.convertGoTypeToTypescript(field.Type.Elem())
		}

		// Extract comment from struct tag if available
		if tag := field.Tag.Get("comment"); tag != "" {
			tsField.Comment = tag
		}

		tsType.Fields = append(tsType.Fields, tsField)
	}

	return tsType
}

// getJsonFieldName extracts the JSON field name from struct tags
func (c *GoToTypeScriptConverter) getJsonFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return c.camelCase(field.Name)
	}

	// Parse JSON tag (e.g., "name,omitempty")
	parts := strings.Split(jsonTag, ",")
	if parts[0] != "" {
		return parts[0]
	}

	return c.camelCase(field.Name)
}

// convertGoTypeToTypescript converts Go types to TypeScript types
func (c *GoToTypeScriptConverter) convertGoTypeToTypescript(goType reflect.Type) string {
	switch goType.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		 reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		 reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		elemType := c.convertGoTypeToTypescript(goType.Elem())
		return elemType + "[]"
	case reflect.Map:
		keyType := c.convertGoTypeToTypescript(goType.Key())
		valueType := c.convertGoTypeToTypescript(goType.Elem())
		
		// TypeScript record type
		if keyType == "string" {
			return fmt.Sprintf("Record<%s, %s>", keyType, valueType)
		}
		return fmt.Sprintf("{ [key: %s]: %s }", keyType, valueType)
	case reflect.Ptr:
		return c.convertGoTypeToTypescript(goType.Elem())
	case reflect.Interface:
		return "any"
	case reflect.Struct:
		// Handle special cases
		if goType == reflect.TypeOf(time.Time{}) {
			return "string" // ISO date string
		}
		
		// For custom structs, use the struct name
		return goType.Name()
	default:
		return "any"
	}
}

// camelCase converts PascalCase to camelCase
func (c *GoToTypeScriptConverter) camelCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// GenerateTypeScript generates TypeScript type definitions from Go types
func (c *GoToTypeScriptConverter) GenerateTypeScript(types map[string]TypeScriptType) string {
	var sb strings.Builder

	// Write header
	sb.WriteString("// Generated TypeScript definitions from Go structs\n")
	sb.WriteString("// DO NOT EDIT - This file is auto-generated\n")
	sb.WriteString("// Generated at: " + time.Now().Format(time.RFC3339) + "\n\n")

	// Write types
	for _, tsType := range types {
		c.writeTypeDefinition(&sb, tsType)
		sb.WriteString("\n")
	}

	return sb.String()
}

// writeTypeDefinition writes a single TypeScript type definition
func (c *GoToTypeScriptConverter) writeTypeDefinition(sb *strings.Builder, tsType TypeScriptType) {
	// Write comment if available
	if tsType.Comment != "" {
		sb.WriteString("/**\n * " + tsType.Comment + "\n */\n")
	}

	if tsType.IsEnum {
		// Write enum
		sb.WriteString(fmt.Sprintf("export enum %s {\n", tsType.Name))
		for _, value := range tsType.Values {
			sb.WriteString(fmt.Sprintf("  %s,\n", value))
		}
		sb.WriteString("}")
	} else if tsType.IsInterface {
		// Write interface
		sb.WriteString(fmt.Sprintf("export interface %s", tsType.Name))
		
		// Write extends clause
		if tsType.Extends != "" {
			sb.WriteString(fmt.Sprintf(" extends %s", tsType.Extends))
		}

		// Write implements clause
		if len(tsType.Implements) > 0 {
			sb.WriteString(fmt.Sprintf(" implements %s", strings.Join(tsType.Implements, ", ")))
		}

		sb.WriteString(" {\n")

		// Write fields
		for _, field := range tsType.Fields {
			if field.Comment != "" {
				sb.WriteString(fmt.Sprintf("  /** %s */\n", field.Comment))
			}
			
			optional := ""
			if field.Optional {
				optional = "?"
			}
			
			sb.WriteString(fmt.Sprintf("  %s%s: %s;\n", field.Name, optional, field.Type))
		}

		sb.WriteString("}")
	} else {
		// Write type alias
		sb.WriteString(fmt.Sprintf("export type %s = %s;", tsType.Name, tsType.Type))
	}
}

// main function to generate TypeScript definitions
func main() {
	converter := NewConverter()

	// Convert ecosystem types
	ecosystemInfoType := reflect.TypeOf(ecosystem.EcosystemInfo{})
	converter.types["EcosystemInfo"] = converter.ConvertStruct(ecosystemInfoType, "EcosystemInfo")

	detectedLanguageType := reflect.TypeOf(ecosystem.DetectedLanguage{})
	converter.types["DetectedLanguage"] = converter.ConvertStruct(detectedLanguageType, "DetectedLanguage")

	// Add PluginEcosystemMap type
	converter.types["PluginEcosystemMap"] = TypeScriptType{
		Name:        "PluginEcosystemMap",
		Type:        "Record<string, EcosystemInfo>",
		IsInterface: false,
	}

	// Add MergeStrategy enum
	converter.types["MergeStrategy"] = TypeScriptType{
		Name:   "MergeStrategy",
		Type:   "enum",
		IsEnum: true,
		Values: []string{
			"UNION = 'union'",
			"INTERSECTION = 'intersection'",
			"PRIORITY = 'priority'",
		},
	}

	// Generate TypeScript
	typescript := converter.GenerateTypeScript(converter.types)

	// Determine output path
	outputDir := "."
	if len(os.Args) > 1 {
		outputDir = os.Args[1]
	}

	outputFile := filepath.Join(outputDir, "ecosystem.gen.ts")

	// Write to file
	err := os.WriteFile(outputFile, []byte(typescript), 0644)
	if err != nil {
		fmt.Printf("Error writing TypeScript file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated TypeScript definitions: %s\n", outputFile)

	// Also generate a JSON schema for runtime validation
	jsonSchema := map[string]interface{}{
		"$schema": "http://json-schema.org/draft-07/schema#",
		"title":   "Ecosystem Types",
		"types":   converter.types,
	}

	jsonBytes, err := json.MarshalIndent(jsonSchema, "", "  ")
	if err != nil {
		fmt.Printf("Error generating JSON schema: %v\n", err)
		os.Exit(1)
	}

	schemaFile := filepath.Join(outputDir, "ecosystem.schema.json")
	err = os.WriteFile(schemaFile, jsonBytes, 0644)
	if err != nil {
		fmt.Printf("Error writing JSON schema file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated JSON schema: %s\n", schemaFile)
}