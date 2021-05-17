package notes

import (
	"bytes"
	"os"
	"path"
	"reflect"

	"testing"
)

func TestFindContains(t *testing.T) {
	n := NotesFile{Lines: []string{"line1", "line2", "line3"}}

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

func TestFromFile_ReadsFile(t *testing.T) {
	want := []string{"line1", "line2"}
	p := writeToTempFile(t, want)

	got := ReadFromFile(p)
	if !reflect.DeepEqual(got.Lines, want) {
		t.Errorf("Wanted %q, got %q", want, got.Lines)
	}
}

func TestSaveWithBackup(t *testing.T) {
	n := NotesFile{
		Path:  path.Join(t.TempDir(), "notes-file.txt"),
		Lines: []string{"line1", "line2"},
	}

	bkp, err := n.SaveWithBackup()
	if err != nil {
		t.Error(err)
	}

	want := "line1\nline2\n"
	got, _ := os.ReadFile(n.Path)

	if string(got) != want {
		t.Errorf("Wanted %q, Got %q", want, got)
	}

	compareFiles(t, n.Path, bkp)
}

func TestBackup(t *testing.T) {
	want := []string{"line1", "line2"}
	path := writeToTempFile(t, want)
	note := ReadFromFile(path)
	bkp, err := note.Backup()
	if err != nil {
		t.Fatal(err)
	}

	compareFiles(t, note.Path, bkp)
}

func compareFiles(t *testing.T, f1 string, f2 string) {
	c1, err := os.ReadFile(f1)
	if err != nil {
		t.Fatalf("Failed to read %s, %v", f1, err)
	}

	c2, err := os.ReadFile(f2)
	if err != nil {
		t.Fatalf("Failed to read %s, %v", f2, err)
	}

	if bytes.Compare(c1, c2) != 0 {
		t.Errorf("Contents of %s: %q is not equal to contents of %s %q", f1, c1, f2, c2)
	}
}

func writeToTempFile(t *testing.T, content []string) string {
	p := path.Join(os.TempDir(), "temp-file.txt")
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

var InsertTests = []struct {
	dst    []string
	index  int
	values []string
	want   []string
}{
	{[]string{"line1", "line2"}, 1, []string{"line3", "line4"}, []string{"line1", "line3", "line4", "line2"}},
	{[]string{"line1", "line2"}, -1, []string{"line3", "line4"}, []string{"line3", "line4", "line1", "line2"}},
	{[]string{"line1", "line2"}, 0, []string{"line3", "line4"}, []string{"line3", "line4", "line1", "line2"}},
	{[]string{"line1", "line2"}, 2, []string{"line3", "line4"}, []string{"line1", "line2", "line3", "line4"}},
	{[]string{"line1", "line2"}, 5, []string{"line3", "line4"}, []string{"line1", "line2", "line3", "line4"}},
	{[]string{"line1", "line2"}, 0, nil, []string{"line1", "line2"}},
	{[]string{}, 0, []string{"line1", "line2"}, []string{"line1", "line2"}},
}

func TestInsert(t *testing.T) {
	for _, tt := range InsertTests {
		n := NotesFile{Lines: tt.dst}
		n.Insert(tt.values, tt.index)
		if !reflect.DeepEqual(n.Lines, tt.want) {
			t.Errorf("n.Insert(%v, %d); n.Lines = %v; want %v", tt.values, tt.index, n.Lines, tt.want)
		}
	}
}

func TestInsertStrings(t *testing.T) {
	for _, tt := range InsertTests {
		if got := insertString(tt.dst, tt.index, tt.values...); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("insertString(%v, %d, %v) = %v; want %v", tt.dst, tt.index, tt.values, got, tt.want)
		}
	}
}
