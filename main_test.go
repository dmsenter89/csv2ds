package main

import (
	"reflect"
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

func Test_collectColumnAsString(t *testing.T) {
	sampleCSV := [][]string{{"Name", "Sex", "Age", "Height", "Weight"},
		{"Alfred", "M", "14", "69", "112.5"},
		{"Alice", "F", "13", "56.5", "84"},
		{"Barbara", "F", "13", "65.3", "-98"}}

	type args struct {
		records   [][]string
		colNumber int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Letters only", args{sampleCSV, 1}, "MFF"},
		{"Ints Only", args{sampleCSV, 2}, "141313"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collectColumnAsString(tt.args.records, tt.args.colNumber); got != tt.want {
				t.Errorf("collectColumnAsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isStringOnlyNumeric(t *testing.T) {

	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Character string only", args{"AlfredAliceBarbara"}, false},
		{"Ints only", args{"141313"}, true},
		{"Ints and dollar sign", args{"1413$13"}, false},
		{"Floats with signs", args{"69-56.5+65.3"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isStringOnlyNumeric(tt.args.input); got != tt.want {
				t.Errorf("isStringOnlyNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initializeCSVData(t *testing.T) {
	sampleCSV := [][]string{{"Name", "Sex", "Age", "Height", "Weight"},
		{"Alfred", "M", "14", "69", "112.5"},
		{"Alice", "F", "13", "56.5", "84"},
		{"Barbara", "F", "13", "65.3", "-98"}}
	sampleRecords := [][]string{{"Alfred", "M", "14", "69", "112.5"},
		{"Alice", "F", "13", "56.5", "84"},
		{"Barbara", "F", "13", "65.3", "-98"}}
	sampleHeader := []string{"Name", "Sex", "Age", "Height", "Weight"}
	sampleNumeric := []bool{false, false, true, true, true}

	type args struct {
		filename   string
		csvrecords [][]string
	}
	tests := []struct {
		name string
		args args
		want CSVData
	}{
		{"Basic CSV, name to be fixed", args{"sample data", sampleCSV},
			CSVData{"sample_data", sampleHeader, sampleRecords, sampleNumeric}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initializeCSVData(tt.args.filename, tt.args.csvrecords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initializeCSVData() = %v, want %v", got, tt.want)
			}
		})
	}
}
