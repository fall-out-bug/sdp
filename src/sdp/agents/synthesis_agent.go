package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fall-out-bug/sdp/src/sdp/synthesis"
	"gopkg.in/yaml.v3"
)

// ContractSynthesizer manages contract synthesis from requirements
type ContractSynthesizer struct {
	supervisor *synthesis.Supervisor
}

// ContractRequirements represents parsed feature requirements
type ContractRequirements struct {
	FeatureName string         `yaml:"feature_name"`
	Endpoints   []EndpointSpec `yaml:"endpoints"`
}

// EndpointSpec represents an API endpoint specification
type EndpointSpec struct {
	Path     string     `yaml:"path"`
	Method   string     `yaml:"method"`
	Request  SchemaSpec `yaml:"request"`
	Response SchemaSpec `yaml:"response"`
}

// SchemaSpec represents a request/response schema
type SchemaSpec struct {
	Fields []FieldSpec `yaml:"fields"`
}

// FieldSpec represents a field in a schema
type FieldSpec struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"`
	Required bool   `yaml:"required"`
}

// OpenAPIContract represents an OpenAPI 3.0 contract
type OpenAPIContract struct {
	OpenAPI string     `yaml:"openapi"`
	Info    InfoSpec   `yaml:"info"`
	Paths   PathsSpec  `yaml:"paths"`
}

// InfoSpec represents OpenAPI info block
type InfoSpec struct {
	Title   string `yaml:"title"`
	Version string `yaml:"version"`
}

// PathsSpec represents OpenAPI paths block
type PathsSpec map[string]PathSpec

// PathSpec represents a single path in OpenAPI
type PathSpec map[string]OperationSpec

// OperationSpec represents an HTTP operation in OpenAPI
type OperationSpec struct {
	Summary    string      `yaml:"summary"`
	RequestBody   *RequestSpec  `yaml:"requestBody,omitempty"`
	Responses ResponsesSpec `yaml:"responses"`
}

// RequestSpec represents OpenAPI request body
type RequestSpec struct {
	Required bool                 `yaml:"required"`
	Content  map[string]MediaSpec `yaml:"content"`
}

// ResponsesSpec represents OpenAPI responses
type ResponsesSpec map[string]ResponseSpec

// ResponseSpec represents an OpenAPI response
type ResponseSpec struct {
	Description string                 `yaml:"description"`
	Content     map[string]MediaSpec   `yaml:"content"`
}

// MediaSpec represents OpenAPI media type
type MediaSpec struct {
	Schema SchemaRefSpec `yaml:"schema"`
}

// SchemaRefSpec represents OpenAPI schema reference
type SchemaRefSpec struct {
	Type       string                    `yaml:"type,omitempty"`
	Properties map[string]PropertySpec   `yaml:"properties,omitempty"`
	Required   []string                  `yaml:"required,omitempty"`
}

// PropertySpec represents a property in schema
type PropertySpec struct {
	Type string `yaml:"type"`
}

// EndpointProposal represents a proposed endpoint change
type EndpointProposal struct {
	Path   string
	Method string
}

// NewContractSynthesizer creates a new contract synthesizer
func NewContractSynthesizer() *ContractSynthesizer {
	engine := synthesis.DefaultRuleEngine()
	supervisor := synthesis.NewSupervisor(engine, 3) // max 3 agents

	return &ContractSynthesizer{
		supervisor: supervisor,
	}
}

// AnalyzeRequirements parses a requirements markdown file
func (cs *ContractSynthesizer) AnalyzeRequirements(reqPath string) (*ContractRequirements, error) {
	// Read requirements file
	content, err := os.ReadFile(reqPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read requirements: %w", err)
	}

	// Extract feature name from filename
	featureName := strings.TrimSuffix(filepath.Base(reqPath), "-requirements.md")
	featureName = strings.TrimPrefix(featureName, "sdp-")

	// Validate feature name to prevent injection attacks
	if err := validateFeatureName(featureName); err != nil {
		return nil, fmt.Errorf("invalid feature name in path %q: %w", reqPath, err)
	}

	// Parse endpoints from markdown
	endpoints, err := cs.parseEndpointsFromMarkdown(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoints: %w", err)
	}

	return &ContractRequirements{
		FeatureName: featureName,
		Endpoints:   endpoints,
	}, nil
}

