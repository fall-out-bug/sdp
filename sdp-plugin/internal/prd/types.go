package prd

// ProjectType represents different types of projects
type ProjectType string

const (
	// Service - Web service/API with Docker
	Service ProjectType = "service"
	// CLI - Command-line interface tool
	CLI ProjectType = "cli"
	// Library - Code library/package
	Library ProjectType = "library"
	// Go - Go project/module
	Go ProjectType = "go"
	// Python - Python project
	Python ProjectType = "python"
	// Java - Java project
	Java ProjectType = "java"
	// Unknown - Unknown project type
	Unknown ProjectType = "unknown"
)

// String returns the string representation
func (p ProjectType) String() string {
	return string(p)
}

// Value returns the enum value (for compatibility)
func (p ProjectType) Value() string {
	return string(p)
}
