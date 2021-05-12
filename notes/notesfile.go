package notes

import (
	"bufio"
	"os"
	"time"
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

func createIfIsNotExist(p string) error {
	_, err := os.Stat(p)
	if os.IsNotExist(err) {
		_, err := os.Create(p)
		return err
	}

	return err
}

// Find the latest entry in the Journal file
func (f *NotesFile) FindPosition(t time.Time) int {

	// Looks for the latest entry before the current day entry, and get the position two lines above it
	position, found := f.positionBefore(t)
	if found {
		if position > 1 {
			return position - 1
		}

		return 0
	}

	// If no later entry found look for the current day entry,
	// This could happen if there is only one entry in the file, or change of formats
	// If none were found that means either its an empty file or formats are not recognized
	// and it will default back to the 0th position
	position, _ = f.positionOn(t)
	return position
}

func (f *NotesFile) positionOn(t time.Time) (int, bool) {
	return 3, true
}

func (f *NotesFile) positionBefore(t time.Time) (int, bool) {
	return 0, false
}
