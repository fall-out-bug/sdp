package agents

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fall-out-bug/sdp/src/sdp/synthesis"
	"gopkg.in/yaml.v3"
)

// NewContractSynthesizer creates a new contract synthesizer
func NewContractSynthesizer() *ContractSynthesizer {
	engine := synthesis.DefaultRuleEngine()
	supervisor := synthesis.NewSupervisor(engine, 3)

	return &ContractSynthesizer{
		supervisor: supervisor,
	}
}

// AnalyzeRequirements parses a requirements markdown file
func (cs *ContractSynthesizer) AnalyzeRequirements(reqPath string) (*ContractRequirements, error) {
	content, err := os.ReadFile(reqPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read requirements: %w", err)
	}

	featureName := strings.TrimSuffix(filepath.Base(reqPath), "-requirements.md")
	featureName = strings.TrimPrefix(featureName, "sdp-")

	if err := validateFeatureName(featureName); err != nil {
		return nil, fmt.Errorf("invalid feature name in path %q: %w", reqPath, err)
	}

	endpoints, err := cs.parseEndpointsFromMarkdown(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse endpoints: %w", err)
	}

	return &ContractRequirements{
		FeatureName: featureName,
		Endpoints:   endpoints,
	}, nil
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

// ApplySynthesisRules applies synthesis rules to resolve conflicts
func (cs *ContractSynthesizer) ApplySynthesisRules(proposals []*synthesis.Proposal) (*synthesis.SynthesisResult, error) {
	synthesizer := synthesis.NewSynthesizer()

	for _, proposal := range proposals {
		synthesizer.AddProposal(proposal)
	}

	return synthesizer.Synthesize()
}

// WriteContract writes the agreed contract to a YAML file
func (cs *ContractSynthesizer) WriteContract(contract *OpenAPIContract, outputPath string) error {
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := yaml.Marshal(contract)
	if err != nil {
		return fmt.Errorf("failed to marshal contract: %w", err)
	}

	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write contract: %w", err)
	}

	return nil
}

// SynthesizeContract performs end-to-end contract synthesis
func (cs *ContractSynthesizer) SynthesizeContract(featureName, reqPath, outputPath string) (*synthesis.SynthesisResult, error) {
	requirements, err := cs.AnalyzeRequirements(reqPath)
	if err != nil {
		return nil, fmt.Errorf("analyze requirements failed: %w", err)
	}

	contract, err := cs.ProposeContract(requirements)
	if err != nil {
		return nil, fmt.Errorf("propose contract failed: %w", err)
	}

	proposals := []*synthesis.Proposal{
		synthesis.NewProposal("architect", contract, 1.0, "Initial contract from requirements"),
	}

	result, err := cs.ApplySynthesisRules(proposals)
	if err != nil {
		return nil, fmt.Errorf("apply synthesis rules failed: %w", err)
	}

	finalContract := result.Solution.(*OpenAPIContract)
	if err := cs.WriteContract(finalContract, outputPath); err != nil {
		return nil, fmt.Errorf("write contract failed: %w", err)
	}

	return result, nil
}
