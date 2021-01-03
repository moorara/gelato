package edit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/moorara/gelato/internal/log"
)

// Editor is used for editing text files.
type Editor struct {
	logger *log.ColorfulLogger
}

// NewEditor creates a new editor.
func NewEditor(level log.Level) *Editor {
	logger := log.NewColorful(level)

	return &Editor{
		logger: logger,
	}
}

// ReplaceSpec has the input parameters for the Replace method.
type ReplaceSpec struct {
	PathRE *regexp.Regexp
	OldRE  *regexp.Regexp
	New    string
}

// ReplaceInDir is used for modifying all files in a directory.
func (e *Editor) ReplaceInDir(root string, specs []ReplaceSpec) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			var data []byte

			for _, s := range specs {
				if s.PathRE.MatchString(path) {
					if data == nil {
						e.logger.Yellow.Tracef("Reading %s", path)
						e.logger.Green.Debugf("Editing %s", path)
						if data, err = ioutil.ReadFile(path); err != nil {
							return err
						}
					}

					e.logger.Magenta.Tracef("  Replacing %q with %q", s.OldRE, s.New)
					data = s.OldRE.ReplaceAll(data, []byte(s.New))
				}
			}

			if data != nil {
				e.logger.Yellow.Tracef("Writing back %s", path)
				if err := ioutil.WriteFile(path, data, 0); err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
