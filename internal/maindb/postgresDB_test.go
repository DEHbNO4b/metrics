package maindb

import (
	"testing"
)

func TestNewPostgresDB(t *testing.T) {
	got := NewPostgresDB("")
	if got != nil {
		t.Error("not nil db with nil string")
	}
}
