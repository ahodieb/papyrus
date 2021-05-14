package notes

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type NotesFile struct {
	Path  string
	Lines []string
}

// Read notes file from path
func ReadFromFile(path string) NotesFile {
	return NotesFile{
		Path:  path,
		Lines: readLines(path),
	}
}

func readLines(path string) []string {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var lines []string
	sc := bufio.NewScanner(bytes.NewReader(b))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}

	return lines
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

// Save the current lines into Notes file
// creates a backup file before overwriting the existing file
// returns path to backup file and an error
func (n *NotesFile) SaveWithBackup() (string, error) {
	bkp, err := n.Backup()
	if err != nil {
		return "", err
	}

	err = n.Save()
	if err != nil {
		return "", err
	}

	return bkp, nil
}

// Saves lines to note file, this overwrites the existing file
// unless intentionally required use SaveWithBackup instead for more safety
func (n *NotesFile) Save() error {
	f, err := os.Create(n.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()
	for _, line := range n.Lines {
		_, err = w.WriteString(line + "\n")
	}

	return err
}

// Backup the Notes file, returns path for new backup file
// The backup file will be next to the notes file under ".backups/YYYYmmDD-<HHMM>.txt"
func (n *NotesFile) Backup() (string, error) {
	source, err := os.ReadFile(n.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		} else {
			return "", fmt.Errorf("Backup failed, could not read original file: %s, %s", n.Path, err)
		}
	}

	path := path.Join(filepath.Dir(n.Path), ".backups", time.Now().UTC().Format("20060102-0304.txt"))
	os.MkdirAll(path, os.ModePerm)
	destination, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("Backup failed, could not create backup file: %s, %s", path, err)
	}
	defer destination.Close()

	_, err = destination.Write(source)
	if err != nil {
		return "", fmt.Errorf("Backup failed, could not write to backup file: %s, %s", path, err)
	}

	return path, nil
}
