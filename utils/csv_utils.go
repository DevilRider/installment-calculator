package utils

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func WriteCsv(file string, rows [][]string) {
	csvFile, err := os.Create(file)
	if err != nil {
		logrus.Infof("failed creating file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)

	for _, row := range rows {
		_ = csvWriter.Write(row)
	}
	csvWriter.Flush()
}

func ReadCsv(file string) ([][]string, error) {
	csvFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	var rows [][]string
	reader := csv.NewReader(csvFile)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		rows = append(rows, row)
	}
	return rows, nil
}
