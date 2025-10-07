package aquamarine

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

// FieldTemplateData holds data for a single field in the template.
type FieldTemplateData struct {
	Name        string
	Type        string
	JSONTag     string
	IsID        bool
	Validations []FieldValidationData
}

// FieldValidationData holds data for a single validation rule.
type FieldValidationData struct {
	Name  string
	Value string
}

// ModelTemplateData holds all data needed to render a model template.
type ModelTemplateData struct {
	PackageName  string
	ModelName    string
	Audit        bool
	Fields       []FieldTemplateData
	NeedsFmt     bool
	NeedsStrconv bool
}

// HandlerTemplateData holds all data needed to render a handler template.
type HandlerTemplateData struct {
	PackageName       string
	ModelName         string
	ModelPlural       string
	ModelLower        string
	ModelPluralLower  string
	AuthEnabled       bool
	Audit             bool
	ModulePath        string
	IsChildCollection bool
}

type FeatureGenerator struct {
	Config                   Config
	OutputDir                string
	DevMode                  bool
	Template                 *template.Template
	RepoInterfaceTemplate    *template.Template
	ServiceInterfaceTemplate *template.Template
	SQLiteRepoTemplate       *template.Template
	SQLiteQueriesTemplate    *template.Template
	MongoRepoTemplate        *template.Template
	HandlerTemplate          *template.Template
	ValidatorTemplate        *template.Template
	MainTemplate             *template.Template
	ConfigTemplate           *template.Template
	ConfigYAMLTemplate       *template.Template
	XParamsTemplate          *template.Template
	MakefileTemplate         *template.Template
	AggregateRootTemplate    *template.Template
	ChildCollectionTemplate  *template.Template
}

// NewFeatureGenerator creates a new feature generator.
func NewFeatureGenerator(config Config, outputDir string, devMode bool, assetsFS fs.FS) (*FeatureGenerator, error) {
	tmplFS, err := fs.Sub(assetsFS, "assets/templates")
	if err != nil {
		log.Fatalf("cannot create sub-filesystem for templates: %v", err)
	}

	tmpl, err := template.New("model.tmpl").ParseFS(tmplFS, "model.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse model template: %w", err)
	}

	repoInterfaceTmpl, err := template.New("repo_interface.tmpl").ParseFS(tmplFS, "repo_interface.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse repository interface template: %w", err)
	}

	serviceInterfaceTmpl, err := template.New("service_interface.tmpl").ParseFS(tmplFS, "service_interface.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse service interface template: %w", err)
	}

	sqliteRepoTmpl, err := template.New("repo_sqlite.tmpl").ParseFS(tmplFS, "repo_sqlite.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse SQLite repository template: %w", err)
	}

	sqliteQueriesTmpl, err := template.New("queries_sqlite.tmpl").ParseFS(tmplFS, "queries_sqlite.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse SQLite queries template: %w", err)
	}

	mongoRepoTmpl, err := template.New("repo_mongo.tmpl").ParseFS(tmplFS, "repo_mongo.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse MongoDB repository template: %w", err)
	}

	handlerTmpl, err := template.New("handler.tmpl").ParseFS(tmplFS, "handler.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse handler template: %w", err)
	}

	validatorTmpl, err := template.New("validator.tmpl").ParseFS(tmplFS, "validator.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse validator template: %w", err)
	}

	mainTmpl, err := template.New("main.tmpl").ParseFS(tmplFS, "main.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse main template: %w", err)
	}

	configTmpl, err := template.New("config.go.tmpl").ParseFS(tmplFS, "config.go.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse config go template: %w", err)
	}

	configYAMLTmpl, err := template.New("config.yaml.tmpl").ParseFS(tmplFS, "config.yaml.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse config yaml template: %w", err)
	}

	xparamsTmpl, err := template.New("xparams.go.tmpl").ParseFS(tmplFS, "xparams.go.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse xparams template: %w", err)
	}

	makefileTmpl, err := template.New("Makefile.tmpl").ParseFS(tmplFS, "Makefile.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse Makefile template: %w", err)
	}

	aggregateRootTmpl, err := template.New("aggregate_root.tmpl").ParseFS(tmplFS, "aggregate_root.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse aggregate root template: %w", err)
	}

	childCollectionTmpl, err := template.New("child_collection.tmpl").ParseFS(tmplFS, "child_collection.tmpl")
	if err != nil {
		return nil, fmt.Errorf("cannot parse child collection template: %w", err)
	}

	return &FeatureGenerator{
		Config:                   config,
		OutputDir:                outputDir,
		DevMode:                  devMode,
		Template:                 tmpl,
		RepoInterfaceTemplate:    repoInterfaceTmpl,
		ServiceInterfaceTemplate: serviceInterfaceTmpl,
		SQLiteRepoTemplate:       sqliteRepoTmpl,
		SQLiteQueriesTemplate:    sqliteQueriesTmpl,
		MongoRepoTemplate:        mongoRepoTmpl,
		HandlerTemplate:          handlerTmpl,
		ValidatorTemplate:        validatorTmpl,
		MainTemplate:             mainTmpl,
		ConfigTemplate:           configTmpl,
		ConfigYAMLTemplate:       configYAMLTmpl,
		XParamsTemplate:          xparamsTmpl,
		MakefileTemplate:         makefileTmpl,
		AggregateRootTemplate:    aggregateRootTmpl,
		ChildCollectionTemplate:  childCollectionTmpl,
	}, nil
}

