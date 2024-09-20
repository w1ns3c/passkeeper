package tui

import "testing"

func TestBeautifyCard(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   string
	}{
		{
			name:   "Test 1: valid card",
			number: "1111222233334444",
			want:   "1111 2222 3333 4444",
		},
		{
			name:   "Test 2: short card",
			number: "11112222333344",
			want:   "1111 2222 3333 44",
		},
		{
			name:   "Test 3: shorter card",
			number: "1111222",
			want:   "1111 222",
		},
		{
			name:   "Test 4: one section card",
			number: "1111",
			want:   "1111",
		},
		{
			name:   "Test 5: less than one section card",
			number: "111",
			want:   "111",
		},
		{
			name:   "Test 5: longer than valid card",
			number: "111122223333444455",
			want:   "1111 2222 3333 4444 55",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BeautifyCard(tt.number); got != tt.want {
				t.Errorf("BeautifyCard() = %v, want %v", got, tt.want)
			}
		})
	}
}
