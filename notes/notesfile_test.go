package notes

import (
	"os"
	"path"
	"reflect"

	"testing"
)

func TestFindContains(t *testing.T) {
	n := NotesFile{
		Lines: []string{"line1", "line2", "line3"},
	}
	want := 1

	got, found := n.FindContaining("line2")
	if got != 1 || !found {
		t.Errorf("Wanted %d, Got %d", want, got)
	}

	got, _ = n.FindContaining("ne2")
	if got != 1 || !found {
		t.Errorf("Wanted %d, Got %d", want, got)
	}

	_, found = n.FindContaining("line5")
	if found {
		t.Error("Wanted false, Got true")
	}
}

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
