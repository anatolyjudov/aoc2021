package main

import (
    "fmt"
    "os"
    "bufio"
    "strconv"
    "strings"
    "github.com/TwiN/go-color"
    "sort"
)

func getColor(n int) string {
    r := n % 4
    switch r {
    case 0:
        return color.Green
    case 1:
        return color.Yellow
    case 2:
        return color.Cyan
    case 3:
        return color.Purple
    }
    return color.Red
}

func getNearPoints(x, y int) [4][2]int {
    p := [4][2]int{
            {x - 1, y},
            {x + 1, y},
            {x, y - 1},
            {x, y + 1},
        }
    return p
}

type Field struct {
    xSize, ySize int
    heights      []byte
    basins       []int
    basinSizes   map[int]int
}
    
    func (f *Field) getPoint(x, y int) (byte, bool) {
        if x < 0 || x >= f.xSize || y < 0 || y >= f.ySize {
            return 0, false
        }
        return f.heights[y * f.xSize + x], true
    }    
    func (f *Field) getLowPoint(x, y int) (bool, byte) {
        testPoints := getNearPoints(x, y)
        selfHeight, ok := f.getPoint(x, y)
        if !ok {
            return false, selfHeight
        }
        for _, p := range testPoints {
            h, ok := f.getPoint(p[0], p[1])
            if !ok {
                continue
            }
            if selfHeight >= h {
                return false, selfHeight
            }
        }
        return true, selfHeight
    }

    /*
     * Task 2
     */
    func (f *Field) debug2() {
        for x := 0; x < f.xSize; x++ {
            for y := 0; y < f.ySize; y++ {
                height, _ := f.getPoint(x, y)
                basin := f.getPointBasin(x, y)
                if basin > 0 {
                    fmt.Printf(color.Ize(getColor(basin), "%d"), height)
                } else {
                    fmt.Printf("%d", height)
                }
            }
            fmt.Printf("\r\n")
        }
        //fmt.Println(f.basinSizes)
    }
    func (f *Field) findAllBasins() {
        var lastBasin int
        f.basinSizes = make(map[int]int)
        for x := 0; x < f.xSize; x++ {
            for y := 0; y < f.ySize; y++ {
                h, _ := f.getPoint(x, y)
                if h < 9 && f.getPointBasin(x, y) == 0 {
                    lastBasin++
                    f.exploreBasin(x, y, lastBasin)
                }
            }
        }
    }
    func (f *Field) exploreBasin(x, y, basin int) {
        f.markPointBasin(x, y, basin)
        testPoints := getNearPoints(x, y)
        for _, p := range testPoints {
            h, ok := f.getPoint(p[0], p[1])
            if !ok || h == 9 || f.getPointBasin(p[0], p[1]) > 0 {
                continue
            }
            f.exploreBasin(p[0], p[1], basin)
        }
    }
    func (f *Field) getPointBasin(x, y int) int {
        return f.basins[y * f.xSize + x]
    }
    func (f *Field) markPointBasin(x, y, basin int) {
        f.basins[y * f.xSize + x] = basin
        f.basinSizes[basin]++
    }
    func (f *Field) getMultiplyOfBiggestBasins() int {
        var biggestBasins [3]int
        for _, basinSize := range f.basinSizes {
            for i := 0; i < 3; i++ {
                if biggestBasins[i] < basinSize {
                    biggestBasins[i] = basinSize
                    sort.Ints(biggestBasins[0:3])
                    break
                }
            }
        }
        return biggestBasins[0] * biggestBasins[1] * biggestBasins[2]
    }

    /*
     * Task 1
     */
    func (f *Field) debug1() {
        for x := 0; x < f.xSize; x++ {
            for y := 0; y < f.ySize; y++ {
                hasRisk, risk := f.getRisk(x, y)
                if hasRisk {
                    fmt.Printf(color.Ize(color.Red, "%d"), risk-1)
                } else {
                    fmt.Printf("%d", risk)
                }
            }
            fmt.Printf("\r\n")
        }
    }
    func (f *Field) countAllRisk() int {
        var totalRisk int
        for x := 0; x < f.xSize; x++ {
            for y := 0; y < f.ySize; y++ {
                hasRisk, risk := f.getRisk(x, y)
                if hasRisk {
                    totalRisk += int(risk)
                }
            }
        }
        return totalRisk
    }
    func (f *Field) getRisk(x, y int) (bool, byte) {
        isLow, height := f.getLowPoint(x, y)
        if !isLow {
            return false, height
        }
        return true, height + 1
    }

// File reading
func readMap(f *Field) {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {
        numbers := strings.Split(scanner.Text(), "")

        row := []byte{}
        for _, n := range numbers {
            row = append(row, byte(convInt(n)))
        }

        if f.xSize == 0 {
            f.xSize = len(row)
        } else if f.xSize != len(row) {
            panic("Row with incorrect length found")
        }
        f.ySize++

        f.heights = append(f.heights, row...)
    }
    f.basins = make([]int, len(f.heights))
}

//
func main() {
    var field Field

    readMap(&field)

    task1 := field.countAllRisk()
    fmt.Printf("Task 1: %d\r\n===\r\n", task1)

    field.findAllBasins()
    field.debug2()
    task2 := field.getMultiplyOfBiggestBasins()
    fmt.Printf("Task 2: %d\r\n===\r\n", task2)

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