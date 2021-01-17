package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tamoore/esbx/internal/cmd"
	"github.com/tamoore/esbx/pkg/plugins"
)

var entrypoints cmd.StringList
var outdir string

func main() {
	flag.Var(&entrypoints, "entrypoint", "entrypoints to build")
	flag.StringVar(&outdir, "outdir", "build", "the directory to output to")
	flag.Parse()

	if len(entrypoints) == 0 {
		fmt.Fprintf(os.Stderr, "esbx: at least one entrypoint is required: recv %v\n", entrypoints)
	}

	result := api.Build(api.BuildOptions{
		Color:             0,
		ErrorLimit:        0,
		LogLevel:          0,
		Sourcemap:         0,
		SourcesContent:    0,
		Target:            0,
		Engines:           []api.Engine{},
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Charset:           0,
		TreeShaking:       0,
		JSXFactory:        "",
		JSXFragment:       "",
		Define: map[string]string{
			"process.env.NODE_ENV": "\"production\"",
		},
		Pure:       []string{},
		AvoidTDZ:   false,
		KeepNames:  false,
		GlobalName: "",
		Bundle:     true,
		Splitting:  false,
		Outfile:    "",
		Metafile:   "",
		Outdir:     outdir,
		Outbase:    "",
		Platform:   0,
		Format:     0,
		External:   []string{},
		MainFields: []string{},
		Loader: map[string]api.Loader{
			".html": api.LoaderText,
			".css":  api.LoaderText,
		},
		ResolveExtensions: []string{},
		Tsconfig:          "",
		OutExtensions:     map[string]string{},
		PublicPath:        "",
		Inject:            []string{},
		Banner:            "",
		Footer:            "",
		EntryPoints:       []string(entrypoints),
		Stdin:             &api.StdinOptions{},
		Write:             true,
		Incremental:       false,
		Plugins:           []api.Plugin{plugins.HTTPModules},
	})

	if len(result.Errors) > 0 {
		fmt.Fprintf(os.Stderr, "esbx: issuing building: %v\n", result.Errors)
		os.Exit(1)
	}
}
