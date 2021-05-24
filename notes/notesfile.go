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
func FromText(txt string) NotesFile {
	if txt == "" {
		return NotesFile{}
	}

	return NotesFile{Lines: strings.Split(txt, "\n")}
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

// Insert lines line containing specified substring
func (n *NotesFile) Insert(lines []string, position int) int {
	n.Lines = insertString(n.Lines, position, lines...)
	return position + len(lines)
}

// insertString inserts strings inside a string slice  at the specified position
// TODO: this is a generally useful function, i imagine i'll use it again, might pull it up in a different package
func insertString(dst []string, index int, s ...string) []string {
	if index >= len(dst) {
		return append(dst, s...)
	}

	if index <= 0 {
		return append(s, dst...)
	}

	a := make([]string, len(dst)+len(s)) // Allocate a new slice with new size
	copy(a, dst[:index])                 // Copy the first part of the original slice
	copy(a[index:], s)                   // Insert the strings at the specified position
	copy(a[index+len(s):], dst[index:])  // Copy the remaining of the original slice
	return a
}

// Find line containing specified substring
func (n *NotesFile) FindContains(s string) (int, bool) {
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
	os.MkdirAll(filepath.Dir(path), os.ModePerm)

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

func (n *NotesFile) String() string {
	return strings.Join(n.Lines, "\n")
}

func (n *NotesFile) StringWithLineNumbers() string {
	var lines []string = make([]string, len(n.Lines))
	for i, v := range n.Lines {
		lines[i] = fmt.Sprintf("%0*d. %s", (len(n.Lines)+1)%10, i, v)
	}

	return strings.Join(lines, "\n")
}
