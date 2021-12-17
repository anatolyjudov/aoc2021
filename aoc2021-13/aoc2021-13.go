package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

type Dot struct {
    x, y int
}
    func (d *Dot) putPair(xy []string) {
        d.x = convInt(xy[0])
        d.y = convInt(xy[1])
    }
    func (d *Dot) equal(x, y int) bool {
        return x == d.x && y == d.y
    }

type Command struct {
    x, y int
}
    func (c *Command) load(data []string) {
        if data[0] == "x" {
            c.x = convInt(data[1])
        } else if data[0] == "y" {
            c.y = convInt(data[1])
        } else {
            fmt.Println(data)
            panic("Unknown command")
        }
    }

func fold(dots []Dot, command Command) []Dot {
    var newDots []Dot

    dotsLoop:
    for _, dot := range dots {
        newDot := reflect(dot, command)
        for _, d := range newDots {
            if d == newDot {
                continue dotsLoop
            }
        }
        newDots = append(newDots, newDot)
    }
    return newDots
}

func reflect(dot Dot, command Command) Dot {
    switch {
    case command.x == 0:
        if command.y < dot.y {
            dot.y = command.y << 1 - dot.y
        }
    case command.y == 0:
        if command.x < dot.x {
            dot.x = command.x << 1 - dot.x
        }
    }
    return dot
}

func visualize(dots []Dot) {
    var (
        maxX, maxY int
    )
    for _, dot := range dots {
        if maxX < dot.x {
            maxX = dot.x
        }
        if maxY < dot.y {
            maxY = dot.y
        }
    }
    for y := 0; y <= maxY; y++ {
        rowLoop:
        for x := 0; x <= maxX; x++ {
            for _, d := range dots {
                if d.equal(x, y) {
                    fmt.Print("#")
                    continue rowLoop
                }
            }
            fmt.Print(".")
        }
        fmt.Print("\r\n")
    }
    fmt.Print("\r\n")
}


func main() {

    var (
        dots     []Dot
        commands []Command
    )

    readInput(&dots, &commands)

    for f := 0; f < len(commands); f++ {
        dots = fold(dots, commands[f])
        if f == 0 {
            fmt.Printf("After first fold there are %v dots\r\n", len(dots))
        }
    }

    visualize(dots)
}

// read all input data
func readInput(dots *[]Dot, commands *[]Command) {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    // dots
    for scanner.Scan() {
        row := scanner.Text()
        if row == "" {
            break;
        }
        d := Dot{}
        d.putPair(strings.Split(row, ","))
        *dots = append(*dots, d)
    }

    // commands
    for scanner.Scan() {
        row := scanner.Text()
        if row == "" {
            break;
        }
        c := Command{}
        c.load(strings.Split(strings.Split(row, " ")[2], "="))
        *commands = append(*commands, c)
    }
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
