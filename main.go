package main

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"
)

type CSVData struct {
	dsName    string
	header    []string
	records   [][]string
	isNumeric []bool
}

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	files := os.Args[1:]
	for _, f := range files {
		var fname string = filenameWithoutExtension(f)
		fmt.Printf("'%s'\n", validateMemName(fname))
	}

	fmt.Println("Initial commit line.")
}

func usage() {
	fmt.Printf("Usage: %s file1 [file2...]\n", path.Base(os.Args[0]))
}

func filenameWithoutExtension(filepath string) string {
	var fileName string = path.Base(filepath)
	var fileExtension string = path.Ext(filepath)
	return strings.TrimSuffix(fileName, fileExtension)
}

// validateMemName returns a string that is
// considered a SAS validmemname=compatible string.
// This includes the following rules: 1) the length is
// up to 32 characters, 2) the name may not
// contain blanks or any special characters other than the
// underscore, 3) names must begin with a Latin
// alphabet character or an underscore.
func validateMemName(sourceString string) string {
	var compatibleName string = sourceString

	// ensure that the membername is not more than 32 characters long
	if len(sourceString) > 32 {
		compatibleName = sourceString[:32]
	}

	var re = regexp.MustCompile(`(\W+)`)
	compatibleName = re.ReplaceAllString(compatibleName, "_")

	// must start with latin character or underscore
	re = regexp.MustCompile(`(?m)^([^a-zA-Z_])`)
	var substitution = "_$1"
	compatibleName = re.ReplaceAllString(compatibleName, substitution)

	return compatibleName
}

// collectColumnAsString Iterates over a particular column of the CSV
// record and collects everything into a single long string. This is
// a helper function to assist in determining if a particular column
// is numeric or not using isStringOnlyNumeric.
func collectColumnAsString(records [][]string, colNumber int) string {
	var columnString string

	// skip through the header row while accessing elements
	for _, elem := range records[1:] {
		columnString += elem[colNumber]
	}

	return columnString
}

// isStringOnlyNumeric parses the output of collectColumnAsString
// and returns true if only valid numeric symbols are found,
// false otherwise.
func isStringOnlyNumeric(input string) bool {
	var re = regexp.MustCompile(`[^\d\-\+\.]`)
	// a match means that the string contains an unexpected symbol
	// so we need to negate the bool value in the return.
	return !re.MatchString(input)
}

// Given the name of the CSV file and the [][]string returned by the
// CSV reader, initialize a CSVData element.
func initializeCSVData(filename string, csvrecords [][]string) CSVData {
	var data CSVData
	data.dsName = validateMemName(filename)
	data.records = csvrecords[1:]

	data.header = make([]string, len(csvrecords[0]))
	data.isNumeric = make([]bool, len(data.header))

	for i := range data.header {
		data.header[i] = validateMemName(csvrecords[0][i])

		columnValues := collectColumnAsString(data.records, i)
		data.isNumeric[i] = isStringOnlyNumeric(columnValues)
	}

	return data
}
