package notebook

import (
	"testing"
	"time"
)

func TestNoteInfoComparison(t *testing.T) {
	a := NoteInfo{}
	b := NoteInfo{"", time.Now()}
	if a != b {
		t.Fatalf("a not equal to b")
	}
}
