package main

import (
	"testing"
)

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

func Test_validateMemName(t *testing.T) {
	type args struct {
		fileBase string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Compatible name", args{"data1"}, "data1"},
		{"Long compatible name", args{"Long_Name_That_Is_Compatible_Just_Too_Long"}, "Long_Name_That_Is_Compatible_Jus"},
		{"Starts with a number", args{"45Name"}, "_45Name"},
		{"Starts with a Special Char", args{"!name"}, "_name"},
		{"Contains Spaces", args{"Name With Spaces"}, "Name_With_Spaces"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validateMemName(tt.args.fileBase); got != tt.want {
				t.Errorf("validateMemName() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