const (
	// MaxContentLength limits content size to prevent ReDoS
	MaxContentLength = 100000
	// MaxFieldCount limits number of fields to prevent resource exhaustion
	MaxFieldCount = 100
)

var (
	// allowedHTTPMethods is a whitelist of permitted HTTP methods
	allowedHTTPMethods = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"HEAD":    true,
		"OPTIONS": true,
	}
)

// validateFeatureName checks if a feature name is valid
func validateFeatureName(name string) error {
	// Feature names must be lowercase alphanumeric with dashes only
	// This prevents injection attacks via malicious feature names
	matched, err := regexp.MatchString(`^[a-z0-9-]+$`, name)
	if err != nil {
		return fmt.Errorf("feature name validation failed: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid feature name %q: must contain only lowercase letters, numbers, and dashes", name)
	}
	return nil
}

// validateHTTPMethod checks if an HTTP method is allowed
func validateHTTPMethod(method string) error {
	// Method must contain only uppercase letters (no special characters, no spaces)
	if !regexp.MustCompile(`^[A-Z]+$`).MatchString(method) {
		return fmt.Errorf("invalid HTTP method %q: must contain only uppercase letters", method)
	}

	// Normalize to uppercase for whitelist check
	normalizedMethod := strings.ToUpper(method)
	if !allowedHTTPMethods[normalizedMethod] {
		return fmt.Errorf("invalid HTTP method %q: must be one of %v", method, getAllowedMethods())
	}
	return nil
}

// getAllowedMethods returns the list of allowed HTTP methods
func getAllowedMethods() []string {
	methods := make([]string, 0, len(allowedHTTPMethods))
	for method := range allowedHTTPMethods {
		methods = append(methods, method)
	}
	return methods
}

// validateEndpointPath checks if an endpoint path is valid
func validateEndpointPath(path string) error {
	// Paths must start with /
	if !strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid endpoint path %q: must start with /", path)
	}

	// Paths should not contain .. to prevent directory traversal
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid endpoint path %q: must not contain ..", path)
	}

	// Limit path length to prevent DoS
	if len(path) > 500 {
		return fmt.Errorf("invalid endpoint path %q: too long (max 500 chars)", path)
	}

	// Paths should only contain safe characters
	// Allow: /, alphanumeric, -, _, :, {}, * (for path params and wildcards)
	matched, err := regexp.MatchString(`^/[a-zA-Z0-9\-_{}/*.:]+$`, path)
	if err != nil {
		return fmt.Errorf("endpoint path validation failed: %w", err)
	}
	if !matched {
		return fmt.Errorf("invalid endpoint path %q: contains invalid characters", path)
	}

	return nil
}

// sanitizeFieldName sanitizes field names to prevent injection
func sanitizeFieldName(name string) (string, error) {
	// Field names must be alphanumeric with underscores
	matched, err := regexp.MatchString(`^[a-zA-Z_][a-zA-Z0-9_]*$`, name)
	if err != nil {
		return "", fmt.Errorf("field name validation failed: %w", err)
	}
	if !matched {
		return "", fmt.Errorf("invalid field name %q: must be alphanumeric with underscores, starting with letter or underscore", name)
	}
	return name, nil
}

