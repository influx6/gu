// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gu-io/gu/parsers/cssparser"
)

func main() {
	cdir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	items, err := cssparser.ParseDir(filepath.Join(cdir, "./"))
	if err != nil {
		panic(err)
	}

	jd, err := json.MarshalIndent(items, "", "\t")
	if err != nil {
		panic(err)
	}

	fmt.Printf("CSS: %+s\n", jd)
}
