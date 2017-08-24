//+build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gu-io/gu/shell"
	"github.com/gu-io/gu/shell/parse"
)

func main() {
	rootDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	pkgDir := filepath.Join(rootDir, "detos")
	res, err := parse.ShellResources(pkgDir)
	if err != nil {
		panic(err)
	}

	var manifests []*shell.AppManifest

	for _, rs := range res {
		ms, err := rs.GenManifests()
		if err != nil {
			panic(err)
		}

		manifests = append(manifests, ms)
	}

	jsn, _ := json.MarshalIndent(manifests, "", "\t")
	fmt.Printf("Manifest(Package: %s):\n%s\n", pkgDir, jsn)
}
