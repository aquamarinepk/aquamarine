package main

import (
	"embed"
	"fmt"
	"os"

	"github.com/aquamarinepk/aquamarine/internal/aquamarine"
)

//go:embed assets
var assetsFS embed.FS

func main() {
	app := aquamarine.NewApp(assetsFS)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
