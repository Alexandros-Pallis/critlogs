package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Row struct {
	TargetProcessingTime string
	Request              string
}

func parseToRow(line string) Row {
	info := strings.Fields(line)
	row := Row{
		TargetProcessingTime: info[6],
		Request:              strings.Join([]string{info[12], info[13]}, " "),
	}
	return row
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func filterCritical(rows []Row, criticalLimit string) []Row {
	var criticalRows []Row
	for _, row := range rows {
		TargetProcessingTime, err := strconv.ParseFloat(row.TargetProcessingTime, 64)
		criticalLimit, err := strconv.ParseFloat(criticalLimit, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		if TargetProcessingTime > criticalLimit || row.TargetProcessingTime == "-1" {
			criticalRows = append(criticalRows, row)
		}
	}
	return criticalRows
}

func main() {
	var inputLogFile string
	var outputLogFile string
	var criticalLimit string
	var rows []Row
	flag.StringVar(&inputLogFile, "i", "", "Input log file")
	flag.StringVar(&outputLogFile, "o", "", "Output log file")
	flag.StringVar(&criticalLimit, "climit", "0.5", "Lower critical limit")
	flag.Parse()
	lines, err := readLines(inputLogFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for _, line := range lines {
		row := parseToRow(line)
		if err != nil {
			log.Fatal(err)
			continue
		}
		rows = append(rows, row)
	}

	rows = filterCritical(rows, criticalLimit)

	output, err := os.Create(outputLogFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	for _, row := range rows {
		output.WriteString(fmt.Sprintf("processing time: %s,\nrequest: %s\n\n", row.TargetProcessingTime, row.Request))
	}
	defer output.Close()
	os.Exit(0)
	return
}
