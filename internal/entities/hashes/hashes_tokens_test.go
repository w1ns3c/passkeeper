package hashes

import "testing"

func TestCheckToken(t *testing.T) {
	type args struct {
		token  string
		secret string
	}

	t1 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ5MDYwODcsIlVzZXJJRCI6InVzZXIxIn0.n-n9x8swkhCByPwvn9e6mtoH4sKx53B_qOpODjMmZl0"
	id1 := "user1"
	secret1 := "3ec6b6e734d97a3392e08052c50a167f"

	tests := []struct {
		name       string
		args       args
		wantUserID string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "Test 1",
			args: args{
				token:  t1,
				secret: secret1,
			},
			wantUserID: id1,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := CheckToken(tt.args.token, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("CheckToken() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestExtractUserID(t *testing.T) {
	t1 := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjQ5MDYwODcsIlVzZXJJRCI6InVzZXIxIn0.n-n9x8swkhCByPwvn9e6mtoH4sKx53B_qOpODjMmZl0"
	id1 := "user1"

	tests := []struct {
		name       string
		token      string
		wantUserID string
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name:       "Test 1",
			token:      t1,
			wantUserID: id1,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ExtractUserID(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ExtractUserID() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
