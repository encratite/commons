package commons

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"slices"
)

func ReadCSVFile(reader io.Reader, skipHeader bool, callback func ([]string)) {
	csvReader := csv.NewReader(reader)
	if skipHeader {
		_, _ = csvReader.Read()
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		callback(record)
	}
}

func ReadCSV(path string, skipHeader bool, callback func ([]string)) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to read CSV (%s): %v", path, err)
	}
	defer file.Close()
	ReadCSVFile(file, skipHeader, callback)
	return nil
}

func ReadCSVColumnsFile(reader io.Reader, path string, columns []string, callback func ([]string)) error {
	csvReader := csv.NewReader(reader)
	csvColumns, err := csvReader.Read()
	if err == io.EOF {
		return fmt.Errorf("Failed to read CSV columns from %s: %v", path, err)
	}
	indices := []int{}
	for _, column := range columns {
		index := slices.Index(csvColumns, column)
		if index == -1 {
			return fmt.Errorf("Unable to find column \"%s\" in CSV file %s", column, path)
		}
		indices = append(indices, index)
	}
	line := 2
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		indexRecord := []string{}
		for _, index := range indices {
			if index >= len(record) {
				return fmt.Errorf("Not enough records on line %d in CSV file %s", line, path)
			}
			indexRecord = append(indexRecord, record[index])
		}
		callback(indexRecord)
		line++
	}
	return nil
}

func ReadCSVColumns(path string, columns []string, callback func ([]string)) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("Failed to read CSV (%s): %v", path, err)
	}
	defer file.Close()
	ReadCSVColumnsFile(file, path, columns, callback)
	return nil
}

func WriteCSV(path string, header []string, rows [][]string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	writer.Comma = ','
	err = writer.Write(header)
	if err != nil {
		return err
	}
	for _, row := range rows {
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	writer.Flush()
	err = writer.Error()
	if err != nil {
		return err
	}
	return nil
}