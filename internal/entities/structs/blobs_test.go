package structs

import (
	"testing"
	"time"
)

func TestCard_GetID(t *testing.T) {
	type fields struct {
		Type        BlobType
		ID          string
		Name        string
		Bank        string
		Person      string
		Number      int
		CVC         int
		Expiration  time.Time
		PIN         int
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			want: "secretID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Card{
				Type:        tt.fields.Type,
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				Bank:        tt.fields.Bank,
				Person:      tt.fields.Person,
				Number:      tt.fields.Number,
				CVC:         tt.fields.CVC,
				Expiration:  tt.fields.Expiration,
				PIN:         tt.fields.PIN,
				Description: tt.fields.Description,
			}
			if got := c.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCard_SetID(t *testing.T) {
	type fields struct {
		Type        BlobType
		ID          string
		Name        string
		Bank        string
		Person      string
		Number      int
		CVC         int
		Expiration  time.Time
		PIN         int
		Description string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			args: args{
				id: "secretID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Card{
				Type:        tt.fields.Type,
				ID:          tt.fields.ID,
				Name:        tt.fields.Name,
				Bank:        tt.fields.Bank,
				Person:      tt.fields.Person,
				Number:      tt.fields.Number,
				CVC:         tt.fields.CVC,
				Expiration:  tt.fields.Expiration,
				PIN:         tt.fields.PIN,
				Description: tt.fields.Description,
			}
			c.SetID(tt.args.id)
		})
	}
}

func TestCredential_GetID(t *testing.T) {
	type fields struct {
		Type        BlobType
		ID          string
		Date        time.Time
		Resource    string
		Login       string
		Password    string
		Description string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			want: "secretID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Credential{
				Type:        tt.fields.Type,
				ID:          tt.fields.ID,
				Date:        tt.fields.Date,
				Resource:    tt.fields.Resource,
				Login:       tt.fields.Login,
				Password:    tt.fields.Password,
				Description: tt.fields.Description,
			}
			if got := c.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCredential_SetID(t *testing.T) {
	type fields struct {
		Type        BlobType
		ID          string
		Date        time.Time
		Resource    string
		Login       string
		Password    string
		Description string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			args: args{
				id: "secretID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Credential{
				Type:        tt.fields.Type,
				ID:          tt.fields.ID,
				Date:        tt.fields.Date,
				Resource:    tt.fields.Resource,
				Login:       tt.fields.Login,
				Password:    tt.fields.Password,
				Description: tt.fields.Description,
			}
			c.SetID(tt.args.id)
		})
	}
}

func TestFile_GetID(t *testing.T) {
	type fields struct {
		ID   string
		Type BlobType
		Name string
		Body []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			want: "secretID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				ID:   tt.fields.ID,
				Type: tt.fields.Type,
				Name: tt.fields.Name,
				Body: tt.fields.Body,
			}
			if got := f.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFile_SetID(t *testing.T) {
	type fields struct {
		ID   string
		Type BlobType
		Name string
		Body []byte
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			args: args{
				id: "secretID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &File{
				ID:   tt.fields.ID,
				Type: tt.fields.Type,
				Name: tt.fields.Name,
				Body: tt.fields.Body,
			}
			f.SetID(tt.args.id)
		})
	}
}

func TestNote_GetID(t *testing.T) {
	type fields struct {
		Type BlobType
		ID   string
		Name string
		Date time.Time
		Body string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			want: "secretID",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Note{
				Type: tt.fields.Type,
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
				Date: tt.fields.Date,
				Body: tt.fields.Body,
			}
			if got := n.GetID(); got != tt.want {
				t.Errorf("GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNote_SetID(t *testing.T) {
	type fields struct {
		Type BlobType
		ID   string
		Name string
		Date time.Time
		Body string
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test 1: simple test",
			fields: fields{
				ID: "secretID",
			},
			args: args{
				id: "secretID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Note{
				Type: tt.fields.Type,
				ID:   tt.fields.ID,
				Name: tt.fields.Name,
				Date: tt.fields.Date,
				Body: tt.fields.Body,
			}
			c.SetID(tt.args.id)
		})
	}
}
