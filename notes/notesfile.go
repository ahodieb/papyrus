package notes

import (
	"bufio"
	"os"
	"strings"
)

type NotesFile struct {
	Path  string
	Lines []string
}

// Read all lines in file if exists, or create a new empty one
func ReadOrCreate(path string) (NotesFile, error) {
	err := createIfIsNotExist(path)
	if err != nil {
		return NotesFile{}, err
	}

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return NotesFile{}, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return NotesFile{Path: path, Lines: lines}, nil
}

// Find line containing specified substring
func (n *NotesFile) FindContaining(s string) (int, bool) {
	var finder LineFinder = func(line string) bool { return strings.Contains(line, s) }
	return n.Find(finder)
}

type LineFinder func(line string) bool

// Find line matched by specified finder
func (n *NotesFile) Find(finder LineFinder) (int, bool) {
	for i, line := range n.Lines {
		if finder(line) {
			return i, true
		}
	}

	return 0, false
}

func createIfIsNotExist(p string) error {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		_, err := os.Create(p)
		return err
	}

	return err
}
