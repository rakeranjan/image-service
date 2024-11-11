package apiutils

import "testing"

func TestValidImageFile(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "With valid fileName",
			args: args{
				fileName: "a.jpg",
			},
			want: true,
		},
		{
			name: "With invalid fileName",
			args: args{
				fileName: "a.pdf",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidImageFile(tt.args.fileName); got != tt.want {
				t.Errorf("ValidImageFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
