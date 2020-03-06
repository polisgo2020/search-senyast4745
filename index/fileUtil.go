package main

import (
	"bufio"
	"log"
	"os"
)

func readFileByWords(fn string) ([]string, error) {

	file, err := os.Open(fn)
	if err != nil {
		log.Fatalf("error while ")
		return nil, err
	}
	//noinspection GoUnhandledErrorResult
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	var data []string
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return data, nil
}

func writeDataToFile(str string) error {

	f, err := os.Create(finalDataFile)
	if err != nil {
		return err
	}
	//noinspection GoUnhandledErrorResult
	defer f.Close()

	if _, err := f.WriteString(str); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}
