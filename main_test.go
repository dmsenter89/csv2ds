package main

import "testing"

func Test_filenameWithoutExtension(t *testing.T) {
	type args struct {
		filepath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"No extension", args{"a"}, "a"},
		{"Basic extension", args{"a.csv"}, "a"},
		{"Path and extenson", args{"path/to/somewhere/a.csv"}, "a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filenameWithoutExtension(tt.args.filepath); got != tt.want {
				t.Errorf("filenameWithoutExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}
