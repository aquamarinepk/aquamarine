package aquamarine

import (
	"embed"
	"errors"
	"flag"
	"fmt"
)

// App is the CLI entrypoint for the Aquamarine generator.
type App struct {
	assets embed.FS
}

func NewApp(assetsFS embed.FS) *App {
	return &App{assets: assetsFS}
}

// Run parses CLI args and dispatches to subcommands.
// Supported (scaffold):
func (a *App) Run(args []string) error {
	if len(args) < 2 {
		usage()
		return nil
	}
	sub := args[1]
	switch sub {
	case "generate":
		fs := flag.NewFlagSet("generate", flag.ContinueOnError)
		dev := fs.Bool("dev", false, "generate into development output directory")
		if err := fs.Parse(args[2:]); err != nil {
			return err
		}
		mode := "prod"
		if *dev {
			mode = "dev"
		}
		if err := Generate(a.assets, mode); err != nil {
			return err
		}
		return nil
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		usage()
		return errors.New("unknown command")
	}
}

func usage() {
	fmt.Println("Aquamarine generator")
	fmt.Println("Usage:")
	fmt.Println("  aquamarine generate [--dev]")
	fmt.Println("  aquamarine help")
}
