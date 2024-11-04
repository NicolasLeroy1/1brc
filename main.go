package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	rowIterator := NewRowIterator("measurements.txt")
	var rowData RowData
	DataProcessor := NewDataProcessor()
	for rowIterator.HasNext() {
		row := rowIterator.Next()
		rowData = NewRowData(row)
		DataProcessor.Process(rowData)
	}
	var sortedKeys []string = DataProcessor.SortKeys()
	DataProcessor.Print(sortedKeys)
}

type Row string

type RowData struct {
	StationName string
	Temperature float64
}

func NewRowData(row Row) RowData {
	var rowData RowData
	var err error
	split := strings.Split(string(row), ";")
	rowData.StationName = split[0]
	rowData.Temperature, err = strconv.ParseFloat(split[1], 64)
	if err != nil {
		panic(err)
	}
	return rowData
}

type DataProcessor struct {
	ProcessedData map[string]StationData
}

func (dp *DataProcessor) Process(rowData RowData) {
	if _, ok := dp.ProcessedData[rowData.StationName]; !ok {
		dp.ProcessedData[rowData.StationName] = StationData{
			Min: rowData.Temperature,
			Max: rowData.Temperature,
			Avg: rowData.Temperature,
			n:   1,
		}
	} else {
		stationData := dp.ProcessedData[rowData.StationName]
		stationData.Min = min(stationData.Min, rowData.Temperature)
		stationData.Max = max(stationData.Max, rowData.Temperature)
		stationData.Avg = (stationData.Avg*float64(stationData.n) + rowData.Temperature) / float64(stationData.n+1)
		stationData.n++
		dp.ProcessedData[rowData.StationName] = stationData
	}
}

func (dp *DataProcessor) SortKeys() []string {
	var keys []string
	for key := range dp.ProcessedData {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func (dp *DataProcessor) Print(sortedKeys []string) {
	var stationData StationData
	for _, key := range sortedKeys {
		stationData = dp.ProcessedData[key]
		fmt.Printf("%s=%.1f/%.1f/%.1f\n", key, stationData.Min, stationData.Max, stationData.Avg)
	}
}

type StationData struct {
	Min float64
	Max float64
	Avg float64
	n   int
}

func NewDataProcessor() *DataProcessor {
	return &DataProcessor{
		ProcessedData: make(map[string]StationData),
	}
}

type RowIterator struct {
	file    *os.File
	scanner *bufio.Scanner
}

func (ri *RowIterator) HasNext() bool {
	return ri.scanner.Scan()
}

func (ri *RowIterator) Next() Row {
	return Row(ri.scanner.Text())
}

func NewRowIterator(filepath string) *RowIterator {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	return &RowIterator{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
}
