package aquamarine

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Spec struct {
	Version string `yaml:"version"`
	Project struct {
		Name   string `yaml:"name"`
		Module string `yaml:"module"`
	} `yaml:"project"`
	Runtime struct {
		HTTP struct {
			API struct {
				Host string `yaml:"host"`
				Port int    `yaml:"port"`
			} `yaml:"api"`
			Web struct {
				Host string `yaml:"host"`
				Port int    `yaml:"port"`
			} `yaml:"web"`
		} `yaml:"http"`
		Database struct {
			Engine string `yaml:"engine"`
			DSN    string `yaml:"dsn"`
		} `yaml:"database"`
	} `yaml:"runtime"`
	Feats []Feat `yaml:"feats"`
}

type Feat struct {
	Name  string `yaml:"name"`
	Kind  string `yaml:"kind"`
	Model struct {
		Name   string                    `yaml:"name"`
		Fields map[string]map[string]any `yaml:"fields"`
	} `yaml:"model"`
	Models map[string]struct {
		Fields map[string]map[string]any `yaml:"fields"`
	} `yaml:"models"`
	Service struct {
		Methods []string `yaml:"methods"`
	} `yaml:"service"`
	API struct {
		Routes []Route `yaml:"routes"`
	} `yaml:"api"`
	Pages []Page `yaml:"pages"`
}

type Route struct {
	Method  string `yaml:"method"`
	Path    string `yaml:"path"`
	Handler string `yaml:"handler"`
}

type Page struct {
	Route string   `yaml:"route"`
	Uses  []string `yaml:"uses"`
}

// Generate creates a minimal app skeleton under out/<mode> based on aquamarine.yaml.
// It does not overwrite files outside the out/ tree.
func Generate(assetsFS embed.FS, mode string) error {
	spec, err := readSpec("aquamarine.yaml")
	if err != nil {
		return err
	}
	outRoot := filepath.Join("out", mode)

	fmt.Printf("Generating aquamarine project '%s' in directory: %s\n", spec.Project.Name, outRoot)

	// Create base directory structure
	if err := os.MkdirAll(outRoot, 0o755); err != nil {
		return err
	}

	if err := writeGoMod(outRoot, spec.Project.Module); err != nil {
		return err
	}
	if err := writeAppMain(outRoot, spec); err != nil {
		return err
	}

	// Base dirs
	for _, d := range []string{
		"internal/platform",
		"internal/web",
		"internal/feat",
		"assets/templates",
		"assets/migrations/sqlite",
		"assets/seeds/sqlite",
		"assets/queries/sqlite",
	} {
		if err := os.MkdirAll(filepath.Join(outRoot, d), 0o755); err != nil {
			return err
		}
	}

	for _, f := range spec.Feats {
		featDir := filepath.Join(outRoot, "internal/feat", f.Name)
		if err := os.MkdirAll(featDir, 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(featDir, "README.md"), []byte("# "+f.Name+"\n\nGenerated feat skeleton."), 0o644); err != nil {
			return err
		}
		// TODO: move this to assets/templates/service.tmpl
		svc := fmt.Sprintf("package %s\n\n// Service methods (scaffold): %v\n", f.Name, f.Service.Methods)
		_ = os.WriteFile(filepath.Join(featDir, "service.go"), []byte(svc), 0o644)
	}

	assetsReadme := []byte("# assets\n\nCentralized assets for the generated app.\n")
	_ = os.WriteFile(filepath.Join(outRoot, "assets/README.md"), assetsReadme, 0o644)

	return nil
}

func readSpec(path string) (*Spec, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var s Spec
	if err := yaml.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	if s.Project.Module == "" {
		return nil, errors.New("project.module is required")
	}
	return &s, nil
}

func writeGoMod(outRoot, module string) error {
	// TODO: move this to assets/templates/go.mod.tmpl
	rel := "../.." // from out/<mode> to repo root
	content := fmt.Sprintf(`module %s

go 1.22

require github.com/aquamarinepk/aquamarine v0.0.0

replace github.com/aquamarinepk/aquamarine => %s
`, module, rel)
	return safeWriteFile(filepath.Join(outRoot, "go.mod"), []byte(content), 0o644)
}

func writeAppMain(outRoot string, cfg *Spec) error {
	// TODO: move this to assets/templates/main.tmpl
	main := fmt.Sprintf(`package main

import (
  "context"
  "embed"
  "log"
  "net/http"
  "os/signal"
  "syscall"

  "github.com/go-chi/chi/v5"
  "github.com/aquamarinepk/aquamarine/pkg/lib/am"
)

//go:embed assets
var assetsFS embed.FS

func main() {
  apiPort := ":8081"
  webPort := ":%d"

  logger := am.NewLogger("info")

  ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
  defer cancel()

  apiRouter := chi.NewRouter()
  webRouter := chi.NewRouter()
  
  apiRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("api ok"))
  })
  
  webRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte("web ok"))
  })

  var deps []any

  starts, stops := am.Setup(ctx, apiRouter, webRouter, deps...)

  if err := am.Start(ctx, starts, stops); err != nil {
    log.Fatal(err)
  }

  servers := []am.Server{
    {Name: "API", Addr: apiPort, Handler: apiRouter},
    {Name: "Web", Addr: webPort, Handler: webRouter},
  }

  am.StartServers(servers, logger)

  <-ctx.Done()
  
  am.GracefulShutdown(servers, stops, logger)
}
`, cfg.Runtime.HTTP.Web.Port)
	return safeWriteFile(filepath.Join(outRoot, "main.go"), []byte(main), 0o644)
}

func safeWriteFile(path string, data []byte, perm fs.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, perm)
}
