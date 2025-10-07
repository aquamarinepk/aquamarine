package aquamarine

type Config struct {
	Version    string             `yaml:"version"`
	Project    ProjectConfig      `yaml:"project"`
	Runtime    RuntimeConfig      `yaml:"runtime,omitempty"`
	Feats      map[string]Feature `yaml:"feats"`
	ModulePath string             `yaml:"-"` // Set during generation
}

// ProjectConfig contains project-level configuration.
type ProjectConfig struct {
	Name   string `yaml:"name"`
	Module string `yaml:"module"`
}

// RuntimeConfig contains runtime configuration for the generated application.
type RuntimeConfig struct {
	HTTP     HTTPConfig     `yaml:"http,omitempty"`
	Database DatabaseConfig `yaml:"database,omitempty"`
}

// HTTPConfig contains HTTP server configuration.
type HTTPConfig struct {
	API APIConfig `yaml:"api,omitempty"`
	Web WebConfig `yaml:"web,omitempty"`
}

// APIConfig contains API server configuration.
type APIConfig struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

// WebConfig contains web server configuration.
type WebConfig struct {
	Host string `yaml:"host,omitempty"`
	Port int    `yaml:"port,omitempty"`
}

// DatabaseConfig contains database configuration.
type DatabaseConfig struct {
	Engine string `yaml:"engine,omitempty"`
}

type Feature struct {
	Name       string               `yaml:"name,omitempty"`
	Kind       string               `yaml:"kind,omitempty"` // "feature", "web", "support"
	Models     map[string]Model     `yaml:"models,omitempty"`
	Service    ServiceConfig        `yaml:"service,omitempty"`
	API        APIFeatureConfig     `yaml:"api,omitempty"`
	Web        WebFeatureConfig     `yaml:"web,omitempty"`
	RepoImpl   []string             `yaml:"repo_impl,omitempty"`
	Auth       *AuthConfig          `yaml:"auth,omitempty"`
	Aggregates map[string]Aggregate `yaml:"aggregates,omitempty"`
}

// ServiceConfig contains service-level configuration.
type ServiceConfig struct {
	Methods []string `yaml:"methods,omitempty"`
}

// APIFeatureConfig contains API configuration for a feature.
type APIFeatureConfig struct {
	Routes []RouteConfig `yaml:"routes,omitempty"`
}

// RouteConfig represents an API route configuration.
type RouteConfig struct {
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Handler string `yaml:"handler"`
}

// WebFeatureConfig contains web configuration for a feature.
type WebFeatureConfig struct {
	Pages []PageConfig `yaml:"pages,omitempty"`
}

// PageConfig represents a web page configuration.
type PageConfig struct {
	Route string   `yaml:"route"`
	Uses  []string `yaml:"uses,omitempty"`
}

// Model represents a domain model.
type Model struct {
	Fields  map[string]Field `yaml:"fields,omitempty"`
	Options *ModelOptions    `yaml:"options,omitempty"`
}

// Field represents a model field.
type Field struct {
	Type        string       `yaml:"type"`
	Validations []Validation `yaml:"validations,omitempty"`
}

// Validation represents a field validation rule.
type Validation struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value,omitempty"`
}

// ModelOptions contains optional model configuration.
type ModelOptions struct {
	Audit bool `yaml:"audit,omitempty"`
}

// AuthConfig contains authentication configuration.
type AuthConfig struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

// Aggregate represents an aggregate root with child collections.
type Aggregate struct {
	Fields       map[string]Field       `yaml:"fields,omitempty"`
	VersionField string                 `yaml:"version_field,omitempty"`
	Audit        bool                   `yaml:"audit,omitempty"`
	Children     map[string]ChildConfig `yaml:"children,omitempty"`
}

// ChildConfig represents a child collection configuration.
type ChildConfig struct {
	Of    string `yaml:"of"` // References a model name
	Audit bool   `yaml:"audit,omitempty"`
}