// parseEndpointsFromMarkdown extracts endpoint specifications from markdown
func (cs *ContractSynthesizer) parseEndpointsFromMarkdown(content string) ([]EndpointSpec, error) {
	var endpoints []EndpointSpec

	// Check content size before processing
	if len(content) > MaxContentLength {
		return nil, fmt.Errorf("content too large for parsing: %d bytes (max %d)", len(content), MaxContentLength)
	}

	// FIXED: Added length limits to prevent ReDoS attacks
	// Regex to match endpoint headers: ### POST /api/v1/telemetry/events or - POST /api/v1/telemetry/events
	// Limit path length to 500 chars to prevent catastrophic backtracking
	// Capture method as non-slash, non-space characters on the same line
	// Use [^/\s]+ to ensure we only capture word-like characters
	endpointRe := regexp.MustCompile(`(?:###|-)\s+([^/\s]{1,20})\s+(/[^\s]{1,500})`)
	matches := endpointRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		// Validate HTTP method (whitelist check) and normalize to uppercase
		// match[1] is the raw method string from the regex
		rawMethod := match[1]
		method := strings.ToUpper(rawMethod)

		// Strict validation: the raw method must be exactly the same as the uppercased version
		// This prevents "GET; DROP TABLE users" from being accepted as "GET"
		if rawMethod != method {
			return nil, fmt.Errorf("invalid HTTP method %q: must be uppercase only, got mixed case or special characters", rawMethod)
		}

		if err := validateHTTPMethod(method); err != nil {
			return nil, fmt.Errorf("invalid HTTP method in endpoint: %w", err)
		}

		// Validate endpoint path format
		path := match[2]
		if err := validateEndpointPath(path); err != nil {
			return nil, fmt.Errorf("invalid endpoint path: %w", err)
		}

		endpoint := EndpointSpec{
			Path:   path,
			Method: method,
			Request:  SchemaSpec{Fields: []FieldSpec{}},
			Response: SchemaSpec{Fields: []FieldSpec{}},
		}

		// Extract request/response fields from markdown lists
		// Format: - Request: {field1, field2} or - field_name: type
		lines := strings.Split(content, "\n")
		inRequestSection := false
		inResponseSection := false

		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.Contains(line, fmt.Sprintf("### %s", match[1])) {
				// Start of this endpoint's section
				inRequestSection = false
				inResponseSection = false
				continue
			}

			if strings.Contains(line, "Request:") && strings.Contains(content, line) {
				inRequestSection = true
				inResponseSection = false

				// FIXED: Limited field content to 500 chars to prevent ReDoS
				// Parse inline field specification: {field1, field2}
				inlineRe := regexp.MustCompile(`Request:\s*\{([^}]{1,500})\}`)
				if inlineMatches := inlineRe.FindStringSubmatch(line); len(inlineMatches) > 1 {
					fields := strings.Split(inlineMatches[1], ",")
					for _, f := range fields {
						f = strings.TrimSpace(f)
						if f != "" {
							// Sanitize field name to prevent injection
							sanitized, err := sanitizeFieldName(f)
							if err != nil {
								return nil, fmt.Errorf("invalid request field name: %w", err)
							}
							endpoint.Request.Fields = append(endpoint.Request.Fields, FieldSpec{
								Name:     sanitized,
								Type:     "string",
								Required: true,
							})
						}
					}
				}
				continue
			}

			if strings.Contains(line, "Response:") {
				inRequestSection = false
				inResponseSection = true

				// FIXED: Limited field content to 500 chars to prevent ReDoS
				// Parse inline field specification
				inlineRe := regexp.MustCompile(`Response:\s*\{([^}]{1,500})\}`)
				if inlineMatches := inlineRe.FindStringSubmatch(line); len(inlineMatches) > 1 {
					fields := strings.Split(inlineMatches[1], ",")
					for _, f := range fields {
						f = strings.TrimSpace(f)
						if f != "" {
							// Sanitize field name to prevent injection
							sanitized, err := sanitizeFieldName(f)
							if err != nil {
								return nil, fmt.Errorf("invalid response field name: %w", err)
							}
							endpoint.Response.Fields = append(endpoint.Response.Fields, FieldSpec{
								Name:     sanitized,
								Type:     "string",
								Required: true,
							})
						}
					}
				}
				continue
			}

			// Parse bullet point fields: - field_name: type
			if strings.HasPrefix(line, "-") {
				// FIXED: Added length limits to prevent ReDoS
				fieldRe := regexp.MustCompile(`-\s*(\w{1,100}):\s*(\w{1,50})`)
				if fieldMatches := fieldRe.FindStringSubmatch(line); len(fieldMatches) > 2 {
					// Enforce field count limit
					if len(endpoint.Request.Fields) >= MaxFieldCount || len(endpoint.Response.Fields) >= MaxFieldCount {
						return nil, fmt.Errorf("too many fields (max %d)", MaxFieldCount)
					}

					// Sanitize field name
					fieldName := fieldMatches[1]
					sanitizedName, err := sanitizeFieldName(fieldName)
					if err != nil {
						return nil, fmt.Errorf("invalid field name in bullet point: %w", err)
					}

					field := FieldSpec{
						Name: sanitizedName,
						Type: fieldMatches[2],
					}

					if inRequestSection {
						endpoint.Request.Fields = append(endpoint.Request.Fields, field)
					} else if inResponseSection {
						endpoint.Response.Fields = append(endpoint.Response.Fields, field)
					}
				}
			}
		}

		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

