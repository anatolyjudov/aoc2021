package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

var xG0, xG1, yG0, yG1 int


func main() {

    var xWin, yWin, yMaxWin, totalHitVelocities int

    readInput()

    for xv := 1; xv <= xG1; xv++ {
        for yv := yG1; yv < 250; yv++ {
            hit, yMax := testTrajectory(xv, yv)
            if hit {
                totalHitVelocities++
            }
            if hit && yMax > yMaxWin {
                yMaxWin = yMax
                xWin = xv
                yWin = yv
            }
        }
    }


    fmt.Printf("Task 1: shot (%v, %v) reaches %v\r\n", xWin, yWin, yMaxWin)
    fmt.Printf("Total velocities found: %v", totalHitVelocities)
}

func testTrajectory(xv, yv int) (hit bool, yMax int) {
    var x, y int
    hit = false
    for ;; {
        x += xv
        y += yv
        if y > yMax {
            yMax = y
        }
        if xv > 0 {
            xv--
        }
        yv--
        if (x >= xG0) && (x <= xG1) && (y <= yG0) && (y >= yG1) {
            return true, yMax
        }
        if x > xG1 || y < yG1 {
            break
        }
    }
    return false, yMax
}


func isNatural(f float64) bool {
    return f == float64(int64(f))
}


// read all input data
func readInput() {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    scanner.Scan()
    input := scanner.Text()
    values := []byte(input)[13:]
    xys := strings.Split(string(values), ", ")
    xInfo := strings.Split(xys[0], "=")
    xBorders := strings.Split(xInfo[1], "..")
    xG0 = convInt(xBorders[0])
    xG1 = convInt(xBorders[1])
    yInfo := strings.Split(xys[1], "=")
    yBorders := strings.Split(yInfo[1], "..")
    yG0 = convInt(yBorders[1])
    yG1 = convInt(yBorders[0])

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