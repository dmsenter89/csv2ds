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
	simpleCSV := [][]string{{"Name", "Sex", "Age", "Height", "Weight"},
		{"Alfred", "M", "14", "69", "112.5"},
		{"Alice", "F", "13", "56.5", "84"},
		{"Barbara", "F", "13", "65.3", "-98"}}
	simpleRecords := [][]string{{"Alfred", "M", "14", "69", "112.5"},
		{"Alice", "F", "13", "56.5", "84"},
		{"Barbara", "F", "13", "65.3", "-98"}}
	simpleHeader := []string{"Name", "Sex", "Age", "Height", "Weight"}
	simpleNumeric := []bool{false, false, true, true, true}

	harderCSV := [][]string{{"Name", "4Sex", "Age!", "Height", "$Weight"},
		{"Alfred", "M", "14", "69", "112.5"},
		{"Alice", "F", "13", "56.5", "84"},
		{"Barbara", "F", "13", "65.3", "-98"}}
	harderCSVheader := []string{"Name", "_4Sex", "Age_", "Height", "_Weight"}

	type args struct {
		filename   string
		csvrecords [][]string
	}
	tests := []struct {
		name string
		args args
		want CSVData
	}{
		{"Simple CSV", args{"sampleData", simpleCSV},
			CSVData{dsName: "sampleData", header: simpleHeader, records: simpleRecords, isNumeric: simpleNumeric, maxLength: []int{7, 1, 2, 4, 5}}},
		{"Basic CSV, name to be fixed", args{"sample data", simpleCSV},
			CSVData{dsName: "sample_data", header: simpleHeader, records: simpleRecords, isNumeric: simpleNumeric, maxLength: []int{7, 1, 2, 4, 5}}},
		{"Harder CSV", args{"!Bad$Name", harderCSV},
			CSVData{dsName: "_Bad_Name", header: harderCSVheader, records: simpleRecords, isNumeric: simpleNumeric, maxLength: []int{7, 1, 2, 4, 5}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initializeCSVData(tt.args.filename, tt.args.csvrecords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initializeCSVData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeDataStepFromCSVData(t *testing.T) {

	var solution string = `data sample;
	infile datalines DSD;
	input Name $ Sex $ Age Height Weight;
	datalines;
Alfred,M,14,69,112.5
Alice,F,13,56.5,84
Barbara,F,13,65.3,-98
;
`

	type args struct {
		data CSVData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Base Case", args{
			CSVData{dsName: "sample",
				header: []string{"Name", "Sex", "Age", "Height", "Weight"},
				records: [][]string{{"Alfred", "M", "14", "69", "112.5"},
					{"Alice", "F", "13", "56.5", "84"},
					{"Barbara", "F", "13", "65.3", "-98"}},
				isNumeric: []bool{false, false, true, true, true}},
		}, solution},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := writeDataStepFromCSVData(tt.args.data); got != tt.want {
				t.Errorf("writeDataStepFromCSVData() = \n`%v`,\n want \n`%v`", got, tt.want)
			}
		})
	}
}

func Test_buildInputStatement(t *testing.T) {
	type args struct {
		header    []string
		isNumeric []bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Base case", args{[]string{"Name", "Sex", "Age", "Height", "Weight"}, []bool{false, false, true, true, true}},
			"input Name $ Sex $ Age Height Weight;"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildInputStatement(tt.args.header, tt.args.isNumeric); got != tt.want {
				t.Errorf("buildInputStatement() = `%v`, want `%v`", got, tt.want)
			}
		})
	}
}

func Test_buildDatalines(t *testing.T) {
	type args struct {
		records [][]string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Base case", args{[][]string{{"Alfred", "M", "14", "69", "112.5"},
			{"Alice", "F", "13", "56.5", "84"},
			{"Barbara", "F", "13", "65.3", "-98"}}},
			`Alfred,M,14,69,112.5
Alice,F,13,56.5,84
Barbara,F,13,65.3,-98
`}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildDatalines(tt.args.records); got != tt.want {
				t.Errorf("buildDatalines() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_maxLengthOfColumn(t *testing.T) {
	type args struct {
		records [][]string
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"Simple Base Case", args{records: [][]string{{"a", "bb", "ccc"}, {"ddd", "e", "ffff"}, {"ggggggg", "hhhh", "i"}, {"jj", "kk", "lll"}}}, []int{7, 4, 4}},
		{"SimpleRecords Case", args{records: [][]string{{"Alfred", "M", "14", "69", "112.5"},
			{"Alice", "F", "13", "56.5", "84"},
			{"Barbara", "F", "13", "65.3", "-98"}}}, []int{7, 1, 2, 4, 5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxLengthOfColumn(tt.args.records); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("maxLengthOfColumn() = %v, want %v", got, tt.want)
			}
		})
	}
}
