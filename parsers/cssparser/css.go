package cssparser

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CSSItems defines a map of CSSItem objects which define selected css files and their
// respective contents and connections.
type CSSItems struct {
	Styles  []CSSItem        `json:"styles"`
	Indexes map[string]int   `json:"indexes"`
	Dirs    map[string][]int `json:"dirs"`
}

// CSSItem defines a struct which is returned when passing a directories of css files
// marked by .css extensions. It gleens out inclusion directives by scanner the top of
// the files data for any '/* #include "boxer.css",...,"[relative_file_paths]" */', which
// will indicate which other files should be included with the giving css item.
type CSSItem struct {
	Dir      string   `json:"dir"`
	Path     string   `json:"path"`
	Content  []byte   `json:"content"`
	Includes []string `json:"includes"`
}

// ParseDir returns a new instance of all CSS files located within the provided directory.
func ParseDir(dir string) (*CSSItems, error) {
	var items CSSItems
	items.Styles = make([]CSSItem, 0)
	items.Dirs = make(map[string][]int)
	items.Indexes = make(map[string]int)

	// Walk directory pulling contents into css items.
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err := walkDir(&items, dir, path, info, err); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	// Run through all dirs and find related subdirectories.
	for key, subs := range items.Dirs {

		for _, item := range items.Styles {

			if !strings.HasPrefix(item.Dir, key) {
				continue
			}

			subs = append(subs, items.Indexes[item.Path])
		}

		items.Dirs[key] = subs
	}

	return &items, nil
}

var allowedExtensions = []string{".css", ".gcss"}

// validExension returns true/false if the extension provide is a valid acceptable one
// based on the allowedExtensions string slice.
func validExtension(ext string) bool {
	for _, es := range allowedExtensions {
		if es != ext {
			continue
		}

		return true
	}

	return false
}

// walkDir adds the giving path if it matches certain criterias into the items map.
func walkDir(items *CSSItems, root string, path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if !info.Mode().IsRegular() {
		return nil
	}

	// Is file an exension we allow else skip.
	if !validExtension(filepath.Ext(path)) {
		return nil
	}

	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}

	baseDir := filepath.Dir(rel)
	if baseDir == "" {
		baseDir = "/"
	}

	items.Dirs[baseDir] = make([]int, 0)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	var item CSSItem
	item.Path = rel
	item.Dir = filepath.Dir(rel)
	item.Content = data

	reader := bytes.NewReader(data)
	bufReader := bufio.NewReader(reader)

	if line, err := bufReader.ReadString('\n'); err == nil {
		if strings.Contains(line, "#include") {
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "/*")
			line = strings.TrimSuffix(line, "*/")
			line = strings.TrimSpace(line)
			line = strings.TrimPrefix(line, "#include")
			line = strings.TrimSpace(line)

			files := strings.SplitN(line, ",", -1)
			for _, file := range files {
				item.Includes = append(item.Includes, strings.TrimSpace(file))
			}
		}
	}

	styleLen := len(items.Styles)
	items.Indexes[rel] = styleLen
	items.Styles = append(items.Styles, item)

	return nil
}
