package journal

import (
	"bufio"
	"os"
	"time"
)

type Journal struct {
	Path  string
	Lines []string
}

// Read all lines in file
func Read(p string) (Journal, error) {
	err := createIfIsNotExist(p)
	if err != nil {
		return Journal{}, err
	}

	file, err := os.Open(p)
	defer file.Close()

	if err != nil {
		return Journal{}, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return Journal{Path: p, Lines: lines}, nil
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
func (j *Journal) FindPosition(t time.Time) int {

	// Looks for the latest entry before the current day entry, and get the position two lines above it
	position, found := j.positionBefore(t)
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
	position, _ = j.positionOn(t)
	return position
}

func (j *Journal) positionOn(t time.Time) (int, bool) {
	return 3, true
}

func (j *Journal) positionBefore(t time.Time) (int, bool) {
	return 0, false
}
