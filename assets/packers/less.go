// +build !js

package packers

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gu-io/gu/assets"
	"github.com/influx6/faux/process"
)

// LessPacker defines an implementation for parsing .less files into css files using the less compiler in nodejs.
// WARNING: Requires Nodejs to be installed.
type LessPacker struct {
	MainFile string
	Options  map[string]string
}

// Pack process all files present in the FileStatment slice and returns WriteDirectives
// which conta ins expected outputs for these files.
func (less LessPacker) Pack(statements []assets.FileStatement, dir assets.DirStatement) ([]assets.WriteDirective, error) {
	var directives []assets.WriteDirective

	// If main less file has being set then attempt to find main file.
	if less.MainFile == "" {
		for _, statement := range statements {
			if err := processStatement(statement, less, &directives); err != nil {
				return nil, err
			}
		}

		return directives, nil
	}

	for _, statement := range statements {
		if statement.Path != less.MainFile {
			continue
		}

		if err := processStatement(statement, less, &directives); err != nil {
			return nil, err
		}
	}

	return directives, nil
}

func processStatement(statement assets.FileStatement, less LessPacker, directives *[]assets.WriteDirective) error {
	fileExt := filepath.Ext(statement.Path)
	cssFileName := filepath.Join(filepath.Dir(statement.Path), strings.Replace(filepath.Base(statement.Path), fileExt, ".css", 1))
	cssAbsFileName := filepath.Join(filepath.Dir(statement.AbsPath), strings.Replace(filepath.Base(statement.Path), fileExt, ".css", 1))

	cssFileName = strings.Replace(cssFileName, "less/", "css/", 1)
	cssAbsFileName = strings.Replace(cssAbsFileName, "less/", "css/", 1)

	var args []string

	for option, value := range less.Options {
		args = append(args, option, value)
	}

	args = append(args, filepath.Clean(statement.AbsPath))

	cmd := process.Command{
		Args:  args,
		Name:  filepath.Join(guSrcNodeModulesBin, "lessc"),
		Level: process.RedAlert,
	}

	var errBuf, outBuf bytes.Buffer

	ctx, cancl := context.WithTimeout(context.Background(), time.Minute)
	defer cancl()

	if err := cmd.Run(ctx, &outBuf, &errBuf, nil); err != nil {
		return fmt.Errorf("Command Execution Failed: %+q\n Response: %+q", err, errBuf.String())
	}

	*directives = append(*directives, assets.WriteDirective{
		OriginPath:    cssFileName,
		OriginAbsPath: cssAbsFileName,
		Writer:        bytes.NewReader(outBuf.Bytes()),
	})

	return nil
}
