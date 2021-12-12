package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

type Field [1000][1000]int8

    func (f *Field) apply(l Line, ortho bool) {
        steps, xStep, yStep := l.getStepping(ortho)
        for s := 0; s < steps; s++ {
            x := l.x1 + s * xStep
            y := l.y1 + s * yStep
            f[x][y] += 1
        }        
    }

    func (f *Field) countPoints() int {
        var count int
        for x := 0; x < 1000; x++ {
            for y := 0; y < 1000; y++ {
                if f[x][y] > 1 {
                    count++
                    //fmt.Printf("%d %d has %d intersections\r\n", x, y, f[x][y])
                }
            }
        }
        return count
    }

type Line struct {
    x1, y1, x2, y2 int
}

    func (l *Line) getStepping(ortho bool) (int, int, int) {
        var (
            steps, xStep, yStep int
        )
        if l.x1 == l.x2 {
            // vertical
            steps = mathMod(l.y2 - l.y1) + 1
            xStep = 0
            if (l.y2 > l.y1) {
               yStep = 1
            } else {
               yStep = -1
            }
        } else if l.y1 == l.y2 {
            // horizontal
            steps = mathMod(l.x2 - l.x1) + 1
            if (l.x2 > l.x1) {
               xStep = 1
            } else {
               xStep = -1
            }
            yStep = 0
        } else if mathMod(l.y2 - l.y1) == mathMod(l.x2 - l.x1) {
            // diagonal
            if !ortho {
               steps = mathMod(l.x2 - l.x1) + 1
            } else {
               steps = 0
            }
            if (l.x2 > l.x1) {
               xStep = 1
            } else {
               xStep = -1
            }
            if (l.y2 > l.y1) {
               yStep = 1
            } else {
               yStep = -1
            }
        } else {
            panic("Incorrect line given")
        }
        return steps, xStep, yStep
    }

func readLines() []Line {
    var res []Line

    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        res = append(res, parseLine(scanner.Text()))
    }

    return res
}

func parseLine(str string) Line {
    var l Line
    pointsStr := strings.Split(str, " -> ")
    p := strings.Split(pointsStr[0], ",")
    l.x1 = convInt(p[0])
    l.y1 = convInt(p[1])
    p = strings.Split(pointsStr[1], ",")
    l.x2 = convInt(p[0])
    l.y2 = convInt(p[1])
    return l
}

//
func main() {

    var field Field
    lines := readLines()

    for i := 0; i < len(lines); i++ {
        field.apply(lines[i], false)
    }

    task1 := field.countPoints()
    fmt.Printf("Result is %d.\r\n", task1)
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

func mathMod(i int) int {
    if i < 0 {
        return -1 * i
    }
    return i
}