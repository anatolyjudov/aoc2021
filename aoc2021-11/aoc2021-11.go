package main

import (
    "fmt"
    "os"
    "bufio"
    "github.com/TwiN/go-color"
)

// Task 1
var flashCounter int

// Task 2
var stepFlashesCounter int

type Octopuses [10][10]int

    func (o *Octopuses) print() {
        for y := 0; y < 10; y++ {
            for x := 0; x < 10; x++ {
                str := color.Ize(color.White, "%d")
                if o[y][x] > 0 {
                    str = color.Ize(color.Gray, "%d")
                }
                fmt.Printf(str, o[y][x])
            }
            fmt.Print("\r\n")
        }
    }

    func (o *Octopuses) liveStep() {
        stepFlashesCounter = 0
        // grow
        for y := 0; y < 10; y++ {
            for x := 0; x < 10; x++ {
                o[y][x]++
            }
        }
        // flash
        for y := 0; y < 10; y++ {
            for x := 0; x < 10; x++ {
                if o[y][x] > 9 {
                    o.flash(x, y)
                }
            }
        }
    }

    func (o *Octopuses) flash(x, y int) {
        o[y][x] = 0
        for _, p := range getNearPoints(x, y) {
            o.affect(p[0], p[1])
        }
        flashCounter++
        stepFlashesCounter++
    }    

    func (o *Octopuses) affect(x, y int) {
        if !isValid(x, y) {
            return
        }
        if o[y][x] == 0 {
            return
        } else if o[y][x] < 9 {
            o[y][x]++
        } else {
            o.flash(x, y)
        }
    }

func isValid(x, y int) bool {
    return x >= 0 && x < 10 && y >= 0 && y < 10
}

func getNearPoints(x, y int) [8][2]int {
    p := [8][2]int{
            {x - 1, y - 1},
            {x - 1, y},
            {x, y - 1},
            {x + 1, y - 1},
            {x - 1, y + 1},
            {x, y + 1},
            {x + 1, y},
            {x + 1, y + 1},
        }
    return p
}

func readOctopuses(octopuses *Octopuses) {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    y := 0;
    for scanner.Scan() {
        for x, b := range scanner.Bytes() {
            octopuses[y][x] = int(b - 48)
        }
        y++
    }

}

//
func main() {

    var octopuses Octopuses
    readOctopuses(&octopuses)

    var step int

    for step = 1; step <= 1000; step++ {
        fmt.Printf("Step %d\r\n", step)
        octopuses.liveStep()
        octopuses.print()
        fmt.Printf("\r\nStep flashes: %d", stepFlashesCounter)
        if stepFlashesCounter == 100 {
            fmt.Printf("\r\nIt's a final step: %d", step)
            os.Exit(0)
        }
    }

    fmt.Printf("\r\nTotal flashes: %d", flashCounter)    
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}