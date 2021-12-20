package main

import (
    "fmt"
    "os"
    "bufio"
    //"strings"
    "strconv"
    "github.com/TwiN/go-color"
)

type FishNum5 struct {
    s [][]int
}
    func (fn *FishNum5) isNull() bool {
        for l := 0; l < 5; l++ {
            for n := 0; n < (1 << (l + 1)); n++ {
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
            var colorized string
            if fn.s[level][position] > -1 {
                return fmt.Sprintf("%d", fn.s[level][position])
            } else if level == 4 {
                return "null"
            }
            colorized = color.Ize(color.Gray, "[%v, %v]")
            if level == 3 {
                colorized = color.Ize(color.Red, "[%v, %v]")
            }
            return fmt.Sprintf(colorized, fn.printPart(level + 1, position << 1), fn.printPart(level + 1, position << 1 + 1))
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
    func (fn *FishNum5) reduce() {
        for ;; {
            if fn.explode() {
                continue
            }
            if fn.split() {
                continue
            }
            break
        }
    }
        func (fn *FishNum5) split() bool {
            for pos := 0; pos < (1 << 5); pos++ {
                if fn.splitAbove(pos) {
                    return true
                }
            }
            return false
        }
            func (fn *FishNum5) splitAbove(pos int) bool {
                for l := (5 - 2); l >= 0; l-- {
                    pos = pos >> 1
                    if fn.s[l][pos] >= 10 {
                        fn.s[l + 1][pos << 1] = fn.s[l][pos] >> 1
                        fn.s[l + 1][(pos << 1) + 1] = fn.s[l][pos] - (fn.s[l][pos] >> 1)
                        fn.s[l][pos] = -1
                        return true
                    }
                }
                return false
            }
        func (fn *FishNum5) explode() bool {
            ll := 5 - 1
            for i := 0; i < (1 << 5); i += 2 {
                if fn.s[ll][i] != -1 {
                    // explode i-position
                    fn.explodeLeft(i, fn.s[ll][i])
                    fn.explodeRight(i + 1, fn.s[ll][i + 1])
                    fn.s[ll][i] = -1
                    fn.s[ll][i + 1] = -1
                    fn.s[ll - 1][i >> 1] = 0
                    return true
                }
            }
            return false
        }
        func (fn *FishNum5) explodeLeft(pos int, value int) {
            ll := 5 - 1
            if pos == 0 {
                return
            }
            for l := ll; l >= 0; l-- {
                i := (pos - 1) >> (ll - l)
                if fn.s[l][i] != -1 {
                    fn.s[l][i] += value
                    return
                }
            }
        }
        func (fn *FishNum5) explodeRight(pos int, value int) {
            ll := 5 - 1
            if pos == ((1 << 5) - 1) {
                return
            }
            for l := ll; l >= 0; l-- {
                i := (pos + 1) >> (ll - l)
                if fn.s[l][i] != -1 {
                    fn.s[l][i] += value
                    return
                }
            }
        }
    func (fn *FishNum5) magnitude() int {
        return fn.magnitudePos(0, 0) * 3 + fn.magnitudePos(0, 1) * 2
    }
        func (fn *FishNum5) magnitudePos(level, pos int) int {
            if fn.s[level][pos] > -1 {
                return fn.s[level][pos]
            }
            return fn.magnitudePos(level + 1, pos << 1) * 3 + fn.magnitudePos(level + 1, (pos << 1) + 1) * 2
        }

func add(f1, f2 FishNum5) (f FishNum5) {
    f = makeFishNum5()
    for l := 0; l < (5 - 1); l++ {
        for n := 0; n < (1 << (l + 1)); n++ {
            f.s[l + 1][n] = f1.s[l][n]
            f.s[l + 1][n + 1 << (l + 1)] = f2.s[l][n]
        }
    }
    f.reduce()
    return f
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

    sum := fns[0]
    for i := 1; i < len(fns); i++ {
        sum = add(sum, fns[i])
    }
    task1 := sum.magnitude()
    fmt.Println(task1)

    task2 := task1
    for n1 := 0; n1 < len(fns); n1++ {
        for n2 := 0; n2 < len(fns); n2++ {
            if n2 == n1 {
                continue
            }
            sum = add(fns[n1], fns[n2])
            task1 = sum.magnitude()
            if task1 > task2 {
                task2 = task1
            }
        }
    }
    fmt.Println(task2)

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
