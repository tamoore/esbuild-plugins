package plugins

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/evanw/esbuild/pkg/api"
)

var cacheDIR = filepath.Join(xdg.CacheHome, "esbuild")

// HTTPModules provides the http.Plugin to be used with esbuild plugins
var HTTPModules = api.Plugin{
	Name: "http",
	Setup: func(build api.PluginBuild) {
		// Intercept import paths starting with "http:" and "https:" so
		// esbuild doesn't attempt to map them to a file system location.
		// Tag them with the "http-url" namespace to associate them with
		// this plugin.
		build.OnResolve(api.OnResolveOptions{Filter: `^https?://`},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{
					Path:      args.Path,
					Namespace: "http-url",
				}, nil
			})

		// We also want to intercept all import paths inside downloaded
		// files and resolve them against the original URL. All of these
		// files will be in the "http-url" namespace. Make sure to keep
		// the newly resolved URL in the "http-url" namespace so imports
		// inside it will also be resolved as URLs recursively.
		build.OnResolve(api.OnResolveOptions{Filter: ".*", Namespace: "http-url"},
			func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				base, err := url.Parse(args.Importer)
				if err != nil {
					return api.OnResolveResult{}, err
				}
				relative, err := url.Parse(args.Path)
				if err != nil {
					return api.OnResolveResult{}, err
				}
				return api.OnResolveResult{
					Path:      base.ResolveReference(relative).String(),
					Namespace: "http-url",
				}, nil
			})

		// When a URL is loaded, we want to actually download the content
		// from the internet. This has just enough logic to be able to
		// handle the example import from unpkg.com but in reality this
		// would probably need to be more complex.
		build.OnLoad(api.OnLoadOptions{Filter: ".*", Namespace: "http-url"},
			func(args api.OnLoadArgs) (api.OnLoadResult, error) {
				var contents string
				sum := sha256.Sum256([]byte(args.Path))
				filePath := filepath.Join(cacheDIR, fmt.Sprintf("%x", sum))

				// Create the cache directory if it doesn't exist
				if _, err := os.Stat(cacheDIR); os.IsNotExist(err) {
					os.MkdirAll(cacheDIR, 0700)
				}

				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					// Download the esm package from URL
					res, err := http.Get(args.Path)

					fmt.Printf("Downloading: %s\n", args.Path)

					if err != nil {
						return api.OnLoadResult{}, err
					}
					defer res.Body.Close()
					bytes, err := ioutil.ReadAll(res.Body)
					if err != nil {
						return api.OnLoadResult{}, err
					}
					contents = string(bytes)
					ioutil.WriteFile(filePath, bytes, 0655)

				} else {
					// Read from cache

					bytes, err := ioutil.ReadFile(filePath)
					if err != nil {
						return api.OnLoadResult{}, err
					}
					contents = string(bytes)
				}

				return api.OnLoadResult{Contents: &contents}, nil
			})
	},
}
