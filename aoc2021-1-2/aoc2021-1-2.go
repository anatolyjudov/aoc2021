package main

import (
	"fmt"
	"os"
	"strconv"
	"bufio"
)

type Frame struct {
	d1, d2, d3, d4 int
}

func (frame *Frame) add(measure int) {
	frame.d1 = frame.d2
	frame.d2 = frame.d3
	frame.d3 = frame.d4
	frame.d4 = measure
}

func (frame *Frame) checkIncreased() (bool, bool) {
	if frame.d1 == 0 {
		return false, true
	}
	if frame.d4 > frame.d1 {
		return true, true
	}
	return false, true
}

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

    increasesCount := 0
    frame := Frame{0,0,0,0}

    for scanner.Scan() {
    	depth, err := strconv.Atoi(scanner.Text())
    	if err != nil {
    		panic(err)
    	}
    	frame.add(depth);

    	isIncreased, applicable := frame.checkIncreased()
    	if applicable && isIncreased {
			increasesCount += 1;
    	}
    }

    fmt.Println(increasesCount)
}