func (fg *FeatureGenerator) GenerateModels() error {
	for featName, feat := range fg.Config.Feats {
		for modelName, model := range feat.Models {
			if model.Fields == nil {
				model.Fields = make(map[string]Field)
			}
			fmt.Printf("  - Generating model: %s/%s\n", featName, modelName)

			packageName := featName
			// Individual model files: user.go, list.go, order.go, etc...
			modelFileName := strings.ToLower(modelName) + ".go"
			// Path: internal/feat/{featName}/{model}.go
			modelPath := filepath.Join(fg.OutputDir, "internal", "feat", featName, modelFileName)

			data := ModelTemplateData{
				PackageName: packageName,
				ModelName:   modelName,
				Audit:       false,
				Fields:      []FieldTemplateData{},
			}
			if model.Options != nil {
				data.Audit = model.Options.Audit
			}

			var needsFmt bool
			var needsStrconv bool

			for fieldName, field := range model.Fields {
				goType := mapGoType(field.Type)
				fieldData := FieldTemplateData{
					Name:        capitalizeFirst(fieldName),
					Type:        goType,
					JSONTag:     toSnakeCase(fieldName),
					IsID:        false,
					Validations: []FieldValidationData{},
				}

				for _, v := range field.Validations {
					fieldData.Validations = append(fieldData.Validations, FieldValidationData{
						Name:  v.Name,
						Value: v.Value,
					})
					switch v.Name {
					case "min_length", "max_length", "min", "max":
						needsFmt = true
						needsStrconv = true
					}
				}
				data.Fields = append(data.Fields, fieldData)
			}

			data.NeedsFmt = needsFmt
			data.NeedsStrconv = needsStrconv

			dir := filepath.Dir(modelPath)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("cannot create directory for model %s: %w", modelName, err)
			}

			file, err := os.Create(modelPath)
			if err != nil {
				return fmt.Errorf("cannot create model file %s: %w", modelPath, err)
			}
			defer file.Close()

			if err := fg.Template.Execute(file, data); err != nil {
				return fmt.Errorf("cannot execute model template for %s: %w", modelName, err)
			}
			fmt.Printf("    - Created %s\n", modelPath)
		}
	}
	return nil
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func capitalizeFirst(str string) string {
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// mapGoType maps YAML types to Go types.
func mapGoType(yamlType string) string {
	switch yamlType {
	case "text", "string", "email":
		return "string"
	case "bool":
		return "bool"
	case "uuid":
		return "uuid.UUID"
	case "int":
		return "int"
	case "int64":
		return "int64"
	case "float64":
		return "float64"
	default:
		return "any"
	}
}
