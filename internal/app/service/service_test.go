package service

import "testing"

func TestService_RetShort(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			got, err := s.RetShort(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetShort() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RetShort() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_RetLong(t *testing.T) {
	type fields struct {
		storage Storage
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}
			got, err := s.RetLong(tt.args.shortURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("RetLong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RetLong() got = %v, want %v", got, tt.want)
			}
		})
	}
}
