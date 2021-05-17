package notes

import (
	"bytes"
	"os"
	"path"
	"reflect"

	"testing"
)

func TestReadFromFile(t *testing.T) {
	tests := []struct {
		path string
		want []string
	}{
		{tmpFileWithContent(t, "line1", "line2"), []string{"line1", "line2"}},
		{"does-not-exist.txt", nil},
		{"", nil},
	}

	for _, tt := range tests {
		n := ReadFromFile(tt.path)
		if !reflect.DeepEqual(n.Lines, tt.want) {
			t.Errorf("Expected %q, got %q", tt.want, n.Lines)
		}
	}
}

func TestBackup(t *testing.T) {
	tests := []NotesFile{
		ReadFromFile(tmpFileWithContent(t, "line1", "line2")),
		ReadFromFile("does-not-exist"),
	}

	for _, tt := range tests {
		bkp, err := tt.Backup()
		if err != nil {
			t.Fatal(err)
		}
		compareFiles(t, tt.Path, bkp)
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		n       NotesFile
		changes []string
	}{
		{ReadFromFile(tmpFileWithContent(t, "line1", "line2")), []string{"line1", "line2", "line3"}},
		{NotesFile{Path: path.Join(t.TempDir(), "does-not-exist")}, []string{"line1"}},
	}

	for _, tt := range tests {
		tt.n.Lines = tt.changes
		if err := tt.n.Save(); err != nil {
			t.Error(err)
		}

		saved := ReadFromFile(tt.n.Path)
		if !reflect.DeepEqual(saved.Lines, tt.changes) {
			t.Errorf("Saved lines: %v, want %v", saved.Lines, tt.changes)
		}
	}
}

func TestSaveWithBackup(t *testing.T) {
	tests := []struct {
		n       NotesFile
		changes []string
	}{
		{ReadFromFile(tmpFileWithContent(t, "line1", "line2")), []string{"line1", "line2", "line3"}},
		{NotesFile{Path: path.Join(t.TempDir(), "does-not-exist")}, []string{"line1"}},
	}

	for _, tt := range tests {
		original := tt.n.Lines
		tt.n.Lines = tt.changes
		bkp, err := tt.n.SaveWithBackup()
		if err != nil {
			t.Error(err)
		}

		saved := ReadFromFile(tt.n.Path)
		if !reflect.DeepEqual(saved.Lines, tt.changes) {
			t.Errorf("Saved lines: %v, want %v", saved.Lines, tt.changes)
		}

		if bkp != "" {
			backup := ReadFromFile(bkp)
			if !reflect.DeepEqual(backup.Lines, original) {
				t.Errorf("Backup for %q: %v, want %v", tt.n.Path, backup.Lines, original)
			}
		}

	}
}

func TestFindContaining(t *testing.T) {
	tests := []struct {
		lines []string
		term  string
		index int
		found bool
	}{
		{[]string{"line1", "line2", "line3"}, "line2", 1, true},
		{[]string{"line1", "line2", "line3"}, "ne2", 1, true},
		{[]string{"line1", "line2", "line3"}, "line5", 0, false},
	}

	for _, tt := range tests {
		n := NotesFile{Lines: tt.lines}
		if index, found := n.FindContains(tt.term); index != tt.index || found != tt.found {
			t.Errorf("n.FindContains(%q) = (%d, %t) want (%d, %t)", tt.term, index, found, tt.index, tt.found)
		}
	}
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

func tmpFileWithContent(t *testing.T, content ...string) string {
	p := path.Join(t.TempDir(), "temp-file.txt")
	f, err := os.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for _, line := range content {
		_, _ = f.WriteString(line + "\n")
	}

	return p
}

func compareFiles(t *testing.T, f1 string, f2 string) {
	if isNotExist(f1) && isNotExist(f2) {
		return
	}

	b1, err := os.ReadFile(f1)
	if err != nil {
		t.Fatalf("Failed to read %q, %v", f1, err)
	}

	b2, err := os.ReadFile(f2)
	if err != nil {
		t.Fatalf("Failed to read %q, %v", f2, err)
	}

	if !bytes.Equal(b1, b2) {
		t.Errorf("Contens of %q and %q are different\n%q\n%q", f1, f2, b1, b2)
	}
}

func isNotExist(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return true
	}
	return false
}
