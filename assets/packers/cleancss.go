// +build !js

package packers

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gu-io/gu/assets"
	"github.com/influx6/faux/process"
)

var (
	inGOPATH            = os.Getenv("GOPATH")
	inGOPATHSrc         = filepath.Join(inGOPATH, "src")
	guSrc               = filepath.Join(inGOPATHSrc, "github.com/gu-io/gu")
	guSrcNodeModules    = filepath.Join(inGOPATHSrc, "github.com/gu-io/gu/node_modules")
	guSrcNodeModulesBin = filepath.Join(inGOPATHSrc, "github.com/gu-io/gu/node_modules/.bin/")
)

// CleanCSSPacker defines an implementation for parsing css files.
// WARNING: Requires Nodejs to be installed.
type CleanCSSPacker struct {
	Options map[string]string
}

// Pack process all files present in the FileStatment slice and returns WriteDirectives
// which conta ins expected outputs for these files.
func (cess CleanCSSPacker) Pack(statements []assets.FileStatement, dir assets.DirStatement) ([]assets.WriteDirective, error) {
	var directives []assets.WriteDirective

	for _, statement := range statements {
		if err := processCleanStatement(statement, cess, &directives); err != nil {
			return nil, err
		}
	}

	return directives, nil
}

func processCleanStatement(statement assets.FileStatement, cess CleanCSSPacker, directives *[]assets.WriteDirective) error {
	var args []string

	for option, value := range cess.Options {
		args = append(args, option, value)
	}

	args = append(args, filepath.Clean(statement.AbsPath))

	cmd := process.Command{
		Args:  args,
		Name:  filepath.Join(guSrcNodeModulesBin, "cleancss"),
		Level: process.RedAlert,
	}

	var errBuf, outBuf bytes.Buffer

	ctx, cancl := context.WithTimeout(context.Background(), time.Minute)
	defer cancl()

	if err := cmd.Run(ctx, &outBuf, &errBuf, nil); err != nil {
		return fmt.Errorf("Command Execution Failed: %+q\n Response: %+q", err, errBuf.String())
	}

	*directives = append(*directives, assets.WriteDirective{
		OriginPath:    statement.Path,
		OriginAbsPath: statement.AbsPath,
		Writer:        bytes.NewReader(outBuf.Bytes()),
	})

	return nil
}
