package main

import (
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/tamoore/esbx/pkg/plugins"
)

func main() {
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
		Outdir:     "build",
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
		EntryPoints:       []string{"src/app.jsx"},
		Stdin:             &api.StdinOptions{},
		Write:             true,
		Incremental:       false,
		Plugins:           []api.Plugin{plugins.HTTPModules},
	})

	if len(result.Errors) > 0 {
		log.Fatal(result.Errors)
		os.Exit(1)
	}
}
