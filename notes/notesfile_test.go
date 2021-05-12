package notes

import (
	"os"
	"path"
	"reflect"
	"time"

	// "reflect"
	"testing"
)

func TestReadOrCreate_CreatesNewFileIfExists(t *testing.T) {
	f := path.Join(t.TempDir(), "temp-file.txt")
	defer os.Remove(f)

	ReadOrCreate(f)

	if _, err := os.Stat(f); os.IsNotExist(err) {
		t.Error("New notes file was not created")
	}
}

func TestReadOrCreate_ReadsFile(t *testing.T) {
	want := []string{"line1", "line2"}
	p := tempFileWithContent(t, want)

	got, _ := ReadOrCreate(p)
	if !reflect.DeepEqual(got.Lines, want) {
		t.Errorf("Wanted %v, got %v", want, got.Lines)
	}
}

func TestPositionOn_GetsPosition(t *testing.T) {
	lines := []string{
		"### Thu 2021/05/13",
		"Entry 3",
		"",
		"### Wed 2021/05/12",
		"Entry 2",
		"",
		"### Tue 2021/05/11",
		"Entry 1",
		"",
	}

	n := NotesFile{
		Lines: lines,
	}

	d := time.Date(2021, 5, 12, 0, 0, 0, 0, time.UTC)
	want := 4
	got, _ := n.positionOn(d)
	if got != want {
		t.Errorf("Wanted %d got %d", want, got)
	}
}

func tempFileWithContent(t *testing.T, content []string) string {
	p := path.Join(t.TempDir(), "temp-file.txt")
	f, err := os.Create(p)
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}

	for _, line := range content {
		_, _ = f.WriteString(line + "\n")
	}

	return p
}
