package qdrantkit

import "testing"

func TestCreateCollection(t *testing.T) {
	type args struct {
		name                    string
		createCollectionRequest CreateCollectionRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{name: "xxx", createCollectionRequest: CreateCollectionRequest{Vectors: Vectors{
				Size:     111,
				Distance: "Cosine",
			}}},
			wantErr: false,
		},
	}
	client := New("http://localhost:6333", "test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.CreateCollection(tt.args.name, tt.args.createCollectionRequest); (err != nil) != tt.wantErr {
				t.Errorf("CreateCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetCollection(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "1",
			args:    args{name: "xxx"},
			wantErr: false,
		},
	}
	client := New("http://localhost:6333", "test")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.GetCollection(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("GetCollection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
