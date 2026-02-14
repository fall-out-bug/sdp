package agents

import (
	"fmt"
	"regexp"
	"strings"
)

// parseEndpointsFromMarkdown extracts endpoint specifications from markdown
func (cs *ContractSynthesizer) parseEndpointsFromMarkdown(content string) ([]EndpointSpec, error) {
	var endpoints []EndpointSpec

	if len(content) > MaxContentLength {
		return nil, fmt.Errorf("content too large for parsing: %d bytes (max %d)", len(content), MaxContentLength)
	}

	endpointRe := regexp.MustCompile(`(?:###|-)\s+([^/\s]{1,20})\s+(/[^\s]{1,500})`)
	matches := endpointRe.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		endpoint, err := cs.parseEndpoint(match, content)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints, nil
}

// parseEndpoint parses a single endpoint from regex match
func (cs *ContractSynthesizer) parseEndpoint(match []string, content string) (EndpointSpec, error) {
	rawMethod := match[1]
	method := strings.ToUpper(rawMethod)

	if rawMethod != method {
		return EndpointSpec{}, fmt.Errorf("invalid HTTP method %q: must be uppercase only", rawMethod)
	}

	if err := validateHTTPMethod(method); err != nil {
		return EndpointSpec{}, fmt.Errorf("invalid HTTP method in endpoint: %w", err)
	}

	path := match[2]
	if err := validateEndpointPath(path); err != nil {
		return EndpointSpec{}, fmt.Errorf("invalid endpoint path: %w", err)
	}

	endpoint := EndpointSpec{
		Path:     path,
		Method:   method,
		Request:  SchemaSpec{Fields: []FieldSpec{}},
		Response: SchemaSpec{Fields: []FieldSpec{}},
	}

	cs.parseFields(content, match[1], &endpoint)
	return endpoint, nil
}

// parseFields extracts request/response fields from markdown
func (cs *ContractSynthesizer) parseFields(content, method string, endpoint *EndpointSpec) {
	lines := strings.Split(content, "\n")
	inRequestSection := false
	inResponseSection := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, fmt.Sprintf("### %s", method)) {
			inRequestSection = false
			inResponseSection = false
			continue
		}

		if strings.Contains(line, "Request:") && strings.Contains(content, line) {
			inRequestSection = true
			inResponseSection = false
			cs.parseInlineFields(line, "Request:", &endpoint.Request)
			continue
		}

		if strings.Contains(line, "Response:") {
			inRequestSection = false
			inResponseSection = true
			cs.parseInlineFields(line, "Response:", &endpoint.Response)
			continue
		}

		if strings.HasPrefix(line, "-") {
			cs.parseBulletField(line, inRequestSection, inResponseSection, endpoint)
		}
	}
}

// parseInlineFields parses inline field specification {field1, field2}
func (cs *ContractSynthesizer) parseInlineFields(line, prefix string, schema *SchemaSpec) {
	pattern := prefix + `\s*\{([^}]{1,500})\}`
	inlineRe := regexp.MustCompile(pattern)
	if inlineMatches := inlineRe.FindStringSubmatch(line); len(inlineMatches) > 1 {
		fields := strings.Split(inlineMatches[1], ",")
		for _, f := range fields {
			f = strings.TrimSpace(f)
			if f != "" {
				sanitized, err := sanitizeFieldName(f)
				if err == nil {
					schema.Fields = append(schema.Fields, FieldSpec{
						Name:     sanitized,
						Type:     "string",
						Required: true,
					})
				}
			}
		}
	}
}

// parseBulletField parses bullet point field: - field_name: type
func (cs *ContractSynthesizer) parseBulletField(line string, inRequest, inResponse bool, endpoint *EndpointSpec) {
	fieldRe := regexp.MustCompile(`-\s*(\w{1,100}):\s*(\w{1,50})`)
	if fieldMatches := fieldRe.FindStringSubmatch(line); len(fieldMatches) > 2 {
		if len(endpoint.Request.Fields) >= MaxFieldCount || len(endpoint.Response.Fields) >= MaxFieldCount {
			return
		}

		fieldName := fieldMatches[1]
		sanitizedName, err := sanitizeFieldName(fieldName)
		if err != nil {
			return
		}

		field := FieldSpec{
			Name: sanitizedName,
			Type: fieldMatches[2],
		}

		if inRequest {
			endpoint.Request.Fields = append(endpoint.Request.Fields, field)
		} else if inResponse {
			endpoint.Response.Fields = append(endpoint.Response.Fields, field)
		}
	}
}
