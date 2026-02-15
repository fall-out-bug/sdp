package agents

import (
	"strings"
)

// SchemaInferrer infers request/response schemas from code
type SchemaInferrer struct{}

// ContractGenerator enhances contracts with inferred schemas
type ContractGenerator struct {
	inferrer *SchemaInferrer
}

// NewSchemaInferrer creates a new schema inferrer
func NewSchemaInferrer() *SchemaInferrer {
	return &SchemaInferrer{}
}

// NewContractGenerator creates a new contract generator
func NewContractGenerator() *ContractGenerator {
	return &ContractGenerator{
		inferrer: NewSchemaInferrer(),
	}
}

// mapGoTypeToJSON maps Go types to JSON types
func mapGoTypeToJSON(goType string) string {
	switch goType {
	case "string":
		return "string"
	case "int", "int32", "int64", "uint32", "uint64":
		return "integer"
	case "float32", "float64":
		return "number"
	case "bool":
		return "boolean"
	default:
		if strings.HasPrefix(goType, "[]") {
			return "array"
		}
		if strings.HasPrefix(goType, "map[") {
			return "object"
		}
		return "object"
	}
}

// mapTSTypeToJSON maps TypeScript types to JSON types
func mapTSTypeToJSON(tsType string) string {
	tsType = strings.TrimSuffix(tsType, "?")
	tsType = strings.TrimSpace(tsType)

	switch tsType {
	case "string":
		return "string"
	case "number":
		return "number"
	case "boolean":
		return "boolean"
	default:
		if strings.HasPrefix(tsType, "Array<") {
			return "array"
		}
		if strings.HasPrefix(tsType, "Record<") {
			return "object"
		}
		return "object"
	}
}
