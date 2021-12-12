package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
    "sort"
)

var positions []int

func readPositions() {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    scanner.Scan();
    positionsData := strings.Split(scanner.Text(), ",")
    for p := 0; p < len(positionsData); p++ {
        pos := convInt(positionsData[p])
        positions = append(positions, pos)
    }
}

func median() int {
    l := len(positions)
    if l % 2 == 0 {
        return positions[l >> 1]
    } else {
        return (positions[(l + 1) >> 1] + positions[(l - 1) >> 1]) >> 1
    }
}

func avg() int {
    l := len(positions)
    sum := 0
    for p := 0; p < len(positions); p++ {
        sum += positions[p]
    }
    return sum / l
}

func countFuel(aim int) int {
    var sum int
    for p := 0; p < len(positions); p++ {
        if (aim < positions[p]) {
            sum += positions[p] - aim
        } else {
            sum += aim - positions[p]
        }
    }
    return sum
}

func countFuel2(aim int) int {
    var sum int
    for p := 0; p < len(positions); p++ {
        if (aim < positions[p]) {
            sum += calculateMove2(positions[p] - aim)
        } else {
            sum += calculateMove2(aim - positions[p])
        }
    }
    return sum
}

func calculateMove2(diff int) int {
    var sum int
    for i := diff; i > 0; i-- {
        sum += i
    }
    return sum
}

//
func main() {

    readPositions()
    fmt.Printf("Initial positions are %d.\r\n", positions)

    sort.Ints(positions)
    fmt.Printf("Sorted positions are %d.\r\n", positions)

    m := median()
    fmt.Printf("Median is %d\r\n", m)

    task1 := countFuel(m)
    fmt.Printf("Task 1 result is %d\r\n", task1)

    avg := avg()
    fmt.Printf("Average is %d\r\n", avg)

    task2 := countFuel2(avg)
    fmt.Printf("Task 2 result is %d\r\n", task2)
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