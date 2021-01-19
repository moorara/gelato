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

// Remove deletes files and folders using glob patterrns.
func (e *Editor) Remove(globs ...string) error {
	for _, glob := range globs {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return err
		}

		for _, match := range matches {
			if err := os.RemoveAll(match); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveSpec has the input parameters for the MoveFile method.
type MoveSpec struct {
	Src  string
	Dest string
}

// Move moves a file from a destination to a source
func (e *Editor) Move(specs ...MoveSpec) error {
	for _, s := range specs {
		if err := os.Rename(s.Src, s.Dest); err != nil {
			return err
		}
	}

	return nil
}

// ReplaceSpec has the input parameters for the Replace method.
type ReplaceSpec struct {
	PathRE *regexp.Regexp
	OldRE  *regexp.Regexp
	New    string
}

// ReplaceInDir is used for modifying all files in a directory.
func (e *Editor) ReplaceInDir(root string, specs ...ReplaceSpec) error {
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
