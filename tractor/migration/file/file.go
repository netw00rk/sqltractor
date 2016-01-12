// Package file contains functions for low-level migration files handling.
package file

import (
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"regexp"
	"strconv"
	"strings"

	"github.com/netw00rk/sqltractor/tractor/migration/direction"
)

var filenameRegex = regexp.MustCompile(`^([0-9]+)_(.*)\.(up|down)\..*$`)

type ContentFunc func() ([]byte, error)

// File represents one file on disk.
// Example: 001_initial_plan_to_do_sth.up.sql
type File struct {
	// the name of the file
	FileName string

	// version parsed from filename
	Version uint64

	// the actual migration name parsed from filename
	Name string

	// function that reads file content
	ContentFunc ContentFunc

	// UP or DOWN migration
	Direction direction.Direction

	content []byte
}

func NewFile(fileName string, contentFunc ContentFunc) (*File, error) {
	version, name, d, err := parseFilenameSchema(fileName)
	if err != nil {
		return nil, err
	}

	return &File{
		FileName:    fileName,
		Version:     version,
		Name:        name,
		ContentFunc: contentFunc,
		Direction:   d,
		content:     nil,
	}, nil

}

// ReadContent reads the file's content if the content is empty
func (f *File) Content() ([]byte, error) {
	if len(f.content) == 0 {
		content, err := f.ContentFunc()
		if err != nil {
			return nil, err
		}
		f.content = content
	}
	return f.content, nil
}

// parseFilenameSchema parses the filename
func parseFilenameSchema(filename string) (version uint64, name string, d direction.Direction, err error) {
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
