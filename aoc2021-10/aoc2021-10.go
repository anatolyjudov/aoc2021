package main

import (
    "fmt"
    "os"
    "bufio"
    "sort"
)

var task1sum int

var task2lines []int

var closedTags = map[string]string{
    "(": ")",
    "[": "]",
    "{": "}",
    "<": ">",
}

func getClosedByte(b byte) byte {
    r := []byte(closedTags[string([]byte{b})])
    return r[0]
}

type Stack struct {
    data []byte
    pos int
}
    func (s *Stack) push(v string) {
        s.pos++
        if (len(s.data) <= s.pos-1) {
            s.data = append(s.data, byt(v))
        } else {
            s.data[s.pos-1] = byt(v)
        }
    }
    func (s *Stack) pop() string {
        ret := char(s.data[s.pos-1])
        s.data[s.pos-1] = 0
        s.pos--
        return ret
    }
    func (s *Stack) flush() []byte {
        var f []byte
        for i := 1; i < s.pos + 1; i++ {
            closedByte := getClosedByte(s.data[s.pos - i])
            f = append(f, closedByte)
        }
        return f
    }

func flushCost(b []byte) int {
    var s int
    for i := 0; i < len(b); i++ {
        s = s * 5 + task2Cost(char(b[i]))
    }
    return s
}

func char(b byte) string {
    return string([]byte{b})
}
func byt(s string) byte {
    b := []byte(s)
    return b[0]
}

func analyzeLine(line []byte) bool {
    var s Stack

    fmt.Printf("%v -> ", string(line))

    for i := 0; i < len(line); i++ {
        ch := char(line[i])
        if ch == "{" || ch == "[" || ch == "(" || ch == "<" {
            s.push(ch)
        } else if ch == "}" || ch == "]" || ch == ")" || ch == ">" {
            lastOpen := s.pop()
            expected := closedTags[lastOpen]
            if expected != ch {
                fmt.Printf("Error at %v, expected %v, %v found\r\n", i, expected, ch)
                task1sum += task1Cost(ch)
                return false
            }
        }
    }

    fmt.Printf("No errors found\r\n")
    flush := s.flush()
    thisFlushCost := flushCost(flush)
    fmt.Printf("flush: %v, costs: %v \r\n", string(flush), thisFlushCost)
    task2lines = append(task2lines, thisFlushCost)

    return true
}

func task1Cost(ch string) int {
    switch ch {
    case ")":
        return 3
    case "]":
        return 57
    case "}":
        return 1197
    case ">":
        return 25137
    }
    return 0
}

func task2Cost(ch string) int {
    switch ch {
    case ")":
        return 1
    case "]":
        return 2
    case "}":
        return 3
    case ">":
        return 4
    }
    return 0
}

//
func main() {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)

    task1sum = 0
    for scanner.Scan() {
        analyzeLine(scanner.Bytes())
    }

    fmt.Printf("\r\nTask 1: %d\r\n===\r\n", task1sum)

    sort.Ints(task2lines)
    task2 := task2lines[len(task2lines) / 2]
    fmt.Println(task2lines)
    fmt.Printf("\r\nTask 2: %d\r\n===\r\n", task2)
}

// file opening routine
func openFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        panic(err)
    }
    return file
}