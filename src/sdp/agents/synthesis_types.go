package agents

// ContractSynthesizer manages contract synthesis from requirements
type ContractSynthesizer struct {
	supervisor any // *synthesis.Supervisor - using any to avoid import cycle
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
	OpenAPI string    `yaml:"openapi"`
	Info    InfoSpec  `yaml:"info"`
	Paths   PathsSpec `yaml:"paths"`
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
	Summary     string        `yaml:"summary"`
	RequestBody *RequestSpec  `yaml:"requestBody,omitempty"`
	Responses   ResponsesSpec `yaml:"responses"`
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
	Description string               `yaml:"description"`
	Content     map[string]MediaSpec `yaml:"content"`
}

// MediaSpec represents OpenAPI media type
type MediaSpec struct {
	Schema SchemaRefSpec `yaml:"schema"`
}

// SchemaRefSpec represents OpenAPI schema reference
type SchemaRefSpec struct {
	Type       string                  `yaml:"type,omitempty"`
	Properties map[string]PropertySpec `yaml:"properties,omitempty"`
	Required   []string                `yaml:"required,omitempty"`
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
