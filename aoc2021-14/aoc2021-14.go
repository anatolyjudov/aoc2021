package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

func apply(polymer string, subs *map[[2]byte]byte) string {
    var newPolymer strings.Builder

    polymerBytes := []byte(polymer)

    subss := *subs

    for i := 0; i < len(polymerBytes); i++ {
        if i > 0 {
            pair := [2]byte{polymerBytes[i - 1], polymerBytes[i]}
            if insert, ok := subss[pair]; ok {
                newPolymer.WriteByte(insert)
                newPolymer.WriteByte(polymerBytes[i])
                continue
            }
        }
        newPolymer.WriteByte(polymerBytes[i])
    }

    return newPolymer.String()
}

func task1count(polymer string) int {
    var counts map[byte]int

    counts = make(map[byte]int)

    polymerBytes := []byte(polymer)
    for i := 0; i < len(polymerBytes); i++ {
        if c, ok := counts[polymerBytes[i]]; ok {
            counts[polymerBytes[i]] = c + 1
        } else {
            counts[polymerBytes[i]] = 1
        }
    }

    var (
        max, min int
    )
    for _, count := range counts {
        if count > max {
            max = count
        }
        if min == 0 || count < min {
            min = count
        }
    }

    return max - min
}

func main() {

    polymer, subs := readInput()

    fmt.Println(polymer, subs)

    for step := 0; step < 40; step++ {
        polymer = apply(polymer, &subs)
        //fmt.Printf("\r\nStep %v:\r\n%v", step, polymer)
    }
    
    //fmt.Printf("\r\nFinal polymer:\r\n%v", polymer)

    fmt.Printf("\r\nTask 1 count: %v", task1count(polymer))    
}

// read all input data
func readInput() (template string, subs map[[2]byte]byte) {
    file := openFile("./input.txt")
    scanner := bufio.NewScanner(file)
    // template
    scanner.Scan()
    template = scanner.Text()
    // skip empty string
    scanner.Scan()
    // substitutions
    subs = make(map[[2]byte]byte)
    for scanner.Scan() {
        p := strings.Split(scanner.Text(), " -> ")
        pair := [2]byte{[]byte(p[0])[0], []byte(p[0])[1]}
        subs[pair] = []byte(p[1])[0]
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