// ProposeContract generates an initial OpenAPI contract from requirements
func (cs *ContractSynthesizer) ProposeContract(requirements *ContractRequirements) (*OpenAPIContract, error) {
	contract := &OpenAPIContract{
		OpenAPI: "3.0.0",
		Info: InfoSpec{
			Title:   fmt.Sprintf("%s API", strings.Title(requirements.FeatureName)),
			Version: "1.0.0",
		},
		Paths: make(PathsSpec),
	}

	// Convert endpoints to OpenAPI paths
	for _, endpoint := range requirements.Endpoints {
		path := cs.endpointToPathSpec(endpoint)
		contract.Paths[endpoint.Path] = path
	}

	return contract, nil
}

// endpointToPathSpec converts an endpoint spec to OpenAPI path spec
func (cs *ContractSynthesizer) endpointToPathSpec(endpoint EndpointSpec) PathSpec {
	pathSpec := make(PathSpec)

	operation := OperationSpec{
		Summary: fmt.Sprintf("%s %s", endpoint.Method, endpoint.Path),
		Responses: ResponsesSpec{
			"200": ResponseSpec{
				Description: "Success",
				Content: map[string]MediaSpec{
					"application/json": {
						Schema: cs.schemaSpecToSchemaRef(endpoint.Response),
					},
				},
			},
		},
	}

	// Add request body for POST/PUT/PATCH
	if endpoint.Method == "POST" || endpoint.Method == "PUT" || endpoint.Method == "PATCH" {
		operation.RequestBody = &RequestSpec{
			Required: true,
			Content: map[string]MediaSpec{
				"application/json": {
					Schema: cs.schemaSpecToSchemaRef(endpoint.Request),
				},
			},
		}
	}

	pathSpec[strings.ToLower(endpoint.Method)] = operation
	return pathSpec
}

// schemaSpecToSchemaRef converts a schema spec to OpenAPI schema reference
func (cs *ContractSynthesizer) schemaSpecToSchemaRef(schema SchemaSpec) SchemaRefSpec {
	if len(schema.Fields) == 0 {
		return SchemaRefSpec{Type: "object"}
	}

	properties := make(map[string]PropertySpec)
	required := []string{}

	for _, field := range schema.Fields {
		properties[field.Name] = PropertySpec{Type: field.Type}
		if field.Required {
			required = append(required, field.Name)
		}
	}

	return SchemaRefSpec{
		Type:       "object",
		Properties: properties,
		Required:   required,
	}
}

// ApplySynthesisRules applies synthesis rules to resolve conflicts in proposals
func (cs *ContractSynthesizer) ApplySynthesisRules(proposals []*synthesis.Proposal) (*synthesis.SynthesisResult, error) {
	// Use the internal synthesizer
	synthesizer := synthesis.NewSynthesizer()

	for _, proposal := range proposals {
		synthesizer.AddProposal(proposal)
	}

	return synthesizer.Synthesize()
}

// WriteContract writes the agreed contract to a YAML file
func (cs *ContractSynthesizer) WriteContract(contract *OpenAPIContract, outputPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}

	// Write file
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write contract: %w", err)
	}

	return nil
}

// SynthesizeContract performs end-to-end contract synthesis
func (cs *ContractSynthesizer) SynthesizeContract(featureName, reqPath, outputPath string) (*synthesis.SynthesisResult, error) {
	// Step 1: Analyze requirements
	requirements, err := cs.AnalyzeRequirements(reqPath)
	if err != nil {
		return nil, fmt.Errorf("analyze requirements failed: %w", err)
	}

	// Step 2: Propose initial contract
	contract, err := cs.ProposeContract(requirements)
	if err != nil {
		return nil, fmt.Errorf("propose contract failed: %w", err)
	}

	// Step 3: Request agent reviews (in parallel)
	// For now, we'll simulate this with the synthesis rules
	// In a full implementation, this would use the Task tool to spawn agents

	// Step 4: Apply synthesis rules
	// For now, we'll use unanimous rule (no conflicts)
	proposals := []*synthesis.Proposal{
		synthesis.NewProposal(
			"architect",
			contract,
			1.0,
			"Initial contract from requirements",
		),
	}

	result, err := cs.ApplySynthesisRules(proposals)
	if err != nil {
		return nil, fmt.Errorf("apply synthesis rules failed: %w", err)
	}

	// Step 5: Write agreed contract
	finalContract := result.Solution.(*OpenAPIContract)
	if err := cs.WriteContract(finalContract, outputPath); err != nil {
		return nil, fmt.Errorf("write contract failed: %w", err)
	}

	return result, nil
}
