package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func count(r io.Reader, countLines bool, countBytes bool) int {
	scanner := bufio.NewScanner(r)

	if !countLines {
		scanner.Split(bufio.ScanWords)
	} else if countBytes {
		scanner.Split(bufio.ScanBytes)
	}

	wc := 0

	for scanner.Scan() {
		wc++
	}

	return wc
}

func main() {
	// -l count lines instead of words
	lines := flag.Bool("l", false, "Count lines")
	// -b count bytes instead of words
	bytes := flag.Bool("b", false, "Count bytes")

	flag.Parse()

	fmt.Println(count(os.Stdin, *lines, *bytes))
}
