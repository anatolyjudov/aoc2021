package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
    //"strconv"
    "sort"
)

const A = 97
const B = 98
const C = 99
const D = 100
const E = 101
const F = 102
const G = 103

func getDigit(digit string) byte {
    var segments = strings.Split(digit, "")
    sort.Strings(segments)
    digit = strings.Join(segments, "")
    switch digit {
    case "abcefg":
        return 0
    case "cf":
        return 1
    case "acdeg":
        return 2
    case "acdfg":
        return 3
    case "bcdf":
        return 4
    case "abdfg":
        return 5
    case "abdefg":
        return 6
    case "acf":
        return 7
    case "abcdefg":
        return 8
    case "abcdfg":
        return 9
    }

    fmt.Println(digit)
    panic("Can't recognize digit")
    return 0
}

type Display struct {
    patterns [10][]byte
    output   [4][]byte
    key      map[byte]byte
}

    // Prints all info about Display (useful for debug)
    func (d *Display) print() {
        for w := 0; w < 10; w++ {
            fmt.Print(string(d.patterns[w]))
            fmt.Print(" ")
        }
        fmt.Print(" -> ")
        for w := 0; w < 4; w++ {
            fmt.Print(string(d.output[w]))
            fmt.Print(" ")
        }
        fmt.Print("(")
        for i, v := range d.key {
            fmt.Printf("%c -> %c,", i, v)
        }
        fmt.Print(")")
        fmt.Print(" --> ")
        fmt.Println(d.decode())
        fmt.Print("\r\n")
    }
    
    // Finds a key for decoding the display
    func (d *Display) findKey() {
        var counts map[byte]int
        counts = make(map[byte]int)
        for p := 0; p < 10; p++ {
            for _, b := range d.patterns[p] {
                counts[b]++
            }
        }
        d.key = make(map[byte]byte)
        // E, B, F
        for w, count := range counts {
            switch count {
            case 4:
                d.key[E] = w
            case 6:
                d.key[B] = w
            case 9:
                d.key[F] = w
            }
        }
        // C
        for _, p := range d.patterns {
            if len(p) != 2 {
                continue
            }
            if p[0] == d.key[F] {
                d.key[C] = p[1]
            } else if p[1] == d.key[F] {
                d.key[C] = p[0]
            } else {
                panic("C search failed")
            }
        }
        // A, D
        for _, p := range d.patterns {
            if len(p) == 3 {
                for i := 0; i < 3; i++ {
                    if p[i] == d.key[C] || p[i] == d.key[F] {
                        continue
                    }
                    d.key[A] = p[i]
                    break
                }
            }
            if len(p) == 4 {
                for i := 0; i < 4; i++ {
                    if p[i] == d.key[B] || p[i] == d.key[C] || p[i] == d.key[F] {
                        continue
                    }
                    d.key[D] = p[i]
                    break
                }
            }
        }
        // G
        for _, p := range d.patterns {
            if len(p) == 7 {
                for i := 0; i < 7; i++ {
                    if p[i] == d.key[A] || p[i] == d.key[B] || p[i] == d.key[C] || p[i] == d.key[D] || p[i] == d.key[E] || p[i] == d.key[F] {
                        continue
                    }
                    d.key[G] = p[i]
                    break
                }
            }
        }

        rkey := make(map[byte]byte)
        for a, b := range d.key {
            rkey[b] = a
        }
        d.key = rkey
    }

    // Decodes the output using the key
    func (d *Display) decode() int {
        var segments [4]byte
        var b strings.Builder
        for i := 0; i < 4; i++ {
            b.Reset()
            for _, o := range d.output[i] {
                b.WriteByte(d.key[o])
            }
            segments[i] = getDigit(b.String())
        }
        output := 0
        for i := 0; i < 4; i++ {
            if (i == 0) {
                output += int(segments[i])
                continue
            }
            output = output * 10 + int(segments[i])
        }
        return output
    }

// count digits 1, 4, 7, or 8 in the output
// task 1 could be made without decoding
func countTask1(displays *[]Display) int {
    var count int
    for _, d := range *displays {
        for i := 0; i < 4; i++ {
            l := len(d.output[i])
            if l == 2 || l == 4 || l == 3 || l == 7 {
                count++
            }
        }
    }
    return count
}

// File reading
func readDisplays() []Display {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    var displays []Display
    for scanner.Scan() {
        displays = append(displays, parseDisplay(scanner.Text()))
    }
    return displays
}
    func parseDisplay(data string) Display {
        var d Display
        parts := strings.Split(data, " ")
        for w := 0; w < 10; w++ {
            d.patterns[w] = []byte(parts[w])
        }
        for w := 11; w < 15; w++ {
            d.output[w - 11] = []byte(parts[w])
        }
        return d
    }

//
func main() {
    displays := readDisplays()

    task1 := countTask1(&displays)

    fmt.Printf("Task 1: %d\r\n===\r\n", task1)

    task2 := 0
    for _, display := range displays {
        display.findKey()
        task2 += display.decode()
    }
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