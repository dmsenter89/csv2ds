// The CSV2DS Program
//
// This program provides a command line tool for turning a CSV file
// into a SAS data step. If a list of CSV files is provided, multiple
// data steps are generated with each data set named like the basename
// of the CSV file(s). If "-" is passed as an argument, the CSV file
// is read from STDIN. This allows usage of this program with a pipe.
// All output is printed to STDOUT.
//
// This software uses minimal parsing to guess whether a particular
// entry is numeric. If an entry is determined to be a string, the
// appropriate LENGTH statement will be set. No length statements are
// set for numeric columns. Since only minimal parsing is done, you may
// need to add INFORMAT or other statements to get the desired result.
//
// Author: Michael Senter, PhD
package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

// A CSVData struct holds the key information needed
// to populate a data step template.
type CSVData struct {
	dsName    string     // a SAS valid member name to be used for the SAS data set
	header    []string   // the header row of the CSV file
	records   [][]string // the records in the CSV file
	isNumeric []bool     // whether a particular column is numeric
	maxLength []int      // maximum length of a string column
}

// The processing logic behind main is to first check if arguments are
// passed to the program. If not, the function prints a usage statement
// and exits. If arguments are given, the are processed in turn and populate
// a template string. After all files are procesed, the template string is
// printed to STDOUT.
func main() {
	var args int = len(os.Args)
	var output string

	if args > 1 { // process all files given in CLI
		for i := range os.Args[1:] {
			filename := os.Args[i+1]
			ds := processFile(filename)
			output += ds
		}
	} else {
		usage()
		os.Exit(1)
	}

	fmt.Println(output)

}

func usage() {
	fmt.Printf("Usage: %s file1 [file2...]\n", path.Base(os.Args[0]))
	fmt.Printf("\nConverts one or more CSV files to a SAS Data Step using the datalines statement.\n")
	fmt.Printf("Output is written to stdout. The data set name will be the basename of fileD\n")
	fmt.Printf("without the extension. If fileD equals '-' the CSV data is read from stdin.\n")
}

// Returns the file's base name without the extension.
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

// maxLengthOfColumn traverses each column of a CSV record to find
// the entry that consists of the longest string.
func maxLengthOfColumn(records [][]string) []int {
	maxLength := make([]int, len(records[0]))

	for _, row := range records {
		for i, entry := range row {
			var length int = len(entry)
			if length > maxLength[i] {
				maxLength[i] = length
			}
		}
	}

	return maxLength
}

// Given the name of the CSV file and the [][]string returned by the
// CSV reader, initialize a CSVData element.
func initializeCSVData(filename string, csvrecords [][]string) CSVData {
	var data CSVData
	if filename == "-" {
		data.dsName = "SAMPLEDATA"
	} else {
		data.dsName = validateMemName(filenameWithoutExtension(filename))
	}

	data.records = csvrecords[1:]

	data.header = make([]string, len(csvrecords[0]))
	data.isNumeric = make([]bool, len(data.header))
	data.maxLength = maxLengthOfColumn(data.records)

	for i := range data.header {
		data.header[i] = validateMemName(csvrecords[0][i])

		columnValues := collectColumnAsString(data.records, i)
		data.isNumeric[i] = isStringOnlyNumeric(columnValues)
	}

	return data
}

// writeDataStepFromCSVData uses the fields in CSVData to generate
// a complete template data step that can be run in SAS.
func writeDataStepFromCSVData(data CSVData) string {
	var template string = fmt.Sprintf("data %s;\n", data.dsName)

	template += fmt.Sprintln("\tinfile datalines DSD;")

	var lenstatement string = buildLengthStatement(data)

	if lenstatement != "" {
		template += fmt.Sprintf("\t%s\n", lenstatement)
	}

	template += fmt.Sprintf("\t%s\n", buildInputStatement(data.header, data.isNumeric))
	template += fmt.Sprintln("\tdatalines;")

	template += fmt.Sprint(buildDatalines(data.records))

	template += ";\n"

	return template
}

// By default, SAS stores character variables as 8 bytes. A length
// statement is used to specify that a longer string is meant to be
// stored. The buildLengthStatement function iterates over the maximum
// column item lengths and generates the appropriate length statement.
// If no length statement is neccesary, an empty string is returned.
func buildLengthStatement(data CSVData) string {
	var statement string = "length"
	for i, elem := range data.maxLength {
		if elem > 8 && !data.isNumeric[i] {
			statement += fmt.Sprintf(" %s $%d", data.header[i], elem)
		}
	}

	statement += ";"

	// if no character is longer than 8 bytes, return emtpy string
	if statement == "length;" {
		statement = ""
	}

	return statement
}

// buildInputStatement generates the input statement for the SAS data
// step without preceeding tab or newline. Adds the '$' for after the
// name of string variables.
func buildInputStatement(header []string, isNumeric []bool) string {
	var statement string = "input"

	for i, name := range header {
		statement += fmt.Sprintf(" %s", name)

		if !isNumeric[i] {
			statement += " $"
		}
	}

	statement += ";"

	return statement
}

// buildDatalines writes the CSV input for the actual records
// back to a string for use in the datalines statement.
func buildDatalines(records [][]string) string {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.WriteAll(records)
	return buf.String()
}

// processFile is the main driver of this program. It reads the relevant
// file, initializes a CSVData object and generates the data step template.
func processFile(filename string) string {
	// Contents is a byte slice holding an entire CSV record, including header.
	var contents []byte
	if filename == "-" {
		// read from STDIN
		contents = readSTDIN()
	} else {
		contents = readFile(filename)
	}

	// records will contain the entire CSV file contents
	var records [][]string = readCSV(contents)
	var data CSVData = initializeCSVData(filename, records)
	var ds string = writeDataStepFromCSVData(data)

	return ds
}

// Reads the entire contents of a file into a byte slice.
func readFile(filepath string) []byte {
	content, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error - cannot open file %s.\n", filepath)
		os.Exit(2)
	}
	return content
}

// Reads the STDIN into a byte slice. This function is
// used to read from a pipe.
func readSTDIN() []byte {
	reader := bufio.NewReader(os.Stdin)
	pipe, err := io.ReadAll(reader)
	if err != nil {
		fmt.Println("Error reading from STDIN.")
		os.Exit(2)
	}
	return pipe
}

// Wrapper around reading a CSV from a byte slice.
func readCSV(content []byte) [][]string {
	reader := csv.NewReader(bytes.NewReader(content))

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error - cannot read contents of current file.\n")
		os.Exit(3)
	}
	return records
}
