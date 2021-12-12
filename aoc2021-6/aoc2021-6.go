package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

var fishes [9]int64

func readFishes() {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    scanner.Scan();
    fishesData := strings.Split(scanner.Text(), ",")
    for f := 0; f < len(fishesData); f++ {
        age := convInt(fishesData[f])
        if age > 9 {
            panic("Bad fish age")
        }
        fishes[age]++
    }
}

func fishesLiveAnotherDay() {
    var newFishes [9]int64

    newFishes[6] = fishes[0]
    newFishes[8] = fishes[0]
    for f := 0; f < 8; f++ {
        newFishes[f] += fishes[f+1]
    }

    fishes = newFishes
}

func sumFishes() int64 {
    var sum int64
    for f := 0; f < 9; f++ {
        sum += fishes[f]
    }
    return sum
}

//
func main() {

    readFishes()

    fmt.Printf("Initial fishes are %d.\r\n", fishes)

    var task1, task2 int64

    for d := 0; d < 256; d++ {
        fishesLiveAnotherDay()
        fmt.Printf("Day %d. Fishes are %d \r\n", d, fishes)
        if d == 79 {
            task1 = sumFishes()
        }
    }

    task2 = sumFishes();

    fmt.Printf("Task 1 result is %d\r\n", task1)
    fmt.Printf("Task 2 result is %d", task2)
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}

func convInt(s string) int {
    value, err := strconv.Atoi(s)
    if err != nil {
        fmt.Println("Bad number given: %d", s)
        panic("Fatal error when reading input file")
    }
    return value
}