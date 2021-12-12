package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    "strconv"
)

type Board struct {
    numbers [25]int
    marks   [25]int
    hasWon  bool
}

    func (b *Board) loadNumbers(numbers [5]int, line int) {
        start := line * 5
        for i := 0; i < 5; i++ {
            b.numbers[start + i] = numbers[i]
        }
    }

    func (b *Board) checkLine(line int) bool {
        start := line * 5
        for n := 0; n < 5; n++ {
            if b.marks[start + n] == 0 {
                return false
            }
        }
        return true
    }
    func (b *Board) checkColumn(col int) bool {
        for l := 0; l < 5; l++ {
            if b.marks[l * 5 + col] == 0 {
                return false
            }
        }
        return true
    }

    func (b *Board) checkWin() bool {
        for i := 0; i < 5; i++ {
            if b.checkLine(i) || b.checkColumn(i) {
                b.hasWon = true
                return true
            }
        }
        b.hasWon = false
        return false
    }

    func (b *Board) mark(drawn int) {
        for i := 0; i < 25; i++ {
            if b.numbers[i] == drawn {
                b.marks[i] = 1
            }
        }
    }

    func (b *Board) unmarkedSum() int {
        var sum int
        for i := 0; i < 25; i++ {
            if b.marks[i] == 0 {
                sum += b.numbers[i]
            }
        }
        return sum
    }

//
func main() {
    
    drawns, boards := readInputs()

    //fmt.Println(drawns)
    //fmt.Println(boards)

    for d, drawn := range drawns {
        fmt.Printf("draw %d, %d drawned\r\n", d, drawn)
        for b, _ := range boards {
            if (boards[b].hasWon) {
                continue
            }
            boards[b].mark(drawn)
            if boards[b].checkWin() {
                fmt.Printf("Board %d has won on the drawn %d\r\n", b, d)
                //fmt.Println(boards[b])
                sum := boards[b].unmarkedSum()
                fmt.Printf("Its sum is %d\r\n", sum)
                res := sum * drawn
                fmt.Printf("And result is %d\r\n", res)
                //os.Exit(0)
            }
        }
    }

}

// read all data
func readInputs() (drawns []int, boards []Board) {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    // first string is inputs
    scanner.Scan()
    drawnsLine := scanner.Text()
    drawnsScanned := strings.Split(drawnsLine, ",")
    for d := 0; d < len(drawnsScanned); d++ {
        value, err := strconv.Atoi(drawnsScanned[d])
        if err != nil {
            panic(err)
        }
        drawns = append(drawns, value)
    }

    // read boards
    var (
        currentLine int
        currentBoard Board
    )
    for scanner.Scan() {
        currentString := scanner.Text()
        if currentString == "" {
            currentLine  = 0
            currentBoard = Board{}
        } else {
            currentBoard.loadNumbers(parseNumbers(currentString), currentLine)
            if currentLine == 4 {
                boards = append(boards, currentBoard)
            }
            currentLine++
        }
    }

    return drawns, boards
}

func parseNumbers(numbersString string) [5]int {
    var res [5]int

    parts := strings.Split(numbersString, " ")

    currentNumber := 0
    for i := 0; i < len(parts); i++ {
        if parts[i] != "" {
            value, err := strconv.Atoi(parts[i])
            if err != nil {
                fmt.Println("Bad numbers given: %d", numbersString)
                panic("Fatal error when reading input file")
            }
            res[currentNumber] = value
            currentNumber++;
            if (currentNumber >= 5) {
                break
            }
        }
    }
    return res
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}