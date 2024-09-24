package tui

import (
	"testing"

	"passkeeper/internal/entities/config"
	"passkeeper/internal/entities/structs"
)

func TestGenNoteShortName(t *testing.T) {
	var suffix = " ..."
	tests := []struct {
		name     string
		noteName string
		want     string
	}{
		{
			name:     "Test 1: short note",
			noteName: "simple note",
			want:     "simple note",
		},
		{
			name:     "Test 2: long note",
			noteName: "simple super 11111 note",
			want:     "simple super 11111 note"[:config.MaxNameLen-len(suffix)] + suffix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenNoteShortName(&structs.Note{Name: tt.noteName}); got != tt.want {
				t.Errorf("GenNoteShortName() = %v, want %v", got, tt.want)
			}
		})
	}
}
