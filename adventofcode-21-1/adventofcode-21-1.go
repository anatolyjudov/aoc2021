package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
)

func openFile(filename string) *os.File {
	file, err := os.Open(filename)
    if err != nil {
    	panic(err)
    }
    return file
}

func main() {

	file := openFile("./input.txt")

    scanner := bufio.NewScanner(file)

    increasesCount := -1;
    previousDepth := -1;

    for scanner.Scan() {
    	depth, err := strconv.Atoi(scanner.Text())
    	if err != nil {
    		panic(err)
    	}
    	if (depth > previousDepth) {
    		increasesCount += 1;
    	}
    	previousDepth = depth;
    }

    fmt.Println(increasesCount)
}