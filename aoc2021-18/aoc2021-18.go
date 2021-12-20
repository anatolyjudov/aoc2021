package main

import (
    "fmt"
    "os"
    "bufio"
    //"strings"
    "strconv"
)

type FishNum5 struct {
    s [][]int
}
    func (fn *FishNum5) isNull() bool {
        for l := 0; l < 5; l++ {
            for n := 0; n < (1 << l); n++ {
                if fn.s[l][n] != -1 {
                    return false
                }
            }
        }
        return true
    }
    func (fn FishNum5) String() (ret string) {
        if fn.isNull() {
            return "null"
        }
        return fmt.Sprintf("[%v, %v]", fn.printPart(0, 0), fn.printPart(0, 1))
    }
        func (fn FishNum5) printPart(level, position int) (ret string) {
            if fn.s[level][position] > -1 {
                return fmt.Sprintf("%d", fn.s[level][position])
            } else if level == 4 {
                return "null"
            }
            return fmt.Sprintf("[%v, %v]", fn.printPart(level + 1, position << 1), fn.printPart(level + 1, position << 1 + 1))
        }
    func (fn *FishNum5) readToPosition(level, position int, str string, strPos int) int {
        if string(str[strPos]) == "[" {
            // it's a pair
            strPos++
            strPos = fn.readToPosition(level + 1, position << 1, str, strPos)
            if string(str[strPos]) != "," {
                panic(fmt.Sprintf("Incorrect character: divider symbol must be , at %v", strPos))
            }
            strPos++
            strPos = fn.readToPosition(level + 1, position << 1 + 1, str, strPos)
            if string(str[strPos]) != "]" {
                panic(fmt.Sprintf("Incorrect character: closing symbol must be ] at %v", strPos))
            }
            strPos++
            return strPos
        }
        // it's a literal
        pos0 := strPos
        for string(str[strPos]) != "," && string(str[strPos]) != "]" {
            strPos++
        }
        literal := string(str[pos0:strPos])
        fn.s[level][position] = convInt(literal)
        return strPos
    }

func makeFishNum5fromString(str string) FishNum5 {
    fn := makeFishNum5()
    strPos := 0
    if string(str[strPos]) != "[" {
        panic("Incorrect character: first symbol must be [")
    }
    strPos++
    strPos = fn.readToPosition(0, 0, str, strPos)
    if string(str[strPos]) != "," {
        panic("Incorrect character: divider symbol must be ,")
    }
    strPos++
    strPos = fn.readToPosition(0, 1, str, strPos)
    if string(str[strPos]) != "]" {
        panic("Incorrect character: last symbol must be ,")
    }
    return fn
}

func makeFishNum5() FishNum5 {
    var fn FishNum5
    for l := 0; l < 5; l++ {
        level := make([]int, 1 << (l + 1))
        for n := 0; n < len(level); n++ {
            level[n] = -1
        }
        fn.s = append(fn.s, level)
    }
    return fn
}

func main() {
    fns := readInput()

    for _, fn := range fns {
        fmt.Println(fn)
    }
}


// read all input data
func readInput() (fns []FishNum5) {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fns = append(fns, makeFishNum5fromString(scanner.Text()))
    }
    return
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
