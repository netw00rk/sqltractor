// Package file contains functions for low-level migration files handling.
package file

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/netw00rk/sqltractor/tractor/direction"
	"github.com/netw00rk/sqltractor/tractor/file/reader"
)


// File represents one file on disk.
// Example: 001_initial_plan_to_do_sth.up.sql
type File struct {
	// absolute path to file
	Path string

	// the name of the file
	FileName string

	// version parsed from filename
	Version uint64

	// the actual migration name parsed from filename
	Name string

	// content of the file
	Content []byte

	// UP or DOWN migration
	Direction direction.Direction
}

// ReadContent reads the file's content if the content is empty
func (f *File) ReadContent() error {
	if len(f.Content) == 0 {
		content, err := defaultFileReader.ReadFileContent(path.Join(f.Path, f.FileName))
		if err != nil {
			return err
		}
		f.Content = content
	}
	return nil
}

// parseFilenameSchema parses the filename
func parseFilenameSchema(filename string, filenameRegex *regexp.Regexp) (version uint64, name string, d direction.Direction, err error) {
	matches := filenameRegex.FindStringSubmatch(filename)
	if len(matches) != 4 {
		return 0, "", 0, errors.New("Unable to parse filename schema")
	}

	version, err = strconv.ParseUint(matches[1], 10, 0)
	if err != nil {
		return 0, "", 0, errors.New(fmt.Sprintf("Unable to parse version '%v' in filename schema", matches[0]))
	}

	if matches[3] == "up" {
		d = direction.Up
	} else if matches[3] == "down" {
		d = direction.Down
	} else {
		return 0, "", 0, errors.New(fmt.Sprintf("Unable to parse up|down '%v' in filename schema", matches[3]))
	}

	return version, matches[2], d, nil
}

// LineColumnFromOffset reads data and returns line and column integer
// for a given offset.
func LineColumnFromOffset(data []byte, offset int) (line, column int) {
	// TODO is there a better way?
	fs := token.NewFileSet()
	tf := fs.AddFile("", fs.Base(), len(data))
	tf.SetLinesForContent(data)
	pos := tf.Position(tf.Pos(offset))
	return pos.Line, pos.Column
}

// LinesBeforeAndAfter reads n lines before and after a given line.
// Set lineNumbers to true, to prepend line numbers.
func LinesBeforeAndAfter(data []byte, line, before, after int, lineNumbers bool) []byte {
	// TODO(mattes): Trim empty lines at the beginning and at the end
	// TODO(mattes): Trim offset whitespace at the beginning of each line, so that indentation is preserved
	startLine := line - before
	endLine := line + after
	lines := bytes.SplitN(data, []byte("\n"), endLine+1)

	if startLine < 0 {
		startLine = 0
	}
	if endLine > len(lines) {
		endLine = len(lines)
	}

	selectLines := lines[startLine:endLine]
	newLines := make([][]byte, 0)
	lineCounter := startLine + 1
	lineNumberDigits := len(strconv.Itoa(len(selectLines)))
	for _, l := range selectLines {
		lineCounterStr := strconv.Itoa(lineCounter)
		if len(lineCounterStr)%lineNumberDigits != 0 {
			lineCounterStr = strings.Repeat(" ", lineNumberDigits-len(lineCounterStr)%lineNumberDigits) + lineCounterStr
		}

		lNew := l
		if lineNumbers {
			lNew = append([]byte(lineCounterStr+": "), lNew...)
		}
		newLines = append(newLines, lNew)
		lineCounter += 1
	}

	return bytes.Join(newLines, []byte("\n"))
}


func init() {
	SetFileReader(reader.IOFileReader{})
}