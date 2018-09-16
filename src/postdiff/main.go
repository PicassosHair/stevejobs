package main

import (
	"bufio"
	"flag"
	"os"
	"strings"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

// This helper is to parse the diff result file.
// Only keep changed lines.
func main() {
	inPath := flag.String("in", "", "The diff json file to parse.")
	outPath := flag.String("out", "", "The destination to generate file.")

	flag.Parse()

	inFile, err := os.Open(*inPath)
	defer inFile.Close()
	checkErr(err)

	writeFile, err := os.Create(*outPath)
	defer writeFile.Close()
	checkErr(err)

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "<") || strings.HasPrefix(line, ">") {
			writeFile.WriteString(line[2:])
			writeFile.WriteString("\n")
		}
	}
	writeFile.Sync()

	checkErr(err)
}
